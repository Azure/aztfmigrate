package types

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
	"strings"
)

type Reference struct {
	Name  string
	Value interface{}
}

func (r Reference) GetStringValue() string {
	switch v := r.Value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type Output struct {
	OldName string
	NewName string
	Value   interface{}
}

func (r Output) GetStringValue() string {
	switch v := r.Value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type Instance struct {
	Index      interface{}
	ApiVersion string
	ResourceId string
	Outputs    []Output
}

type GenericResource struct {
	Label            string
	Instances        []Instance
	ResourceType     string
	Block            *hclwrite.Block
	References       []Reference
	InputProperties  []string
	OutputProperties []string
	Migrated         bool
}

func (r GenericResource) OldAddress(index interface{}) string {
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

func (r GenericResource) NewAddress(index interface{}) string {
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

func (r GenericResource) EmptyImportConfig() string {
	if r.IsMultipleResources() {
		if !r.IsForEach() {
			return fmt.Sprintf(`resource "%s" "%s" {
  count = %d
}`, r.ResourceType, r.Label, len(r.Instances))
		}
		keys := make([]string, 0)
		for _, instance := range r.Instances {
			keys = append(keys, fmt.Sprintf(`"%s"`, instance.Index))
		}
		return fmt.Sprintf(`resource "%s" "%s" {
	  for_each = toset([%s])
}`, r.ResourceType, r.Label, strings.Join(keys, ","))
	}
	return fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
}

func (r GenericResource) IsMultipleResources() bool {
	return len(r.Instances) != 0 && r.Instances[0].Index != nil
}

func (r GenericResource) IsForEach() bool {
	if len(r.Instances) != 0 && r.Instances[0].Index != nil {
		if _, ok := r.Instances[0].Index.(string); ok {
			return true
		}
	}
	return false
}

type GenericUpdateResource struct {
	ApiVersion       string
	Label            string
	OldLabel         string
	Id               string
	ResourceType     string
	Change           *tfjson.Change
	Block            *hclwrite.Block
	Migrated         bool
	References       []Reference
	Outputs          []Output
	InputProperties  []string
	OutputProperties []string
}

func (r GenericUpdateResource) OldAddress() string {
	return fmt.Sprintf("azapi_update_resource.%s", r.OldLabel)
}

func (r GenericUpdateResource) NewAddress() string {
	return fmt.Sprintf("%s.%s", r.ResourceType, r.Label)
}

func (r GenericUpdateResource) EmptyImportConfig() string {
	return fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
}
