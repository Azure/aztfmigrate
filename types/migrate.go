package types

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclwrite"
	tfjson "github.com/hashicorp/terraform-json"
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

type GenericResource struct {
	Label        string
	Id           string
	ResourceType string
	Block        *hclwrite.Block
	Migrated     bool
	References   []Reference
	Outputs      []Output
}

func (r GenericResource) OldAddress() string {
	return fmt.Sprintf("azurerm-restapi_resource.%s", r.Label)
}

func (r GenericResource) NewAddress() string {
	return fmt.Sprintf("%s.%s", r.ResourceType, r.Label)
}

func (r GenericResource) EmptyImportConfig() string {
	return fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
}

type GenericPatchResource struct {
	Label        string
	OldLabel     string
	Id           string
	ResourceType string
	Change       *tfjson.Change
	Block        *hclwrite.Block
	Migrated     bool
	References   []Reference
	Outputs      []Output
}

func (r GenericPatchResource) OldAddress() string {
	return fmt.Sprintf("azurerm-restapi_patch_resource.%s", r.OldLabel)
}

func (r GenericPatchResource) NewAddress() string {
	return fmt.Sprintf("%s.%s", r.ResourceType, r.Label)
}

func (r GenericPatchResource) EmptyImportConfig() string {
	return fmt.Sprintf("resource \"%s\" \"%s\" {}\n", r.ResourceType, r.Label)
}
