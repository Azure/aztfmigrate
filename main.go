package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
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

func main() {
	log.Printf("[INFO] initializing terraform...")
	workingDirectory, _ := os.Getwd()
	terraform, err := tf.NewTerraform(workingDirectory)
	if err != nil {
		log.Fatal(err)
	}
	migrateGenericResource(terraform, workingDirectory)
	migrateGenericPatchResource(terraform, workingDirectory)
}

func migrateGenericResource(terraform *tf.Terraform, workingDirectory string) {
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
		resources[index].ResourceType = azurerm.GetAzureRMResourceType(resourceId)
	}
	log.Printf("[INFO] found %d azurerm-restapi_resource can migrate to azurerm resource", len(resources))

	// generate import config
	config := ""
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}
	if err := ioutil.WriteFile(filenameImport, []byte(config), 0644); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
		log.Printf("[INFO] migrating resource %s (%d instances) to resource %s...", r.OldAddress(nil), len(r.Instances), r.NewAddress(nil))
		if !r.IsMultipleResources() {
			instance := r.Instances[0]
			log.Printf("[INFO] importing %s to %s and generating config...", r.NewAddress(nil), instance.ResourceId)
			if block, err := importAndGenerateConfig(terraform, r.NewAddress(nil), instance.ResourceId, r.ResourceType); err == nil {
				resources[index].Block = block
				valuePropMap := helper.GetValuePropMap(resources[index].Block, resources[index].NewAddress(nil))
				for i, output := range resources[index].Instances[0].Outputs {
					resources[index].Instances[0].Outputs[i].NewName = valuePropMap[output.GetStringValue()]
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
				log.Printf("[INFO] importing %s to %s...", address, instance.ResourceId)
				if err := terraform.Import(address, instance.ResourceId); err != nil {
					log.Printf("[Error] error importing %s : %s", address, instance.ResourceId)
				}
			}

			// write empty config to temp dir for import
			tempDirectoryCreate(workingDirectory)
			tempPath := filepath.Join(workingDirectory, tempDir)
			tempTerraform, err := tf.NewTerraform(tempPath)
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
				if block, err := importAndGenerateConfig(tempTerraform, fmt.Sprintf("%s.%s_%v", r.ResourceType, r.Label, instance.Index), instance.ResourceId, r.ResourceType); err == nil {
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
				valuePropMap := helper.GetValuePropMap(resources[index].Block, resources[index].NewAddress(instance.Index))
				for j, output := range resources[index].Instances[i].Outputs {
					resources[index].Instances[i].Outputs[j].NewName = valuePropMap[output.GetStringValue()]
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

	// remove from config
	if err := os.Remove(filenameImport); err != nil {
		log.Fatal(err)
	}
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress(nil))
			if err := helper.ReplaceResourceBlock(r.OldAddress(nil), r.Block); err != nil {
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
		}
	}
	if len(outputs) != 0 {
		log.Printf("[INFO] replacing azurerm-restapi resource references with azurerm resoure reference.")
		if err := helper.ReplaceGenericOutputs(outputs); err != nil {
			log.Printf("[ERROR] replacing azurerm-restapi resource references with azurerm resoure reference: %+v", err)
		}
	}
}

func migrateGenericPatchResource(terraform *tf.Terraform, workingDirectory string) {
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
		resources[index].ResourceType = azurerm.GetAzureRMResourceType(resource.Id)
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
	tempTerraform, err := tf.NewTerraform(tempPath)
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0644); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
		log.Printf("[INFO] migrating resource %s to resource %s", r.OldAddress(), r.NewAddress())
		if block, err := importAndGenerateConfig(tempTerraform, r.NewAddress(), r.Id, r.ResourceType); err == nil {
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

	if err := helper.UpdateMigratedResourceBlock(resources); err != nil {
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
		if err := helper.ReplaceGenericOutputs(outputs); err != nil {
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
			if err := helper.ReplaceResourceBlock(r.OldAddress(), nil); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(), err)
			}
		}
	}

	tempDirectoryCleanup(workingDirectory)
}

func importAndGenerateConfig(terraform *tf.Terraform, address string, id string, resourceType string) (*hclwrite.Block, error) {
	tpl, err := terraform.ImportAdd(address, id)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error importing address: %s, id: %s: %+v", address, id, err)
	}
	f, diag := hclwrite.ParseConfig([]byte(tpl), "", hcl.InitialPos)
	if (diag != nil && diag.HasErrors()) || f == nil {
		return nil, fmt.Errorf("[ERROR] parsing the HCL generated by \"terraform add\" of %s: %s", address, diag.Error())
	}

	rb := f.Body().Blocks()[0].Body()
	sch := schema.ProviderSchemaInfo.ResourceSchemas[resourceType]
	if err := azurerm.TuneHCLSchemaForResource(rb, sch); err != nil {
		return nil, fmt.Errorf("[ERROR] tuning hcl config base on schema: %+v", err)
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
		log.Fatalf("creating temp workspace %q: %w", tempPath, err)
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
