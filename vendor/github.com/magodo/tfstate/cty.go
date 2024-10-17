package tfstate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

func UnmarshalToCty(obj map[string]interface{}, t cty.Type) (cty.Value, error) {
	var path cty.Path
	v, err := unmarshal(obj, t, path)
	if err != nil {
		switch err := err.(type) {
		case cty.PathError:
			return cty.NilVal, PathError{err}
		default:
			return cty.NilVal, err
		}
	}
	return v, nil
}

func unmarshal(v interface{}, t cty.Type, path cty.Path) (cty.Value, error) {
	if v == nil {
		return cty.NullVal(t), nil
	}

	if t == cty.DynamicPseudoType {
		_, v, err := unmarshalDynamic(v, path)
		return v, err
	}

	switch {
	case t.IsPrimitiveType():
		val, err := unmarshalPrimitive(v, t, path)
		if err != nil {
			return cty.NilVal, err
		}
		return val, nil
	case t.IsListType():
		return unmarshalList(v, t.ElementType(), path)
	case t.IsSetType():
		return unmarshalSet(v, t.ElementType(), path)
	case t.IsMapType():
		return unmarshalMap(v, t.ElementType(), path)
	case t.IsTupleType():
		return unmarshalTuple(v, t.TupleElementTypes(), path)
	case t.IsObjectType():
		return unmarshalObject(v, t.AttributeTypes(), path)
	// case t.IsCapsuleType():
	// 	return unmarshalCapsule(v, t, path)
	default:
		return cty.NilVal, path.NewErrorf("unsupported type %s", t.FriendlyName())
	}
}

func unmarshalPrimitive(v interface{}, t cty.Type, path cty.Path) (cty.Value, error) {
	switch t {
	case cty.Bool:
		switch v := v.(type) {
		case bool:
			return cty.BoolVal(v), nil
		case string:
			val, err := convert.Convert(cty.StringVal(v), t)
			if err != nil {
				return cty.NilVal, path.NewError(err)
			}
			return val, nil
		default:
			return cty.NilVal, path.NewErrorf("bool is required")
		}
	case cty.Number:
		switch v := v.(type) {
		case json.Number:
			val, err := cty.ParseNumberVal(v.String())
			if err != nil {
				return cty.NilVal, path.NewError(err)
			}
			return val, nil
		case string:
			val, err := cty.ParseNumberVal(v)
			if err != nil {
				return cty.NilVal, path.NewError(err)
			}
			return val, nil
		case float64:
			return cty.NumberFloatVal(v), nil
		default:
			return cty.NilVal, path.NewErrorf("number is required, got %T", v)
		}
	case cty.String:
		switch v := v.(type) {
		case string:
			return cty.StringVal(v), nil
		case json.Number:
			return cty.StringVal(string(v)), nil
		case bool:
			val, err := convert.Convert(cty.BoolVal(v), t)
			if err != nil {
				return cty.NilVal, path.NewError(err)
			}
			return val, nil
		default:
			return cty.NilVal, path.NewErrorf("string is required")
		}
	default:
		// should never happen
		panic("unsupported primitive type")
	}
}

func unmarshalList(v interface{}, ety cty.Type, path cty.Path) (cty.Value, error) {
	l, ok := v.([]interface{})
	if !ok {
		return cty.NilVal, path.NewErrorf("expect a slice, got %T", v)
	}
	var vals []cty.Value
	{
		path := append(path, nil)
		for idx, elem := range l {
			path[len(path)-1] = cty.IndexStep{
				Key: cty.NumberIntVal(int64(idx)),
			}
			el, err := unmarshal(elem, ety, path)
			if err != nil {
				return cty.NilVal, err
			}
			vals = append(vals, el)
		}
	}

	if len(vals) == 0 {
		return cty.ListValEmpty(ety), nil
	}

	return cty.ListVal(vals), nil
}

func unmarshalSet(v interface{}, ety cty.Type, path cty.Path) (cty.Value, error) {
	l, ok := v.([]interface{})
	if !ok {
		return cty.NilVal, path.NewErrorf("expect a slice, got %T", v)
	}
	var vals []cty.Value
	{
		path := append(path, nil)
		for idx, elem := range l {
			path[len(path)-1] = cty.IndexStep{
				Key: cty.NumberIntVal(int64(idx)),
			}
			el, err := unmarshal(elem, ety, path)
			if err != nil {
				return cty.NilVal, err
			}
			vals = append(vals, el)
		}
	}

	if len(vals) == 0 {
		return cty.SetValEmpty(ety), nil
	}

	return cty.SetVal(vals), nil
}

func unmarshalMap(v interface{}, ety cty.Type, path cty.Path) (cty.Value, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return cty.NilVal, path.NewErrorf("expect a map, got %T", v)
	}
	vals := make(map[string]cty.Value)
	{
		path := append(path, nil)
		for k, v := range m {
			path[len(path)-1] = cty.IndexStep{
				Key: cty.StringVal(k),
			}
			el, err := unmarshal(v, ety, path)
			if err != nil {
				return cty.NilVal, err
			}
			vals[k] = el
		}
	}

	if len(vals) == 0 {
		return cty.MapValEmpty(ety), nil
	}

	return cty.MapVal(vals), nil
}

func unmarshalTuple(v interface{}, etys []cty.Type, path cty.Path) (cty.Value, error) {
	l, ok := v.([]interface{})
	if !ok {
		return cty.NilVal, path.NewErrorf("expect a slice, got %T", v)
	}
	var vals []cty.Value
	{
		path := append(path, nil)
		for idx, elem := range l {
			if idx >= len(etys) {
				return cty.NilVal, path[:len(path)-1].NewErrorf("too many tuple elements (need %d)", len(etys))
			}
			path[len(path)-1] = cty.IndexStep{
				Key: cty.NumberIntVal(int64(idx)),
			}
			ety := etys[idx]
			el, err := unmarshal(elem, ety, path)
			if err != nil {
				return cty.NilVal, err
			}
			vals = append(vals, el)
		}
	}

	if len(vals) != len(etys) {
		return cty.NilVal, path[:len(path)-1].NewErrorf("not enough tuple elements (need %d)", len(etys))
	}

	if len(vals) == 0 {
		return cty.EmptyTupleVal, nil
	}

	return cty.TupleVal(vals), nil
}

func unmarshalObject(v interface{}, atys map[string]cty.Type, path cty.Path) (cty.Value, error) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return cty.NilVal, path.NewErrorf("expect a map, got %T", v)
	}
	vals := make(map[string]cty.Value)

	{
		objPath := path           // some errors report from the object's perspective
		path := append(path, nil) // path to a specific attribute

		for k, v := range m {
			var err error

			aty, ok := atys[k]
			if !ok {
				return cty.NilVal, objPath.NewErrorf("unsupported attribute %q", k)
			}

			path[len(path)-1] = cty.GetAttrStep{
				Name: k,
			}

			el, err := unmarshal(v, aty, path)
			if err != nil {
				return cty.NilVal, err
			}

			vals[k] = el
		}
	}

	// Make sure we have a value for every attribute
	for k, aty := range atys {
		if _, exists := vals[k]; !exists {
			vals[k] = cty.NullVal(aty)
		}
	}

	if len(vals) == 0 {
		return cty.EmptyObjectVal, nil
	}

	return cty.ObjectVal(vals), nil
}

func unmarshalDynamic(v interface{}, path cty.Path) (cty.Type, cty.Value, error) {
	if v == nil {
		return cty.DynamicPseudoType, cty.NullVal(cty.DynamicPseudoType), nil
	}

	switch v := v.(type) {
	case bool:
		return cty.Bool, cty.BoolVal(v), nil
	case float64:
		return cty.Number, cty.NumberFloatVal(v), nil
	case string:
		return cty.String, cty.StringVal(v), nil
	case json.Number:
		val, err := cty.ParseNumberVal(v.String())
		if err != nil {
			return cty.NilType, cty.NilVal, path.NewError(err)
		}
		return cty.Number, val, nil
	case []interface{}:
		eTypes := []cty.Type{}
		eVals := []cty.Value{}
		for idx, e := range v {
			path := append(path, cty.IndexStep{
				Key: cty.NumberIntVal(int64(idx)),
			})
			eType, eVal, err := unmarshalDynamic(e, path)
			if err != nil {
				return cty.NilType, cty.NilVal, err
			}
			eTypes = append(eTypes, eType)
			eVals = append(eVals, eVal)
		}
		val := cty.TupleVal(eVals)
		typ := cty.Tuple(eTypes)
		return typ, val, nil
	case map[string]interface{}:
		attrTypes := map[string]cty.Type{}
		attrVals := map[string]cty.Value{}
		for k, v := range v {
			path := append(path, cty.GetAttrStep{
				Name: k,
			})
			attrType, attrVal, err := unmarshalDynamic(v, path)
			if err != nil {
				return cty.NilType, cty.NilVal, err
			}
			attrTypes[k] = attrType
			attrVals[k] = attrVal
		}
		typ := cty.Object(attrTypes)
		val := cty.ObjectVal(attrVals)
		return typ, val, nil
	default:
		return cty.NilType, cty.NilVal, fmt.Errorf("Unhandled type: %T", v)
	}
}

type PathError struct {
	cty.PathError
}

func (e PathError) Error() string {
	if pathStr := pathStr(e.Path); pathStr != "" {
		return pathStr + ": " + e.PathError.Error()
	}
	return e.PathError.Error()
}

func pathStr(path cty.Path) string {
	if len(path) == 0 {
		return ""
	}
	var buf strings.Builder
	for _, step := range path {
		switch step := step.(type) {
		case cty.GetAttrStep:
			fmt.Fprintf(&buf, ".%s", step.Name)
		case cty.IndexStep:
			fmt.Fprintf(&buf, "[%#v]", step.Key)
		default:
			fmt.Fprintf(&buf, "<INVALID: %#v>", step)
		}
	}
	return buf.String()
}
