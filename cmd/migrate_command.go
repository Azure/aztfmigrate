package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Azure/azapi2azurerm/azurerm"
	"github.com/Azure/azapi2azurerm/azurerm/schema"
	"github.com/Azure/azapi2azurerm/helper"
	"github.com/Azure/azapi2azurerm/tf"
	"github.com/Azure/azapi2azurerm/types"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"github.com/zclconf/go-cty/cty"
)

const filenameImport = "imports.tf"
const providerConfig = `
provider "azurerm" {
  features {}
  subscription_id = "%s"
}
`
const tempDir = "temp"

type MigrateCommand struct {
	Ui         cli.Ui
	Verbose    bool
	Strict     bool
	workingDir string
}

func (c *MigrateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("plan")
	fs.BoolVar(&c.Verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.Strict, "strict", false, "strict mode: API versions must be matched")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")

	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c MigrateCommand) Run(args []string) int {
	// AzureRM provider will honor env.var "AZURE_HTTP_USER_AGENT" when constructing for HTTP "User-Agent" header.
	// #nosec G104
	os.Setenv("AZURE_HTTP_USER_AGENT", "mig")
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
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
	resources, updateResources := PlanCommand{ //nolint
		Ui:      c.Ui,
		Verbose: c.Verbose,
		Strict:  c.Strict,
	}.Plan(terraform, false)
	c.MigrateGenericResource(terraform, resources)
	c.MigrateGenericUpdateResource(terraform, updateResources)
	return 0
}

func (c MigrateCommand) Help() string {
	helpText := `
Usage: azapi2azurerm migrate
` + c.Synopsis() + "\nThe Terraform addresses listed in file `azapi2azurerm.ignore` will be ignored during migration.\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c MigrateCommand) Synopsis() string {
	return "Migrate azapi resources to azurerm resources in current working directory"
}

func (c MigrateCommand) MigrateGenericResource(terraform *tf.Terraform, resources []types.GenericResource) {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azapi_resource")

	workingDirectory := terraform.GetWorkingDirectory()

	// import and generate config
	for index, r := range resources {
		log.Printf("[INFO] migrating resource %s (%d instances) to resource %s...", r.OldAddress(nil), len(r.Instances), r.NewAddress(nil))

		// write empty config to temp dir for import
		tempDirectoryCreate(workingDirectory)
		tempPath := filepath.Join(workingDirectory, tempDir)
		tempTerraform, err := tf.NewTerraform(tempPath, c.Verbose)
		if err != nil {
			log.Fatal(err)
		}

		subscriptionId := ""
		for _, instance := range r.Instances {
			if strings.HasPrefix(instance.ResourceId, "/subscriptions/") {
				subscriptionId = strings.Split(instance.ResourceId, "/")[2]
				break
			}
		}
		config := fmt.Sprintf(providerConfig, subscriptionId)
		for _, instance := range r.Instances {
			if !r.IsMultipleResources() {
				config += fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
			} else {
				config += fmt.Sprintf("resource \"%s\" \"%s_%v\" {}\n", r.ResourceType, r.Label, instance.Index)
			}
		}
		if err = os.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0600); err != nil {
			log.Fatal(err)
		}

		if !r.IsMultipleResources() {
			instance := r.Instances[0]
			log.Printf("[INFO] importing %s to %s and generating config...", instance.ResourceId, r.NewAddress(nil))
			if block, err := importAndGenerateConfig(tempTerraform, r.NewAddress(nil), instance.ResourceId, r.ResourceType, false); err == nil {
				resources[index].Block = block
				valuePropMap := helper.GetValuePropMap(resources[index].Block, resources[index].NewAddress(nil))
				for i, output := range resources[index].Instances[0].Outputs {
					resources[index].Instances[0].Outputs[i].NewName = valuePropMap[output.GetStringValue()]
				}
				for i, instance := range r.Instances {
					props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids", "id"}
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
					OldName: resources[index].OldAddress(instance.Index),
					NewName: resources[index].NewAddress(instance.Index),
				})
				props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids", "id"}
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
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress(nil))

			importBlock := hclwrite.NewBlock("import", nil)
			if r.IsMultipleResources() {
				forEachMap := make(map[string]cty.Value)
				for _, instance := range r.Instances {
					switch v := instance.Index.(type) {
					case string:
						forEachMap[instance.ResourceId] = cty.StringVal(v)
					default:
						value, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
						forEachMap[instance.ResourceId] = cty.NumberIntVal(value)
					}
				}
				importBlock.Body().SetAttributeValue("for_each", cty.MapVal(forEachMap))
				importBlock.Body().SetAttributeTraversal("id", hcl.Traversal{hcl.TraverseRoot{Name: "each"}, hcl.TraverseAttr{Name: "key"}})
				importBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: r.ResourceType}, hcl.TraverseAttr{Name: fmt.Sprintf("%s[each.value]", r.Label)}})
			} else {
				importBlock.Body().SetAttributeValue("id", cty.StringVal(r.Instances[0].ResourceId))
				importBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: r.ResourceType}, hcl.TraverseAttr{Name: r.Label}})
			}

			removedBlock := hclwrite.NewBlock("removed", nil)
			removedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: "azapi_resource"}, hcl.TraverseAttr{Name: r.Label}})
			removedLifecycleBlock := hclwrite.NewBlock("lifecycle", nil)
			removedLifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(false))
			removedBlock.Body().AppendBlock(removedLifecycleBlock)

			if err := helper.ReplaceResourceBlock(workingDirectory, r.OldAddress(nil), []*hclwrite.Block{removedBlock, importBlock, r.Block}); err != nil {
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
		log.Printf("[INFO] replacing azapi resource references with azurerm resoure reference.")
		if err := helper.ReplaceGenericOutputs(workingDirectory, outputs); err != nil {
			log.Printf("[ERROR] replacing azapi resource references with azurerm resoure reference: %+v", err)
		}
	}
}

func (c MigrateCommand) MigrateGenericUpdateResource(terraform *tf.Terraform, resources []types.GenericUpdateResource) {
	log.Printf("[INFO] -----------------------------------------------")
	log.Printf("[INFO] task: migrate azapi_update_resource")

	// generate import config
	subscriptionId := ""
	for _, instance := range resources {
		if strings.HasPrefix(instance.Id, "/subscriptions/") {
			subscriptionId = strings.Split(instance.Id, "/")[2]
			break
		}
	}
	config := fmt.Sprintf(providerConfig, subscriptionId)
	for _, resource := range resources {
		config += resource.EmptyImportConfig()
	}

	// save empty import config to temp dir
	workingDirectory := terraform.GetWorkingDirectory()
	tempDirectoryCreate(workingDirectory)
	tempPath := filepath.Join(workingDirectory, tempDir)
	tempTerraform, err := tf.NewTerraform(tempPath, c.Verbose)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(tempPath, filenameImport), []byte(config), 0600); err != nil {
		log.Fatal(err)
	}

	// import and generate config
	newAddrs := make([]string, 0)
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
			newAddrs = append(newAddrs, r.NewAddress())
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
		log.Printf("[INFO] replacing azapi resource references with azurerm resoure reference.")
		if err := helper.ReplaceGenericOutputs(workingDirectory, outputs); err != nil {
			log.Printf("[ERROR] replacing azapi resource references with azurerm resoure reference: %+v", err)
		}
	}

	// remove from config
	for _, r := range resources {
		if r.Migrated {
			log.Printf("[INFO] removing %s from config", r.OldAddress())
			removedBlock := hclwrite.NewBlock("removed", nil)
			removedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: "azapi_update_resource"}, hcl.TraverseAttr{Name: r.OldLabel}})
			removedLifecycleBlock := hclwrite.NewBlock("lifecycle", nil)
			removedLifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(false))
			removedBlock.Body().AppendBlock(removedLifecycleBlock)

			if err := helper.ReplaceResourceBlock(workingDirectory, r.OldAddress(), []*hclwrite.Block{removedBlock}); err != nil {
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
