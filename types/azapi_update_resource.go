package types

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/aztfmigrate/azurerm/coverage"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

var _ AzureResource = &AzapiUpdateResource{}

type AzapiUpdateResource struct {
	ApiVersion       string
	Label            string
	OldLabel         string
	Id               string
	ResourceType     string
	Change           *tfjson.Change
	Block            *hclwrite.Block
	Migrated         bool
	References       []Reference
	outputs          []Output
	InputProperties  []string
	OutputProperties []string
}

func (r *AzapiUpdateResource) StateUpdateBlocks() []*hclwrite.Block {
	blocks := make([]*hclwrite.Block, 0)
	blocks = append(blocks, r.removedBlock())
	return blocks
}

func (r *AzapiUpdateResource) Outputs() []Output {
	res := make([]Output, 0)
	res = append(res, r.outputs...)
	res = append(res, Output{
		OldName: r.OldAddress(nil),
		NewName: r.NewAddress(nil),
	})
	return res
}

func (r *AzapiUpdateResource) MigratedBlock() *hclwrite.Block {
	return nil
}

func (r *AzapiUpdateResource) IsMigrated() bool {
	return r.Migrated
}

func (r *AzapiUpdateResource) GenerateNewConfig(terraform *tf.Terraform) error {
	block, err := importAndGenerateConfig(terraform, r.NewAddress(nil), r.Id, r.ResourceType, true)
	if err != nil {
		return err
	}
	r.Block = block
	valuePropMap := GetValuePropMap(r.Block, r.NewAddress(nil))
	for i := range r.outputs {
		r.outputs[i].NewName = valuePropMap[r.outputs[i].GetStringValue()]
	}
	r.Block = InjectReference(r.Block, r.References)
	r.Migrated = true
	return nil
}

func (r *AzapiUpdateResource) removedBlock() *hclwrite.Block {
	removedBlock := hclwrite.NewBlock("removed", nil)
	removedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: "azapi_update_resource"}, hcl.TraverseAttr{Name: r.OldLabel}})
	removedLifecycleBlock := hclwrite.NewBlock("lifecycle", nil)
	removedLifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(false))
	removedBlock.Body().AppendBlock(removedLifecycleBlock)
	return removedBlock
}

func (r *AzapiUpdateResource) TargetProvider() string {
	return "azurerm"
}

func (r *AzapiUpdateResource) CoverageCheck(strictMode bool) error {
	if os.Getenv("AZTF_MIGRATE_SKIP_COVERAGE_CHECK") == "true" {
		return nil
	}
	idPattern, _ := GetIdPattern(r.Id)
	if strictMode {
		azurermApiVersion := coverage.GetApiVersion(idPattern)
		if azurermApiVersion != r.ApiVersion {
			return fmt.Errorf("%s: api-versions are not matched, expect %s, got %s",
				r.OldAddress(nil), r.ApiVersion, azurermApiVersion)
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

func (r *AzapiUpdateResource) OldAddress(_ interface{}) string {
	return fmt.Sprintf("azapi_update_resource.%s", r.OldLabel)
}

func (r *AzapiUpdateResource) NewAddress(_ interface{}) string {
	return fmt.Sprintf("%s.%s", r.ResourceType, r.Label)
}

func (r *AzapiUpdateResource) EmptyImportConfig() string {
	return fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
}
