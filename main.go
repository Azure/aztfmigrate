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

const filenameImportConfig = "imports.tf"

func main() {
	migrateGenericResource()
	migrateGenericPatchResource()
}

func migrateGenericResource() {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azurerm-restapi_resource")
	log.Printf("[INFO] initializing terraform")
	workingDirectory, _ := os.Getwd()
	terraform, err := tf.NewTerraform(workingDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// get azurerm-restapi resource from state
	resources, err := terraform.ListGenericResources()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] found %d azurerm-restapi_resource in state\n", len(resources))

	// get migrated azurerm resource type
	for index, resource := range resources {
		resources[index].ResourceType = azurerm.GetAzureRMResourceType(resource.Id)
	}
	log.Printf("[INFO] found %d azurerm-restapi_resource can migrate to azurerm resource", len(resources))
	for _, r := range resources {
		log.Printf("[INFO] resource %s will migrate to resource %s", r.OldAddress(), r.NewAddress())
	}

	// generate import config
	config := ""
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}
	err = ioutil.WriteFile(filenameImportConfig, []byte(config), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
		if block, err := importAndGenerateConfig(terraform, r.NewAddress(), r.Id, r.ResourceType); err == nil {
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
	if err := os.Remove(filenameImportConfig); err != nil {
		log.Fatal(err)
	}
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress())
			if err := helper.ReplaceResourceBlock(r.OldAddress(), r.Block); err != nil {
				log.Printf("[ERROR] error removing %s from state: %+v", r.OldAddress(), err)
			}
		}
	}

	// replace jsondecode(xxx.output) with azurerm references
	log.Printf("[INFO] replacing azurerm-restapi resource references with azurerm resoure reference.")
	outputs := make([]types.Output, 0)
	for _, r := range resources {
		if r.Migrated {
			outputs = append(outputs, r.Outputs...)
		}
	}
	if err := helper.ReplaceGenericOutputs(outputs); err != nil {
		log.Printf("[ERROR] replacing azurerm-restapi resource references with azurerm resoure reference: %+v", err)
	}
}

func migrateGenericPatchResource() {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azurerm-restapi_patch_resource")
	log.Printf("[INFO] initializing terraform")
	workingDirectory, _ := os.Getwd()
	terraform, err := tf.NewTerraform(workingDirectory)
	if err != nil {
		log.Fatal(err)
	}

	// get azurerm-restapi patch resource from state
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
	for _, r := range resources {
		log.Printf("[INFO] resource %s will migrate to resource %s", r.OldAddress(), r.NewAddress())
	}

	// generate import config
	config := `
provider "azurerm" {
  features {}
}
`
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}

	// save empty import config to temp dir
	tempPath := filepath.Join(workingDirectory, "temp")
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(tempPath); err != nil {
			log.Fatalf("error deleting %s: %+v", tempPath, err)
		}
	}
	if err := os.MkdirAll(tempPath, 0755); err != nil {
		log.Fatalf("creating temp workspace %q: %w", tempPath, err)
	}
	tempTerraform, err := tf.NewTerraform(tempPath)
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile(filepath.Join(tempPath, "main.tf"), []byte(config), 0644); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	for index, r := range resources {
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

	log.Printf("[INFO] replacing azurerm-restapi resource references with azurerm resoure reference.")
	outputs := make([]types.Output, 0)
	for _, r := range resources {
		if r.Migrated {
			outputs = append(outputs, r.Outputs...)
		}
	}
	if err := helper.ReplaceGenericOutputs(outputs); err != nil {
		log.Printf("[ERROR] replacing azurerm-restapi resource references with azurerm resoure reference: %+v", err)
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

	// cleanup temp folder
	if _, err := os.Stat(tempPath); !os.IsNotExist(err) {
		if err := os.RemoveAll(tempPath); err != nil {
			log.Printf("[WARN] error deleting %s: %+v", tempPath, err)
		}
	}
}

func importAndGenerateConfig(terraform *tf.Terraform, address string, id string, resourceType string) (*hclwrite.Block, error) {
	tpl, err := terraform.Import(address, id)
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
