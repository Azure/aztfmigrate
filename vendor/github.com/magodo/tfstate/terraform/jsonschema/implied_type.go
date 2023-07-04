package jsonschema

// This file is a mimic of: https://github.com/hashicorp/terraform/blob/20e9d8e28c32de78ebb8ae0ba2a3a00b70ee89f4/internal/configs/configschema/implied_type.go
// But instead of processing on the configschema.Block as the core does, here it processes the `tfjson.SchemaBlock`.

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/zclconf/go-cty/cty"
)

// SchemaBlockImpliedType returns the cty.Type that would result from decoding a
// configuration block using the receiving block schema.
func SchemaBlockImpliedType(b *tfjson.SchemaBlock) cty.Type {
	return schemaBlockSpecType(b).WithoutOptionalAttributesDeep()
}

func schemaBlockSpecType(b *tfjson.SchemaBlock) cty.Type {
	if b == nil {
		return cty.EmptyObject
	}
	return hcldec.ImpliedType(DecoderSpec(b))
}

// SchemaNestedAttributeTypeImpliedType returns the cty.Type that would result from decoding a
// NestedType Attribute using the SchemaNestedAttributeType.
func SchemaNestedAttributeTypeImpliedType(o *tfjson.SchemaNestedAttributeType) cty.Type {
	return schemaNestedAttributeTypeSpecType(o).WithoutOptionalAttributesDeep()
}

func schemaNestedAttributeTypeSpecType(o *tfjson.SchemaNestedAttributeType) cty.Type {
	if o == nil {
		return cty.EmptyObject
	}

	attrTys := make(map[string]cty.Type, len(o.Attributes))
	for name, attrS := range o.Attributes {
		if attrS.AttributeNestedType != nil {
			attrTys[name] = schemaNestedAttributeTypeSpecType(attrS.AttributeNestedType)
		} else {
			attrTys[name] = attrS.AttributeType
		}
	}
	optAttrs := listOptionalAttrsFromObject(o)

	var ret cty.Type
	if len(optAttrs) > 0 {
		ret = cty.ObjectWithOptionalAttrs(attrTys, optAttrs)
	} else {
		ret = cty.Object(attrTys)
	}
	switch o.NestingMode {
	case tfjson.SchemaNestingModeSingle:
		return ret
	case tfjson.SchemaNestingModeList:
		return cty.List(ret)
	case tfjson.SchemaNestingModeMap:
		return cty.Map(ret)
	case tfjson.SchemaNestingModeSet:
		return cty.Set(ret)
	default: // Should never happen
		return cty.EmptyObject
	}
}
