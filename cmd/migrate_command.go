package cmd

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/azurerm/schema"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/helper"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/tf"
	"github.com/ms-henglu/azurerm-restapi-to-azurerm/types"
)

const filenameImport = "imports.tf"
const providerConfig = `
provider "azurerm" {
  features {}
}
`
const tempDir = "temp"

type MigrateCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (c *MigrateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c MigrateCommand) Run(args []string) int {
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
	c.MigrateGenericResource(terraform, workingDirectory)
	c.MigrateGenericPatchResource(terraform, workingDirectory)
	return 0
}

func (c MigrateCommand) Help() string {
	helpText := `
Usage: azurerm-restapi-to-azurerm migrate
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c MigrateCommand) Synopsis() string {
	return "Migrate azurerm-restapi resources to azurerm resources in current working directory"
}

func (c MigrateCommand) MigrateGenericResource(terraform *tf.Terraform, workingDirectory string) {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azurerm-restapi_resource")

	// get azurerm-restapi resource from state
	log.Printf("[INFO] searching azurerm-restapi_resource...")
	resources, err := terraform.ListGenericResources()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] found %d azurerm-restapi_resource in state\n", len(resources))

	// get migrated azurerm resource type
	for index, resource := range resources {
		resourceId := ""
		for _, instance := range resource.Instances {
			resourceId = instance.ResourceId
			break
		}
		resourceTypes := azurerm.GetAzureRMResourceType(resourceId)
		if len(resourceTypes) == 1 {
			resources[index].ResourceType = resourceTypes[0]
			continue
		}
		log.Printf("[WARN] couldn't find unique resource type for id: %s\npossible values are %s.\nPlease input a azurerm resource type", resourceId, strings.Join(resourceTypes, ", "))
		reader := bufio.NewReader(os.Stdin)
		resourceType, _ := reader.ReadString('\n')
		resources[index].ResourceType = strings.Trim(resourceType, "\r\n")
	}
	log.Printf("[INFO] found %d azurerm-restapi_resource can migrate to azurerm resource", len(resources))

	// generate import config
	config := ""
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}
	if err := ioutil.WriteFile(filepath.Join(workingDirectory, filenameImport), []byte(config), 0644); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
		log.Printf("[INFO] migrating resource %s (%d instances) to resource %s...", r.OldAddress(nil), len(r.Instances), r.NewAddress(nil))
		if !r.IsMultipleResources() {
			instance := r.Instances[0]
			log.Printf("[INFO] importing %s to %s and generating config...", instance.ResourceId, r.NewAddress(nil))
			if block, err := importAndGenerateConfig(terraform, r.NewAddress(nil), instance.ResourceId, r.ResourceType, false); err == nil {
				resources[index].Block = block
				valuePropMap := helper.GetValuePropMap(resources[index].Block, resources[index].NewAddress(nil))
				for i, output := range resources[index].Instances[0].Outputs {
					resources[index].Instances[0].Outputs[i].NewName = valuePropMap[output.GetStringValue()]
				}
				for i, instance := range r.Instances {
					resources[index].Instances[i].Outputs = append(resources[index].Instances[i].Outputs, types.Output{
						OldName: resources[index].OldAddress(instance.Index) + ".resource_id",
						NewName: resources[index].NewAddress(instance.Index) + ".id",
					})
					props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids"}
					for _, prop := range props {
						resources[index].Instances[i].Outputs = append(resources[index].Instances[i].Outputs, types.Output{
							OldName: resources[index].OldAddress(instance.Index) + "." + prop,
							NewName: resources[index].NewAddress(instance.Index) + "." + prop,
						})
					}
				}
				resources[index].Block = helper.InjectReference(resources[index].Block, resources[index].References)
				resources[index].Migrated = true
				log.Printf("[INFO] resource %s has migrated to %s", r.OldAddress(nil), r.NewAddress(nil))
			} else {
				log.Printf("[ERROR] %+v", err)
			}
		} else {
			// import into real state
			for _, instance := range r.Instances {
				address := r.NewAddress(instance.Index)
				log.Printf("[INFO] importing %s to %s...", instance.ResourceId, address)
				if err := terraform.Import(address, instance.ResourceId); err != nil {
					log.Printf("[Error] error importing %s : %s", address, instance.ResourceId)
				}
			}

			// write empty config to temp dir for import
			tempDirectoryCreate(workingDirectory)
			tempPath := filepath.Join(workingDirectory, tempDir)
			tempTerraform, err := tf.NewTerraform(tempPath, c.verbose)
			if err != nil {
				log.Fatal(err)
			}
			config := providerConfig
			for _, instance := range r.Instances {
				config += fmt.Sprintf("resource \"%s\" \"%s_%v\" {}\n", r.ResourceType, r.Label, instance.Index)
			}
			if err = ioutil.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0644); err != nil {
				log.Fatal(err)
			}

			// import and build combined block
			log.Printf("[INFO] generating config...")
			blocks := make([]*hclwrite.Block, 0)
			for _, instance := range r.Instances {
				if block, err := importAndGenerateConfig(tempTerraform, fmt.Sprintf("%s.%s_%v", r.ResourceType, r.Label, instance.Index), instance.ResourceId, r.ResourceType, false); err == nil {
					blocks = append(blocks, block)
				}
			}
			combinedBlock := hclwrite.NewBlock("resource", []string{r.ResourceType, r.Label})
			if r.IsForEach() {
				foreachItems := helper.CombineBlock(blocks, combinedBlock, true)
				foreachConfig := helper.GetForEachConstants(r.Instances, foreachItems)
				combinedBlock.Body().SetAttributeRaw("for_each", helper.GetTokensForExpression(foreachConfig))
			} else {
				_ = helper.CombineBlock(blocks, combinedBlock, false)
				combinedBlock.Body().SetAttributeRaw("count", helper.GetTokensForExpression(fmt.Sprintf("%d", len(r.Instances))))
			}

			resources[index].Block = combinedBlock
			for i, instance := range r.Instances {
				valuePropMap := helper.GetValuePropMap(blocks[i], resources[index].NewAddress(instance.Index))
				for j, output := range resources[index].Instances[i].Outputs {
					resources[index].Instances[i].Outputs[j].NewName = valuePropMap[output.GetStringValue()]
				}
			}
			for i, instance := range r.Instances {
				resources[index].Instances[i].Outputs = append(resources[index].Instances[i].Outputs, types.Output{
					OldName: resources[index].OldAddress(instance.Index) + ".resource_id",
					NewName: resources[index].NewAddress(instance.Index) + ".id",
				})
				resources[index].Instances[i].Outputs = append(resources[index].Instances[i].Outputs, types.Output{
					OldName: resources[index].OldAddress(instance.Index),
					NewName: resources[index].NewAddress(instance.Index),
				})
				props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids"}
				for _, prop := range props {
					resources[index].Instances[i].Outputs = append(resources[index].Instances[i].Outputs, types.Output{
						OldName: resources[index].OldAddress(instance.Index) + "." + prop,
						NewName: resources[index].NewAddress(instance.Index) + "." + prop,
					})
				}
			}
			resources[index].Block = helper.InjectReference(resources[index].Block, resources[index].References)
			resources[index].Migrated = true
			log.Printf("[INFO] resource %s has migrated to %s", r.OldAddress(nil), r.NewAddress(nil))
		}
	}
	tempDirectoryCleanup(workingDirectory)

	// remove from state
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from state", r.OldAddress(nil))
			exec := terraform.GetExec()
			if err := exec.StateRm(context.TODO(), r.OldAddress(nil)); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(nil), err)
			}
		}
	}

	// migrate depends_on, lifecycle, provisioner
	for index, r := range resources {
		if existingBlock, err := helper.GetResourceBlock(workingDirectory, r.OldAddress(nil)); err == nil && existingBlock != nil {
			if attr := existingBlock.Body().GetAttribute("depends_on"); attr != nil {
				resources[index].Block.Body().SetAttributeRaw("depends_on", attr.Expr().BuildTokens(nil))
			}
			for _, block := range existingBlock.Body().Blocks() {
				if block.Type() == "lifecycle" || block.Type() == "provisioner" {
					resources[index].Block.Body().AppendBlock(block)
				}
			}
		}
	}

	// remove from config
	if err := os.Remove(filepath.Join(workingDirectory, filenameImport)); err != nil {
		log.Fatal(err)
	}
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress(nil))
			if err := helper.ReplaceResourceBlock(workingDirectory, r.OldAddress(nil), r.Block); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(nil), err)
			}
		}
	}

	// replace jsondecode(xxx.output) with azurerm references
	outputs := make([]types.Output, 0)
	for _, r := range resources {
		if r.Migrated {
			for _, instance := range r.Instances {
				outputs = append(outputs, instance.Outputs...)
			}
			outputs = append(outputs, types.Output{
				OldName: r.OldAddress(nil),
				NewName: r.NewAddress(nil),
			})
		}
	}
	if len(outputs) != 0 {
		log.Printf("[INFO] replacing azurerm-restapi resource references with azurerm resoure reference.")
		if err := helper.ReplaceGenericOutputs(workingDirectory, outputs); err != nil {
			log.Printf("[ERROR] replacing azurerm-restapi resource references with azurerm resoure reference: %+v", err)
		}
	}
}

func (c MigrateCommand) MigrateGenericPatchResource(terraform *tf.Terraform, workingDirectory string) {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azurerm-restapi_patch_resource")
	log.Printf("[INFO] initializing terraform")

	// get azurerm-restapi patch resource from state
	log.Printf("[INFO] searching azurerm-restapi_patch_resource...")
	resources, err := terraform.ListGenericPatchResources()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] found %d azurerm-restapi_patch_resource\n", len(resources))

	// get migrated azurerm resource type
	for index, resource := range resources {
		resourceTypes := azurerm.GetAzureRMResourceType(resource.Id)
		if len(resourceTypes) == 1 {
			resources[index].ResourceType = resourceTypes[0]
			continue
		}
		log.Printf("[WARN] couldn't find unique resource type for id: %s\npossible values are %s.\nPlease input a azurerm resource type", resource.Id, strings.Join(resourceTypes, ", "))
		reader := bufio.NewReader(os.Stdin)
		resourceType, _ := reader.ReadString('\n')
		resources[index].ResourceType = strings.Trim(resourceType, "\r\n")
	}
	log.Printf("[INFO] found %d azurerm-restapi_patch_resource can migrate to azurerm resource", len(resources))

	// generate import config
	config := providerConfig
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}

	// save empty import config to temp dir
	tempDirectoryCreate(workingDirectory)
	tempPath := filepath.Join(workingDirectory, tempDir)
	tempTerraform, err := tf.NewTerraform(tempPath, c.verbose)
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0644); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
		log.Printf("[INFO] migrating resource %s to resource %s", r.OldAddress(), r.NewAddress())
		if block, err := importAndGenerateConfig(tempTerraform, r.NewAddress(), r.Id, r.ResourceType, true); err == nil {
			resources[index].Block = block
			valuePropMap := helper.GetValuePropMap(resources[index].Block, resources[index].NewAddress())
			for i := range resources[index].Outputs {
				resources[index].Outputs[i].NewName = valuePropMap[resources[index].Outputs[i].GetStringValue()]
			}
			resources[index].Block = helper.InjectReference(resources[index].Block, resources[index].References)
			resources[index].Migrated = true
			log.Printf("[INFO] %s has migrated to %s", r.OldAddress(), r.NewAddress())
		} else {
			log.Printf("[ERROR] %+v", err)
		}
	}

	if err := helper.UpdateMigratedResourceBlock(workingDirectory, resources); err != nil {
		log.Fatal(err)
	}

	outputs := make([]types.Output, 0)
	for _, r := range resources {
		if r.Migrated {
			outputs = append(outputs, r.Outputs...)
		}
	}
	if len(outputs) != 0 {
		log.Printf("[INFO] replacing azurerm-restapi resource references with azurerm resoure reference.")
		if err := helper.ReplaceGenericOutputs(workingDirectory, outputs); err != nil {
			log.Printf("[ERROR] replacing azurerm-restapi resource references with azurerm resoure reference: %+v", err)
		}
	}

	// remove from state
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from state", r.OldAddress())
			exec := terraform.GetExec()
			if err := exec.StateRm(context.TODO(), r.OldAddress()); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(), err)
			}
		}
	}

	// remove from config
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress())
			if err := helper.ReplaceResourceBlock(workingDirectory, r.OldAddress(), nil); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(), err)
			}
		}
	}

	tempDirectoryCleanup(workingDirectory)
}

func importAndGenerateConfig(terraform *tf.Terraform, address string, id string, resourceType string, skipTune bool) (*hclwrite.Block, error) {
	tpl, err := terraform.ImportAdd(address, id)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error importing address: %s, id: %s: %+v", address, id, err)
	}
	f, diag := hclwrite.ParseConfig([]byte(tpl), "", hcl.InitialPos)
	if (diag != nil && diag.HasErrors()) || f == nil {
		return nil, fmt.Errorf("[ERROR] parsing the HCL generated by \"terraform add\" of %s: %s", address, diag.Error())
	}

	if !skipTune {
		rb := f.Body().Blocks()[0].Body()
		sch := schema.ProviderSchemaInfo.ResourceSchemas[resourceType]
		if err := azurerm.TuneHCLSchemaForResource(rb, sch); err != nil {
			return nil, fmt.Errorf("[ERROR] tuning hcl config base on schema: %+v", err)
		}
	}

	return f.Body().Blocks()[0], nil
}

func tempDirectoryCreate(workingDirectory string) {
	tempPath := filepath.Join(workingDirectory, tempDir)
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(tempPath); err != nil {
			log.Fatalf("error deleting %s: %+v", tempPath, err)
		}
	}
	if err := os.MkdirAll(tempPath, 0755); err != nil {
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
