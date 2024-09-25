package schema

// The schema definition is referencing the github.com/hashicorp/terraform-plugin-go/tfprotov6/schema.go@v0.23.0
// As tfprotov6 is compatible to the tfprotov5 (that SDKv2 is using).

import "github.com/zclconf/go-cty/cty"

type ProviderSchema struct {
	Provider          *Schema            `json:"provider,omitempty"`
	ResourceSchemas   map[string]*Schema `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]*Schema `json:"data_source_schemas,omitempty"`
}

type Schema struct {
	Version int64        `json:"schema_version,omitempty"`
	Block   *SchemaBlock `json:"block,omitempty"`
}

type SchemaBlock struct {
	Attributes SchemaAttributes   `json:"attributes,omitempty"`
	BlockTypes SchemaNestedBlocks `json:"block_types,omitempty"`
}

type SchemaAttributes []*SchemaAttribute

type SchemaNestedBlocks []*SchemaNestedBlock

func (attrs SchemaAttributes) Map() map[string]*SchemaAttribute {
	m := map[string]*SchemaAttribute{}
	for _, attr := range attrs {
		m[attr.Name] = attr
	}
	return m
}

func (blocks SchemaNestedBlocks) Map() map[string]*SchemaNestedBlock {
	m := map[string]*SchemaNestedBlock{}
	for _, blk := range blocks {
		m[blk.TypeName] = blk
	}
	return m
}

type SchemaNestedBlockNestingMode int

const (
	SchemaNestedBlockNestingModeInvalid SchemaNestedBlockNestingMode = 0
	SchemaNestedBlockNestingModeSingle  SchemaNestedBlockNestingMode = 1
	SchemaNestedBlockNestingModeList    SchemaNestedBlockNestingMode = 2
	SchemaNestedBlockNestingModeSet     SchemaNestedBlockNestingMode = 3
	SchemaNestedBlockNestingModeMap     SchemaNestedBlockNestingMode = 4
	SchemaNestedBlockNestingModeGroup   SchemaNestedBlockNestingMode = 5
)

type SchemaNestedBlock struct {
	TypeName string                       `json:"type_name,omitempty"`
	Block    *SchemaBlock                 `json:"block,omitempty"`
	Nesting  SchemaNestedBlockNestingMode `json:"nesting_mode,omitempty"`

	MinItems int `json:"min_items,omitempty"`
	MaxItems int `json:"max_items,omitempty"`

	// Extended properties
	// SDKv2 Only
	Required      *bool    `json:"required,omitempty"`
	Optional      *bool    `json:"optional,omitempty"`
	Computed      *bool    `json:"computed,omitempty"`
	ForceNew      *bool    `json:"force_new,omitempty"`
	ConflictsWith []string `json:"conflicts_with,omitempty"`
	ExactlyOneOf  []string `json:"exactly_one_of,omitempty"`
	AtLeastOneOf  []string `json:"at_least_one_of,omitempty"`
	RequiredWith  []string `json:"required_with,omitempty"`
}

type SchemaObject struct {
	Attributes SchemaAttributes
	Nesting    SchemaObjectNestingMode
}

type SchemaObjectNestingMode int32

const (
	SchemaObjectNestingModeInvalid SchemaObjectNestingMode = 0
	SchemaObjectNestingModeSingle  SchemaObjectNestingMode = 1
	SchemaObjectNestingModeList    SchemaObjectNestingMode = 2
	SchemaObjectNestingModeSet     SchemaObjectNestingMode = 3
	SchemaObjectNestingModeMap     SchemaObjectNestingMode = 4
)

type SchemaAttribute struct {
	Name string    `json:"name,omitempty"`
	Type *cty.Type `json:"type,omitempty"`
	// Proto6 Only
	NestedType *SchemaObject `json:"nested_type,omitempty"`

	Required  bool `json:"required,omitempty"`
	Optional  bool `json:"optional,omitempty"`
	Computed  bool `json:"computed,omitempty"`
	Sensitive bool `json:"sensitive,omitempty"`

	// Extended Properties
	// SDKv2: Go types
	// FW: attr.Value
	Default interface{} `json:"default,omitempty"`

	// SDKv2 Only
	ForceNew      *bool    `json:"force_new,omitempty"`
	ConflictsWith []string `json:"conflicts_with,omitempty"`
	ExactlyOneOf  []string `json:"exactly_one_of,omitempty"`
	AtLeastOneOf  []string `json:"at_least_one_of,omitempty"`
	RequiredWith  []string `json:"required_with,omitempty"`
}
