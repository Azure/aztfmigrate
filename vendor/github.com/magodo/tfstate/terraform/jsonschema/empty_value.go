// This is derived from github.com/hashicorp/terraform/internal/configs/configschema/empty_value.go (c395d90b375e2b230384d0c213fe26a06b76222b)

package jsonschema

import (
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

// SchemaBlockEmptyValue returns the "empty value" for the recieving block, which for
// a block type is a non-null object where all of the attribute values are
// the empty values of the block's attributes and nested block types.
func SchemaBlockEmptyValue(b *tfjson.SchemaBlock) cty.Value {
	vals := make(map[string]cty.Value)
	for name, attrS := range b.Attributes {
		vals[name] = SchemaAttributeEmptyValue(attrS)
	}
	for name, blockS := range b.NestedBlocks {
		vals[name] = SchemaBlockTypeEmptyValue(blockS)
	}
	return cty.ObjectVal(vals)
}

// SchemaAttributeEmptyValue returns the "empty value" for the receiving attribute, which is
// the value that would be returned if there were no definition of the attribute
// at all, ignoring any required constraint.
func SchemaAttributeEmptyValue(a *tfjson.SchemaAttribute) cty.Value {
	return cty.NullVal(SchemaAttributeImpliedType(a))
}

//	SchemaBlockTypeEmptyValue returns the "empty value" for when there are zero nested blocks
//
// present of the receiving type.
func SchemaBlockTypeEmptyValue(b *tfjson.SchemaBlockType) cty.Value {
	switch b.NestingMode {
	case tfjson.SchemaNestingModeSingle:
		return cty.NullVal(SchemaBlockImpliedType(b.Block))
	case tfjson.SchemaNestingModeGroup:
		return SchemaBlockEmptyValue(b.Block)
	case tfjson.SchemaNestingModeList:
		if ty := SchemaBlockImpliedType(b.Block); ty.HasDynamicTypes() {
			return cty.EmptyTupleVal
		} else {
			return cty.ListValEmpty(ty)
		}
	case tfjson.SchemaNestingModeMap:
		if ty := SchemaBlockImpliedType(b.Block); ty.HasDynamicTypes() {
			return cty.EmptyObjectVal
		} else {
			return cty.MapValEmpty(ty)
		}
	case tfjson.SchemaNestingModeSet:
		return cty.SetValEmpty(SchemaBlockImpliedType(b.Block))
	default:
		// Should never get here because the above is intended to be exhaustive,
		// but we'll be robust and return a result nonetheless.
		return cty.NullVal(cty.DynamicPseudoType)
	}
}
