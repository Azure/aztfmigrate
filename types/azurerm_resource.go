package types

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/aztfmigrate/helper"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/gertd/go-pluralize"
	_ "github.com/gertd/go-pluralize"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

var _ AzureResource = &AzurermResource{}

type AzurermResource struct {
	OldLabel        string
	NewLabel        string
	OldResourceType string
	NewResourceType string
	Block           *hclwrite.Block
	Instances       []Instance
	References      []Reference
	Migrated        bool
}

func (r *AzurermResource) StateUpdateBlocks() []*hclwrite.Block {
	movedBlock := hclwrite.NewBlock("moved", nil)
	movedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: r.OldResourceType}, hcl.TraverseAttr{Name: r.OldLabel}})
	movedBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: r.NewResourceType}, hcl.TraverseAttr{Name: r.NewLabel}})
	blocks := make([]*hclwrite.Block, 0)
	blocks = append(blocks, movedBlock)
	return blocks
}

func (r *AzurermResource) Outputs() []Output {
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

func (r *AzurermResource) MigratedBlock() *hclwrite.Block {
	return r.Block
}

func (r *AzurermResource) IsMigrated() bool {
	return r.Migrated
}

func (r *AzurermResource) GenerateNewConfig(terraform *tf.Terraform) error {
	if !r.IsMultipleResources() {
		instance := r.Instances[0]
		log.Printf("[INFO] importing %s to %s and generating config...", instance.ResourceId, r.NewAddress(nil))
		if block, err := importAndGenerateConfig(terraform, r.NewAddress(nil), instance.ResourceId, "", true); err == nil {
			r.Block = block
			valuePropMap := GetValuePropMap(r.Block, r.NewAddress(nil))
			for i, output := range r.Instances[0].Outputs {
				r.Instances[0].Outputs[i].NewName = valuePropMap[output.GetStringValue()]
			}
			r.Migrated = true
			log.Printf("[INFO] resource %s has migrated to %s", r.OldAddress(nil), r.NewAddress(nil))
		} else {
			log.Printf("[ERROR] %+v", err)
		}
		r.Block = InjectReference(r.Block, r.References)
	} else {
		// import and build combined block
		log.Printf("[INFO] generating config...")
		blocks := make([]*hclwrite.Block, 0)
		for _, instance := range r.Instances {
			instanceAddress := fmt.Sprintf("%s.%s_%v", r.NewResourceType, r.NewLabel, strings.ReplaceAll(fmt.Sprintf("%v", instance.Index), "/", "_"))
			if block, err := importAndGenerateConfig(terraform, instanceAddress, instance.ResourceId, "", true); err == nil {
				blocks = append(blocks, block)
			}
		}
		combinedBlock := hclwrite.NewBlock("resource", []string{r.NewResourceType, r.NewLabel})
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
			// TODO: improve this, azapi resource should use .output.xxx to access the output properties
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

	r.Block = sortAttributes(r.Block)
	r.Migrated = true
	return nil
}

func (r *AzurermResource) TargetProvider() string {
	return "azapi"
}

func (r *AzurermResource) CoverageCheck(_ bool) error {
	// all properties are supported by azapi resource, so no need to check
	return nil
}

func (r *AzurermResource) OldAddress(index interface{}) string {
	oldAddress := fmt.Sprintf("%s.%s", r.OldResourceType, r.OldLabel)
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

func (r *AzurermResource) NewAddress(index interface{}) string {
	newAddress := fmt.Sprintf("azapi_resource.%s", r.NewLabel)
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

func (r *AzurermResource) EmptyImportConfig() string {
	config := ""
	if r.IsMultipleResources() {
		for _, instance := range r.Instances {
			config += fmt.Sprintf("resource \"azapi_resource\" \"%s_%s\" {}\n", r.NewLabel, strings.ReplaceAll(fmt.Sprintf("%v", instance.Index), "/", "_"))
		}
	} else {
		config += fmt.Sprintf("resource \"azapi_resource\" \"%s\" {}\n", r.NewLabel)
	}
	return config
}

func (r *AzurermResource) IsMultipleResources() bool {
	return len(r.Instances) != 0 && r.Instances[0].Index != nil
}

func (r *AzurermResource) IsForEach() bool {
	if len(r.Instances) != 0 && r.Instances[0].Index != nil {
		if _, ok := r.Instances[0].Index.(string); ok {
			return true
		}
	}
	return false
}

var pluralizeClient = pluralize.NewClient()

func NewLabel(id string, oldLabel string) string {
	resourceType := ResourceTypeOfResourceId(id)
	lastSegment := LastSegment(resourceType)
	// #nosec G404
	return fmt.Sprintf("%s_%s", pluralizeClient.Singular(lastSegment), oldLabel)
}

func LastSegment(input string) string {
	id := strings.Trim(input, "/")
	components := strings.Split(id, "/")
	if len(components) == 0 {
		return ""
	}
	return components[len(components)-1]
}

func ResourceTypeOfResourceId(input string) string {
	if input == "/" {
		return arm.TenantResourceType.String()
	}
	id := input
	if resourceType, err := arm.ParseResourceType(id); err == nil {
		if resourceType.Type != arm.ProviderResourceType.Type {
			return resourceType.String()
		}
	}

	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	resourceType := ""
	provider := ""
	for current := 0; current < len(components)-1; current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return ""
		}

		if key == "providers" {
			provider = value
			resourceType = provider
		} else if len(provider) > 0 {
			resourceType += "/" + key
		}
	}
	return resourceType
}

func sortAttributes(input *hclwrite.Block) *hclwrite.Block {
	output := hclwrite.NewBlock(input.Type(), input.Labels())
	attrList := []string{"count", "for_each", "type", "parent_id", "name", "location", "identity", "body", "tags"}
	usedAttr := make(map[string]bool)
	for _, attr := range attrList {
		if attribute := input.Body().GetAttribute(attr); attribute != nil {
			output.Body().SetAttributeRaw(attr, attribute.Expr().BuildTokens(nil))
			usedAttr[attr] = true
		} else {
			for _, block := range input.Body().Blocks() {
				if block.Type() == attr {
					output.Body().AppendBlock(block)
					usedAttr[attr] = true
				}
			}
		}
	}
	for attrName, attribute := range input.Body().Attributes() {
		if _, ok := usedAttr[attrName]; !ok {
			output.Body().SetAttributeRaw(attrName, attribute.Expr().BuildTokens(nil))
		}
	}
	for _, block := range input.Body().Blocks() {
		if _, ok := usedAttr[block.Type()]; !ok {
			output.Body().AppendBlock(block)
		}
	}
	return output
}

func AzurermIdToAzureId(azurermResourceType string, azurermId string) (string, error) {
	switch azurermResourceType {
	case "azurerm_monitor_diagnostic_setting":
		// input: <target id>|<diagnostic setting name>
		// output: <target id>/providers/Microsoft.Insights/diagnosticSettings/<diagnostic setting name>
		azurermIdSplit := strings.Split(azurermId, "|")
		if len(azurermIdSplit) != 2 {
			return "", fmt.Errorf("invalid id: %s, expected format: <target id>|<diagnostic setting name>", azurermId)
		}
		return fmt.Sprintf("%s/providers/Microsoft.Insights/diagnosticSettings/%s", azurermIdSplit[0], azurermIdSplit[1]), nil
	case "azurerm_role_definition":
		// input: <role definition id>|<scope>
		// output: <role definition id>
		azurermIdSplit := strings.Split(azurermId, "|")
		if len(azurermIdSplit) != 2 {
			return "", fmt.Errorf("invalid id: %s, expected format: <role definition id>|<scope>", azurermId)
		}
		return azurermIdSplit[0], nil

		// add more cases here as needed
	}
	// return azure id
	return azurermId, nil
}
