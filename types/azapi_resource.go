package types

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Azure/aztfmigrate/azurerm"
	"github.com/Azure/aztfmigrate/azurerm/coverage"
	"github.com/Azure/aztfmigrate/azurerm/schema"
	"github.com/Azure/aztfmigrate/helper"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var _ AzureResource = &AzapiResource{}

type AzapiResource struct {
	Label            string
	Instances        []Instance
	ResourceType     string
	Block            *hclwrite.Block
	References       []Reference
	InputProperties  []string
	OutputProperties []string
	Migrated         bool
}

func (r *AzapiResource) Outputs() []Output {
	res := make([]Output, 0)
	for _, instance := range r.Instances {
		res = append(res, instance.Outputs...)
	}
	res = append(res, Output{
		OldName: r.OldAddress(nil),
		NewName: r.NewAddress(nil),
	})
	return res
}

func (r *AzapiResource) MigratedBlock() *hclwrite.Block {
	return r.Block
}

func (r *AzapiResource) IsMigrated() bool {
	return r.Migrated
}

func (r *AzapiResource) ImportBlock() *hclwrite.Block {
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
	return importBlock
}

func (r *AzapiResource) RemovedBlock() *hclwrite.Block {
	removedBlock := hclwrite.NewBlock("removed", nil)
	removedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: "azapi_resource"}, hcl.TraverseAttr{Name: r.Label}})
	removedLifecycleBlock := hclwrite.NewBlock("lifecycle", nil)
	removedLifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(false))
	removedBlock.Body().AppendBlock(removedLifecycleBlock)
	return removedBlock
}

func (r *AzapiResource) GenerateNewConfig(terraform *tf.Terraform) error {
	if !r.IsMultipleResources() {
		instance := r.Instances[0]
		block, err := importAndGenerateConfig(terraform, r.NewAddress(nil), instance.ResourceId, r.ResourceType, false)
		if err != nil {
			return err
		}
		r.Block = block
		valuePropMap := GetValuePropMap(r.Block, r.NewAddress(nil))
		for i, output := range r.Instances[0].Outputs {
			r.Instances[0].Outputs[i].NewName = valuePropMap[output.GetStringValue()]
		}
		for i, instance := range r.Instances {
			props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids", "id"}
			for _, prop := range props {
				r.Instances[i].Outputs = append(r.Instances[i].Outputs, Output{
					OldName: r.OldAddress(instance.Index) + "." + prop,
					NewName: r.NewAddress(instance.Index) + "." + prop,
				})
			}
		}
		r.Block = InjectReference(r.Block, r.References)
	} else {
		// import and build combined block
		blocks := make([]*hclwrite.Block, 0)
		for _, instance := range r.Instances {
			instanceAddress := fmt.Sprintf("%s.%s_%v", r.ResourceType, r.Label, strings.ReplaceAll(fmt.Sprintf("%v", instance.Index), "/", "_"))
			if block, err := importAndGenerateConfig(terraform, instanceAddress, instance.ResourceId, r.ResourceType, false); err == nil {
				blocks = append(blocks, block)
			}
		}
		combinedBlock := hclwrite.NewBlock("resource", []string{r.ResourceType, r.Label})
		if r.IsForEach() {
			foreachItems := CombineBlock(blocks, combinedBlock, true)
			foreachConfig := GetForEachConstants(r.Instances, foreachItems)
			combinedBlock.Body().SetAttributeRaw("for_each", helper.GetTokensForExpression(foreachConfig))
		} else {
			_ = CombineBlock(blocks, combinedBlock, false)
			combinedBlock.Body().SetAttributeRaw("count", helper.GetTokensForExpression(fmt.Sprintf("%d", len(r.Instances))))
		}

		r.Block = combinedBlock
		for i, instance := range r.Instances {
			valuePropMap := GetValuePropMap(blocks[i], r.NewAddress(instance.Index))
			for j, output := range r.Instances[i].Outputs {
				r.Instances[i].Outputs[j].NewName = valuePropMap[output.GetStringValue()]
			}
		}
		for i, instance := range r.Instances {
			r.Instances[i].Outputs = append(r.Instances[i].Outputs, Output{
				OldName: r.OldAddress(instance.Index),
				NewName: r.NewAddress(instance.Index),
			})
			props := []string{"location", "tags", "identity", "identity.0", "identity.0.type", "identity.0.identity_ids", "id"}
			for _, prop := range props {
				r.Instances[i].Outputs = append(r.Instances[i].Outputs, Output{
					OldName: r.OldAddress(instance.Index) + "." + prop,
					NewName: r.NewAddress(instance.Index) + "." + prop,
				})
			}
		}
		r.Block = InjectReference(r.Block, r.References)
	}
	r.Migrated = true
	return nil
}

func (r *AzapiResource) TargetProvider() string {
	return "azurerm"
}

func (r *AzapiResource) CoverageCheck(strictMode bool) error {
	if os.Getenv("AZTF_MIGRATE_SKIP_COVERAGE_CHECK") == "true" {
		return nil
	}
	resourceId := r.Instances[0].ResourceId
	idPattern, _ := GetIdPattern(resourceId)
	if strictMode {
		azurermApiVersion := coverage.GetApiVersion(idPattern)
		if azurermApiVersion != r.Instances[0].ApiVersion {
			return fmt.Errorf("%s: api-versions are not matched, expect %s, got %s",
				r.OldAddress(nil), r.Instances[0].ApiVersion, azurermApiVersion)
		}
	}

	_, uncoveredPut := coverage.GetPutCoverage(r.InputProperties, idPattern)
	_, uncoveredGet := coverage.GetGetCoverage(r.OutputProperties, idPattern)

	if len(uncoveredGet)+len(uncoveredPut) != 0 {
		return fmt.Errorf("%s: input properties not supported: [%v], output properties not supported: [%v]",
			r.OldAddress(nil), strings.Join(uncoveredPut, ", "), strings.Join(uncoveredGet, ", "))
	}
	return nil
}

func (r *AzapiResource) OldAddress(index interface{}) string {
	oldAddress := fmt.Sprintf("azapi_resource.%s", r.Label)
	if index == nil {
		return oldAddress
	}
	switch i := index.(type) {
	case int, int32, int64, float32, float64:
		return fmt.Sprintf(`%s[%v]`, oldAddress, i)
	case string:
		return fmt.Sprintf(`%s["%s"]`, oldAddress, i)
	default:
		return oldAddress
	}
}

func (r *AzapiResource) NewAddress(index interface{}) string {
	newAddress := fmt.Sprintf("%s.%s", r.ResourceType, r.Label)
	if index == nil {
		return newAddress
	}
	switch i := index.(type) {
	case int, int32, int64, float32, float64:
		return fmt.Sprintf(`%s[%v]`, newAddress, i)
	case string:
		return fmt.Sprintf(`%s["%s"]`, newAddress, i)
	default:
		return newAddress
	}
}

func (r *AzapiResource) EmptyImportConfig() string {
	config := ""
	for _, instance := range r.Instances {
		if !r.IsMultipleResources() {
			config += fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
		} else {
			config += fmt.Sprintf("resource \"%s\" \"%s_%s\" {}\n", r.ResourceType, r.Label, strings.ReplaceAll(fmt.Sprintf("%v", instance.Index), "/", "_"))
		}
	}
	return config
}

func (r *AzapiResource) IsMultipleResources() bool {
	return len(r.Instances) != 0 && r.Instances[0].Index != nil
}

func (r *AzapiResource) IsForEach() bool {
	if len(r.Instances) != 0 && r.Instances[0].Index != nil {
		if _, ok := r.Instances[0].Index.(string); ok {
			return true
		}
	}
	return false
}

type Instance struct {
	Index      interface{}
	ApiVersion string
	ResourceId string
	Outputs    []Output
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
