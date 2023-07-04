package schema

import "github.com/zclconf/go-cty/cty"

type ProviderSchema struct {
	// The provder version. This is defined in the vcs.
	Version string

	ResourceSchemas map[string]*Schema `json:"resource_schemas,omitempty"`
}

type Schema struct {
	Block *Block `json:"block,omitempty"`
}

type Block struct {
	Attributes   map[string]*Attribute   `json:"attributes,omitempty"`
	NestedBlocks map[string]*NestedBlock `json:"block_types,omitempty"`
}

type NestingMode int

const (
	NestingModeInvalid NestingMode = iota
	NestingSingle
	NestingGroup
	NestingList
	NestingSet
	NestingMap
)

type NestedBlock struct {
	NestingMode NestingMode `json:"nesting_mode,omitempty"`
	Block       *Block      `json:"block,omitempty"`

	Required bool `json:"required,omitempty"`
	Optional bool `json:"optional,omitempty"`
	Computed bool `json:"computed,omitempty"`
	ForceNew bool `json:"force_new,omitempty"`

	ConflictsWith []string `json:"conflicts_with,omitempty"`
	ExactlyOneOf  []string `json:"exactly_one_of,omitempty"`
	AtLeastOneOf  []string `json:"at_least_one_of,omitempty"`
	RequiredWith  []string `json:"required_with,omitempty"`

	MinItems int `json:"min_items,omitempty"`
	MaxItems int `json:"max_items,omitempty"`
}

type Attribute struct {
	Type cty.Type `json:"type,omitempty"`

	Required bool `json:"required,omitempty"`
	Optional bool `json:"optional,omitempty"`
	Computed bool `json:"computed,omitempty"`
	ForceNew bool `json:"force_new,omitempty"`

	Default   interface{} `json:"default,omitempty"`
	Sensitive bool        `json:"sensitive,omitempty"`

	ConflictsWith []string `json:"conflicts_with,omitempty"`
	ExactlyOneOf  []string `json:"exactly_one_of,omitempty"`
	AtLeastOneOf  []string `json:"at_least_one_of,omitempty"`
	RequiredWith  []string `json:"required_with,omitempty"`
}
