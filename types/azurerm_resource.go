package types

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/aztfmigrate/helper"
	"github.com/Azure/aztfmigrate/tf"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/gertd/go-pluralize"
	_ "github.com/gertd/go-pluralize"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
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
	} else {
		// import and build combined block
		log.Printf("[INFO] generating config...")
		blocks := make([]*hclwrite.Block, 0)
		for _, instance := range r.Instances {
			if block, err := importAndGenerateConfig(terraform, fmt.Sprintf("%s.%s_%v", r.NewResourceType, r.NewLabel, instance.Index), instance.ResourceId, "", true); err == nil {
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
	r.Migrated = true
	return nil
}

func (r *AzurermResource) ImportBlock() *hclwrite.Block {
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
		importBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: r.NewResourceType}, hcl.TraverseAttr{Name: fmt.Sprintf("%s[each.value]", r.NewLabel)}})
	} else {
		importBlock.Body().SetAttributeValue("id", cty.StringVal(r.Instances[0].ResourceId))
		importBlock.Body().SetAttributeTraversal("to", hcl.Traversal{hcl.TraverseRoot{Name: r.NewResourceType}, hcl.TraverseAttr{Name: r.NewLabel}})
	}
	return importBlock
}

func (r *AzurermResource) RemovedBlock() *hclwrite.Block {
	removedBlock := hclwrite.NewBlock("removed", nil)
	removedBlock.Body().SetAttributeTraversal("from", hcl.Traversal{hcl.TraverseRoot{Name: r.OldResourceType}, hcl.TraverseAttr{Name: r.OldLabel}})
	removedLifecycleBlock := hclwrite.NewBlock("lifecycle", nil)
	removedLifecycleBlock.Body().SetAttributeValue("destroy", cty.BoolVal(false))
	removedBlock.Body().AppendBlock(removedLifecycleBlock)
	return removedBlock
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
			config += fmt.Sprintf("resource \"azapi_resource\" \"%s_%v\" {}\n", r.NewLabel, instance.Index)
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

func NewLabel(id string) string {
	resourceType := ResourceTypeOfResourceId(id)
	lastSegment := LastSegment(resourceType)
	// #nosec G404
	return fmt.Sprintf("%s_%d", pluralizeClient.Singular(lastSegment), rand.Intn(100))
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
