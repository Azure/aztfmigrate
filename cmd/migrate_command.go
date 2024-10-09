package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/aztfmigrate/types"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
)

const filenameImport = "imports.tf"

const tempDir = "temp"

type MigrateCommand struct {
	Ui             cli.Ui
	Verbose        bool
	Strict         bool
	workingDir     string
	TargetProvider string
}

func (c *MigrateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.Verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.Strict, "strict", false, "strict mode: API versions must be matched")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.StringVar(&c.TargetProvider, "to", "", "Specify the provider to migrate to. The allowed values are: azurerm and azapi. Default is azurerm.")

	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c *MigrateCommand) Run(args []string) int {
	// AzureRM provider will honor env.var "AZURE_HTTP_USER_AGENT" when constructing for HTTP "User-Agent" header.
	// #nosec G104
	_ = os.Setenv("AZURE_HTTP_USER_AGENT", "mig")
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
	tempDirectoryCreate(workingDirectory)
	tempPath := filepath.Join(workingDirectory, tempDir)
	tempTerraform, err := tf.NewTerraform(tempPath, c.Verbose)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] generating import config...")
	config := importConfig(resources, c.TargetProvider)
	if err = os.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0600); err != nil {
		log.Fatal(err)
	}

	log.Printf("[INFO] migrating resources...")
	for _, r := range resources {
		log.Printf("[INFO] generating new config for resource %s...", r.OldAddress(nil))
		if err := r.GenerateNewConfig(tempTerraform); err != nil {
			log.Printf("[ERROR] %+v", err)
		}
	}

	tempDirectoryCleanup(workingDirectory)

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
			importBlock := r.ImportBlock()
			removedBlock := r.RemovedBlock()
			if err := types.ReplaceResourceBlock(workingDirectory, r.OldAddress(nil), []*hclwrite.Block{removedBlock, importBlock, r.MigratedBlock()}); err != nil {
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

func importConfig(resources []types.AzureResource, targetProvider string) string {
	const providerConfig = `
provider "azurerm" {
  features {}
  subscription_id = "%s"
}
`
	const AzapiProviderConfig = `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}
`

	config := ""
	for _, r := range resources {
		config += r.EmptyImportConfig()
	}
	if targetProvider == "azurerm" {
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
		config = fmt.Sprintf(providerConfig, subscriptionId) + config
	} else {
		config = AzapiProviderConfig + config
	}

	return config
}

func tempDirectoryCreate(workingDirectory string) {
	tempPath := filepath.Join(workingDirectory, tempDir)
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(tempPath); err != nil {
			log.Fatalf("error deleting %s: %+v", tempPath, err)
		}
	}
	if err := os.MkdirAll(tempPath, 0750); err != nil {
		log.Fatalf("creating temp workspace %q: %+v", tempPath, err)
	}
}

func tempDirectoryCleanup(workingDirectory string) {
	tempPath := filepath.Join(workingDirectory, tempDir)
	// cleanup temp folder
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(tempPath); err != nil {
			log.Printf("[WARN] error deleting %s: %+v", tempPath, err)
		}
	}
}
