package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/tf"
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
	Plan(terraform)
	return 0
}

func (c PlanCommand) Help() string {
	helpText := `
Usage: azurerm-restapi-to-azurerm plan
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c PlanCommand) Synopsis() string {
	return "Show azurerm-restapi resources which can migrate to azurerm resources in current working directory"
}

func Plan(terraform *tf.Terraform) {
	// get azurerm-restapi resource from state
	log.Printf("[INFO] searching azurerm-restapi_resource...")
	genericResources, err := terraform.ListGenericResources()
	if err != nil {
		log.Fatal(err)
	}
	// get azurerm-restapi patch resource from state
	log.Printf("[INFO] searching azurerm-restapi_patch_resource...")
	patchResources, err := terraform.ListGenericPatchResources()
	if err != nil {
		log.Fatal(err)
	}

	// TODO: check whether resources can migrate based on coverage report

	// get migrated azurerm resource type

	log.Printf("[INFO] found %d azurerm-restapi_resource can migrate to azurerm resource", len(genericResources))
	for _, resource := range genericResources {
		resourceId := ""
		for _, instance := range resource.Instances {
			resourceId = instance.ResourceId
			break
		}
		resourceTypes := azurerm.GetAzureRMResourceType(resourceId)
		log.Printf("[INFO]resource %s will migrate to %v", resource.OldAddress(nil), resourceTypes)
	}

	// get migrated azurerm resource type
	log.Printf("[INFO] found %d azurerm-restapi_patch_resource can migrate to azurerm resource", len(patchResources))
	for _, resource := range patchResources {
		resourceTypes := azurerm.GetAzureRMResourceType(resource.Id)
		log.Printf("[INFO]resource %s will migrate to %v", resource.OldAddress(), resourceTypes)
	}
}
