package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Azure/aztfmigrate/azurerm"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/aztfmigrate/types"
	"github.com/mitchellh/cli"
)

type PlanCommand struct {
	Ui             cli.Ui
	Verbose        bool
	Strict         bool
	workingDir     string
	varFile        string
	TargetProvider string
}

func (c *PlanCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.Verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.Strict, "strict", false, "strict mode: API versions must be matched")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.StringVar(&c.varFile, "var-file", "", "path to the terraform variable file")
	fs.StringVar(&c.TargetProvider, "to", "", "Specify the provider to migrate to. The allowed values are: azurerm and azapi. Default is azurerm.")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c *PlanCommand) Run(args []string) int {
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
	log.Printf("[INFO] target provider: %s", c.TargetProvider)

	log.Printf("[INFO] initializing terraform...")
	if c.workingDir == "" {
		c.workingDir, _ = os.Getwd()
	}
	terraform, err := tf.NewTerraform(c.workingDir, c.Verbose)
	if err != nil {
		log.Fatal(err)
	}
	c.Plan(terraform, true)
	return 0
}

func (c *PlanCommand) Help() string {
	helpText := `
Usage: aztfmigrate plan
` + c.Synopsis() + "\nThe Terraform addresses listed in file `aztfmigrate.ignore` will be ignored during migration.\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c *PlanCommand) Synopsis() string {
	return "Show terraform resources which can be migrated to azurerm or azapi resources in current working directory"
}

func (c *PlanCommand) Plan(terraform *tf.Terraform, isPlanOnly bool) []types.AzureResource {
	// get azapi resource from state
	log.Printf("[INFO] running terraform plan...")
	p, err := terraform.Plan(&c.varFile)
	if err != nil {
		log.Fatal(err)
	}

	migrationMessage := "The following resources will be migrated:\n"
	unsupportedMessage := "The following resources can't be migrated:\n"
	ignoreMessage := "The following resources will be ignored in migration:\n"
	ignoreSet := make(map[string]bool)
	if file, err := os.ReadFile(path.Join(terraform.GetWorkingDirectory(), "aztfmigrate.ignore")); err == nil {
		lines := strings.Split(string(file), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if len(line) == 0 {
				continue
			}
			ignoreSet[line] = true
			ignoreMessage += fmt.Sprintf("\t%s\n", line)
		}
	}

	res := make([]types.AzureResource, 0)

	for _, item := range types.ListResourcesFromPlan(p) {
		if item.TargetProvider() != c.TargetProvider {
			continue
		}
		if ignoreSet[item.OldAddress(nil)] {
			continue
		}
		if err := item.CoverageCheck(c.Strict); err != nil {
			unsupportedMessage += fmt.Sprintf("\t%s\n", err)
			continue
		}

		switch resource := item.(type) {
		case *types.AzapiResource:
			if len(resource.Instances) == 0 {
				continue
			}
			resourceId := resource.Instances[0].ResourceId
			resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resourceId)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to get resource type for %s: %w", resourceId, err))
			}
			if exact {
				resource.ResourceType = resourceTypes[0]
			} else if !isPlanOnly {
				resource.ResourceType = c.getUserInputResourceType(resourceId, resourceTypes)
			}

			if resource.ResourceType != "" {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %s\n", resource.OldAddress(nil), resource.NewAddress(nil))
			} else {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(nil), strings.Join(resourceTypes, ", "))
			}
			res = append(res, resource)

		case *types.AzapiUpdateResource:
			resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resource.Id)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to get resource type for %s: %w", resource.Id, err))
			}

			if exact {
				resource.ResourceType = resourceTypes[0]
			} else if !isPlanOnly {
				resource.ResourceType = c.getUserInputResourceType(resource.Id, resourceTypes)
			}

			if resource.ResourceType != "" {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %s\n", resource.OldAddress(nil), resource.NewAddress(nil))
			} else {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(nil), strings.Join(resourceTypes, ", "))
			}
			res = append(res, resource)

		case *types.AzurermResource:
			if len(resource.Instances) == 0 {
				continue
			}
			migrationMessage += fmt.Sprintf("\t%s will be replaced with %s\n", resource.OldAddress(nil), resource.NewAddress(nil))
			res = append(res, resource)
		}
	}

	log.Printf("[INFO]\n\nThe tool will perform the following actions:\n\n%s\n%s\n%s\n", migrationMessage, unsupportedMessage, ignoreMessage)
	return res
}

func (c *PlanCommand) getUserInputResourceType(resourceId string, values []string) string {
	c.Ui.Warn(fmt.Sprintf("Couldn't find unique resource type for id: %s\nPossible values are [%s].\nPlease input an azurerm resource type:", resourceId, strings.Join(values, ", ")))
	resourceType := ""
	for {
		reader := bufio.NewReader(os.Stdin)
		resourceType, _ = reader.ReadString('\n')
		resourceType = strings.Trim(resourceType, "\r\n")
		for _, value := range values {
			if value == resourceType {
				return resourceType
			}
		}
		c.Ui.Warn("Invalid input. Please input an azurerm resource type:")
	}
}
