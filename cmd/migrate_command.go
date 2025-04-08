package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Azure/aztfmigrate/helper"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/aztfmigrate/types"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
)

const filenameImport = "imports.tf"

const tempFolderName = "aztfmigrate_temp"

type MigrateCommand struct {
	Ui             cli.Ui
	Verbose        bool
	Strict         bool
	workingDir     string
	varFile        string
	TargetProvider string
}

func (c *MigrateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.Verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.Strict, "strict", false, "strict mode: API versions must be matched")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.StringVar(&c.varFile, "var-file", "", "path to the terraform variable file")
	fs.StringVar(&c.TargetProvider, "to", "", "Specify the provider to migrate to. The allowed values are: azurerm and azapi. Default is azurerm.")

	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c *MigrateCommand) Run(args []string) int {
	// AzureRM provider will honor env.var "AZURE_HTTP_USER_AGENT" when constructing for HTTP "User-Agent" header.
	// #nosec G104
	_ = os.Setenv("AZURE_HTTP_USER_AGENT", "mig")
	// The following env.vars are used to disable enhanced validation and skip provider registration, to speed up the process.
	// #nosec G104
	_ = os.Setenv("ARM_PROVIDER_ENHANCED_VALIDATION", "false")
	// #nosec G104
	_ = os.Setenv("ARM_SKIP_PROVIDER_REGISTRATION", "true")
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	if c.TargetProvider == "" {
		c.TargetProvider = "azurerm"
	}
	if c.TargetProvider != "azapi" && c.TargetProvider != "azurerm" {
		c.Ui.Error("Invalid target provider. The allowed values are: azurerm and azapi.")
		return 1
	}

	if c.workingDir == "" {
		c.workingDir, _ = os.Getwd()
	}
	log.Printf("[INFO] working directory: %s", c.workingDir)

	log.Printf("[INFO] initializing terraform...")
	terraform, err := tf.NewTerraform(c.workingDir, c.Verbose)
	if err != nil {
		log.Fatal(err)
	}

	planCommand := &PlanCommand{ //nolint
		Ui:             c.Ui,
		Verbose:        c.Verbose,
		Strict:         c.Strict,
		workingDir:     c.workingDir,
		varFile:        c.varFile,
		TargetProvider: c.TargetProvider,
	}
	allResources := planCommand.Plan(terraform, false)
	c.MigrateResources(terraform, allResources)
	return 0
}

func (c *MigrateCommand) Help() string {
	helpText := `
Usage: aztfmigrate migrate
` + c.Synopsis() + "\nThe Terraform addresses listed in file `aztfmigrate.ignore` will be ignored during migration.\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c *MigrateCommand) Synopsis() string {
	return "Migrate azapi resources to azurerm resources in current working directory"
}

func (c *MigrateCommand) MigrateResources(terraform *tf.Terraform, resources []types.AzureResource) {
	if len(resources) == 0 {
		return
	}

	workingDirectory := terraform.GetWorkingDirectory()
	// write empty config to temp dir for import
	tempDir := filepath.Join(workingDirectory, tempFolderName)
	if err := os.MkdirAll(tempDir, 0750); err != nil {
		log.Fatalf("creating temp workspace %q: %+v", tempDir, err)
	}
	if err := os.RemoveAll(path.Join(tempDir, "terraform.tfstate")); err != nil {
		log.Printf("[WARN] removing temp workspace %q: %+v", tempDir, err)
	}
	defer func() {
		err := os.RemoveAll(path.Join(tempDir, "terraform.tfstate"))
		if err != nil {
			log.Printf("[ERROR] removing temp workspace %q: %+v", tempDir, err)
		}
	}()
	tempTerraform, err := tf.NewTerraform(tempDir, c.Verbose)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] generating import config...")
	config := ImportConfig(resources, helper.FindHclBlock(workingDirectory, "terraform", nil))
	if err = os.WriteFile(filepath.Join(tempDir, filenameImport), []byte(config), 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] migrating resources...")
	for _, r := range resources {
		log.Printf("[INFO] generating new config for resource %s...", r.OldAddress(nil))
		if err := r.GenerateNewConfig(tempTerraform); err != nil {
			log.Printf("[ERROR] %+v", err)
		}
	}

	log.Printf("[INFO] updating config...")
	updateResources := make([]types.AzapiUpdateResource, 0)
	for _, r := range resources {
		if updateResource, ok := r.(*types.AzapiUpdateResource); ok {
			updateResources = append(updateResources, *updateResource)
		}
	}
	if err := types.UpdateMigratedResourceBlock(workingDirectory, updateResources); err != nil {
		log.Fatal(err)
	}

	// migrate depends_on, lifecycle, provisioner
	for _, r := range resources {
		if existingBlock, err := types.GetResourceBlock(workingDirectory, r.OldAddress(nil)); err == nil && existingBlock != nil {
			migratedBlock := r.MigratedBlock()
			if attr := existingBlock.Body().GetAttribute("depends_on"); attr != nil {
				migratedBlock.Body().SetAttributeRaw("depends_on", attr.Expr().BuildTokens(nil))
			}
			for _, block := range existingBlock.Body().Blocks() {
				if block.Type() == "lifecycle" || block.Type() == "provisioner" {
					migratedBlock.Body().AppendBlock(block)
				}
			}
		}
	}

	// remove from config
	for _, r := range resources {
		if r.IsMigrated() {
			log.Printf("[INFO] removing %s from config", r.OldAddress(nil))
			stateUpdateBlocks := r.StateUpdateBlocks()
			newBlocks := make([]*hclwrite.Block, 0)
			newBlocks = append(newBlocks, stateUpdateBlocks...)
			newBlocks = append(newBlocks, r.MigratedBlock())
			if err := types.ReplaceResourceBlock(workingDirectory, r.OldAddress(nil), newBlocks); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(nil), err)
			}
		}
	}

	log.Printf("[INFO] replacing references with migrated resource...")
	outputs := make([]types.Output, 0)
	for _, r := range resources {
		if r.IsMigrated() {
			outputs = append(outputs, r.Outputs()...)
		}
	}
	if err := types.ReplaceGenericOutputs(workingDirectory, outputs); err != nil {
		log.Printf("[ERROR] replacing outputs: %+v", err)
	}
}

func ImportConfig(resources []types.AzureResource, terraformBlock *hclwrite.Block) string {
	config := `terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}`
	if terraformBlock != nil {
		newFile := hclwrite.NewEmptyFile()
		newFile.Body().AppendBlock(terraformBlock)
		config = string(hclwrite.Format(newFile.Bytes()))
	}

	for _, r := range resources {
		config += r.EmptyImportConfig()
	}
	subscriptionId := ""
	for _, r := range resources {
		switch resource := r.(type) {
		case *types.AzapiResource:
			for _, instance := range resource.Instances {
				if strings.HasPrefix(instance.ResourceId, "/subscriptions/") {
					subscriptionId = strings.Split(instance.ResourceId, "/")[2]
					break
				}
			}
		case *types.AzapiUpdateResource:
			if strings.HasPrefix(resource.Id, "/subscriptions/") {
				subscriptionId = strings.Split(resource.Id, "/")[2]
			}
		}
		if subscriptionId != "" {
			break
		}
	}

	const providerConfig = `
provider "azurerm" {
  features {}
  subscription_id = "%s"
}

provider "azapi" {
}
`

	return fmt.Sprintf(providerConfig, subscriptionId) + config
}
