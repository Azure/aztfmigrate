package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm/coverage"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/helper"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/tf"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/types"
)

type PlanCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (c *PlanCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c PlanCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Printf("[INFO] initializing terraform...")
	workingDirectory, _ := os.Getwd()
	terraform, err := tf.NewTerraform(workingDirectory, c.verbose)
	if err != nil {
		log.Fatal(err)
	}
	c.Plan(terraform, true)
	return 0
}

func (c PlanCommand) Help() string {
	helpText := `
Usage: azurerm-restapi-to-azurerm plan
` + c.Synopsis() + "\nThe Terraform addresses listed in file `azurerm-restapi-to-azurerm.ignore` will be ignored during migration.\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c PlanCommand) Synopsis() string {
	return "Show azurerm-restapi resources which can migrate to azurerm resources in current working directory"
}

func (c PlanCommand) Plan(terraform *tf.Terraform, isPlanOnly bool) ([]types.GenericResource, []types.GenericPatchResource) {
	// get azurerm-restapi resource from state
	log.Printf("[INFO] searching azurerm-restapi_resource & azurerm-restapi_patch_resource...")
	p, err := terraform.Plan()
	if err != nil {
		log.Fatal(err)
	}

	migrationMessage := "The following resources will be migrated:\n"
	unsupportedMessage := "The following resources can't be migrated:\n"
	ignoreMessage := "The following resources will be ignored in migration:\n"
	ignoreSet := make(map[string]bool)
	if file, err := ioutil.ReadFile(path.Join(terraform.GetWorkingDirectory(), "azurerm-restapi-to-azurerm.ignore")); err == nil {
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
	patchResources := make([]types.GenericPatchResource, 0)
	for _, resource := range terraform.ListGenericResources(p) {
		if ignoreSet[resource.OldAddress(nil)] {
			continue
		}
		resourceId := ""
		for _, instance := range resource.Instances {
			resourceId = instance.ResourceId
			break
		}
		resourceTypes := azurerm.GetAzureRMResourceType(resourceId)

		idPattern, _ := helper.GetIdPattern(resourceId)
		_, uncoveredPut := coverage.GetPutCoverage(resource.InputProperties, idPattern)
		_, uncoveredGet := coverage.GetGetCoverage(resource.OutputProperties, idPattern)

		if len(uncoveredGet)+len(uncoveredPut) == 0 {
			if !isPlanOnly {
				if len(resourceTypes) == 1 {
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

	for _, resource := range terraform.ListGenericPatchResources(p) {
		if ignoreSet[resource.OldAddress()] {
			continue
		}
		resourceTypes := azurerm.GetAzureRMResourceType(resource.Id)

		idPattern, _ := helper.GetIdPattern(resource.Id)
		_, uncoveredPut := coverage.GetPutCoverage(resource.InputProperties, idPattern)
		_, uncoveredGet := coverage.GetGetCoverage(resource.OutputProperties, idPattern)

		if len(uncoveredGet)+len(uncoveredPut) == 0 {
			if !isPlanOnly {
				if len(resourceTypes) == 1 {
					resource.ResourceType = resourceTypes[0]
				} else {
					resource.ResourceType = c.getUserInputResourceType(resource.Id, resourceTypes)
				}
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(), resource.ResourceType)
			} else {
				migrationMessage += fmt.Sprintf("\t%s will be replaced with %v\n", resource.OldAddress(), strings.Join(resourceTypes, ", "))
			}
			patchResources = append(patchResources, resource)
		} else {
			unsupportedMessage += fmt.Sprintf("\t%s: input properties not supported: [%v], output properties not supported: [%v]\n",
				resource.OldAddress(), strings.Join(uncoveredPut, ", "), strings.Join(uncoveredGet, ", "))
		}
	}

	log.Printf("[INFO]\n\nThe tool will perform the following actions:\n\n%s\n%s\n%s\n", migrationMessage, unsupportedMessage, ignoreMessage)
	return resources, patchResources
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
