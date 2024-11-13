package types

import (
	"fmt"

	"github.com/Azure/aztfmigrate/tf"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

type AzureResource interface {
	TargetProvider() string
	OldAddress(index interface{}) string
	NewAddress(index interface{}) string

	CoverageCheck(strictMode bool) error
	GenerateNewConfig(terraform *tf.Terraform) error
	EmptyImportConfig() string

	StateUpdateBlocks() []*hclwrite.Block
	MigratedBlock() *hclwrite.Block

	IsMigrated() bool
	Outputs() []Output
}

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
