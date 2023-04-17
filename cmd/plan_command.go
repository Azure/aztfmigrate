package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Azure/azapi2azurerm/azurerm"
	"github.com/Azure/azapi2azurerm/azurerm/coverage"
	"github.com/Azure/azapi2azurerm/helper"
	"github.com/Azure/azapi2azurerm/tf"
	"github.com/Azure/azapi2azurerm/types"
	"github.com/mitchellh/cli"
)

type PlanCommand struct {
	Ui      cli.Ui
	Verbose bool
	Strict  bool
}

func (c *PlanCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.Verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.Strict, "strict", false, "strict mode: API versions must be matched")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c PlanCommand) Run(args []string) int {
	// AzureRM provider will honor env.var "AZURE_HTTP_USER_AGENT" when constructing for HTTP "User-Agent" header.
	// #nosec G104
	os.Setenv("AZURE_HTTP_USER_AGENT", "mig")
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Printf("[INFO] initializing terraform...")
	workingDirectory, _ := os.Getwd()
	terraform, err := tf.NewTerraform(workingDirectory, c.Verbose, false)
	if err != nil {
		log.Fatal(err)
	}
	c.Plan(terraform, true)
	return 0
}

func (c PlanCommand) Help() string {
	helpText := `
Usage: azapi2azurerm plan
` + c.Synopsis() + "\nThe Terraform addresses listed in file `azapi2azurerm.ignore` will be ignored during migration.\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c PlanCommand) Synopsis() string {
	return "Show azapi resources which can migrate to azurerm resources in current working directory"
}

func (c PlanCommand) Plan(terraform *tf.Terraform, isPlanOnly bool) ([]types.GenericResource, []types.GenericUpdateResource) {
	// get azapi resource from state
	log.Printf("[INFO] searching azapi_resource & azapi_update_resource...")
	p, err := terraform.Plan()
	if err != nil {
		log.Fatal(err)
	}

	migrationMessage := "The following resources will be migrated:\n"
	unsupportedMessage := "The following resources can't be migrated:\n"
	ignoreMessage := "The following resources will be ignored in migration:\n"
	ignoreSet := make(map[string]bool)
	if file, err := os.ReadFile(path.Join(terraform.GetWorkingDirectory(), "azapi2azurerm.ignore")); err == nil {
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

	resources := make([]types.GenericResource, 0)
	updateResources := make([]types.GenericUpdateResource, 0)
	for _, resource := range terraform.ListGenericResources(p) {
		if ignoreSet[resource.OldAddress(nil)] || len(resource.Instances) == 0 {
			continue
		}
		resourceId := resource.Instances[0].ResourceId
		resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resourceId)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get resource type for %s: %w", resourceId, err))
		}

		idPattern, _ := helper.GetIdPattern(resourceId)
		if c.Strict {
			azurermApiVersion := coverage.GetApiVersion(idPattern)
			if azurermApiVersion != resource.Instances[0].ApiVersion {
				unsupportedMessage += fmt.Sprintf("\t%s: api-versions are not matched, expect %s, got %s\n",
					resource.OldAddress(nil), resource.Instances[0].ApiVersion, azurermApiVersion)
				continue
			}
		}

		_, uncoveredPut := coverage.GetPutCoverage(resource.InputProperties, idPattern)
		_, uncoveredGet := coverage.GetGetCoverage(resource.OutputProperties, idPattern)

		if len(uncoveredGet)+len(uncoveredPut) == 0 {
			if !isPlanOnly {
				if exact {
					resource.ResourceType = resourceTypes[0]
				} else {
					resource.ResourceType = c.getUserInputResourceType(resourceId, resourceTypes)
				}
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(nil), resource.ResourceType)
			} else {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(nil), strings.Join(resourceTypes, ", "))
			}
			resources = append(resources, resource)
		} else {
			unsupportedMessage += fmt.Sprintf("\t%s: input properties not supported: [%v], output properties not supported: [%v]\n",
				resource.OldAddress(nil), strings.Join(uncoveredPut, ", "), strings.Join(uncoveredGet, ", "))
		}
	}

	for _, resource := range terraform.ListGenericUpdateResources(p) {
		if ignoreSet[resource.OldAddress()] {
			continue
		}
		resourceTypes, exact, err := azurerm.GetAzureRMResourceType(resource.Id)
		if err != nil {
			log.Fatal(fmt.Errorf("failed to get resource type for %s: %w", resource.Id, err))
		}

		idPattern, _ := helper.GetIdPattern(resource.Id)
		if c.Strict {
			azurermApiVersion := coverage.GetApiVersion(idPattern)
			if azurermApiVersion != resource.ApiVersion {
				unsupportedMessage += fmt.Sprintf("\t%s: api-versions are not matched, expect %s, got %s\n",
					resource.OldAddress(), resource.ApiVersion, azurermApiVersion)
				continue
			}
		}
		_, uncoveredPut := coverage.GetPutCoverage(resource.InputProperties, idPattern)
		_, uncoveredGet := coverage.GetGetCoverage(resource.OutputProperties, idPattern)

		if len(uncoveredGet)+len(uncoveredPut) == 0 {
			if !isPlanOnly {
				if exact {
					resource.ResourceType = resourceTypes[0]
				} else {
					resource.ResourceType = c.getUserInputResourceType(resource.Id, resourceTypes)
				}
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(), resource.ResourceType)
			} else {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(), strings.Join(resourceTypes, ", "))
			}
			updateResources = append(updateResources, resource)
		} else {
			unsupportedMessage += fmt.Sprintf("\t%s: input properties not supported: [%v], output properties not supported: [%v]\n",
				resource.OldAddress(), strings.Join(uncoveredPut, ", "), strings.Join(uncoveredGet, ", "))
		}
	}

	log.Printf("[INFO]\n\nThe tool will perform the following actions:\n\n%s\n%s\n%s\n", migrationMessage, unsupportedMessage, ignoreMessage)
	return resources, updateResources
}

func (c PlanCommand) getUserInputResourceType(resourceId string, values []string) string {
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
