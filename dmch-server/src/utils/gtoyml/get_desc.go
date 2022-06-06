package gtoyml

import (
	"encoding"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

var (
	jsonMarshalerType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	textMarshalerType = reflect.TypeOf((*encoding.TextMarshaler)(nil)).Elem()
)

// func genField(field reflect.StructField) (map[string]interface{}, error) {

// 	desc, err := genForType(field.Type)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return desc,
// }

func (c *DocGenerator) genDescForType(rt reflect.Type) (map[string]interface{}, error) {
	if rt.Implements(textMarshalerType) {
		return map[string]interface{}{
			"type": "string",
		}, nil
	}

	switch rt.Kind() {
	case reflect.String:
		return map[string]interface{}{
			"type": "string",
		}, nil
	case reflect.Bool:
		return map[string]interface{}{
			"type": "boolean",
		}, nil
	case reflect.Float32:
		return map[string]interface{}{
			"type":   "number",
			"format": "float",
		}, nil
	case reflect.Float64:
		return map[string]interface{}{
			"type":   "number",
			"format": "double",
		}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]interface{}{
			"type":   "integer",
			"format": rt.String(),
		}, nil
	case // all number types without specific format
		reflect.Complex64, reflect.Complex128:
		return map[string]interface{}{
			"type":   "number",
			"format": rt.String(),
		}, nil

	case reflect.Array:
		desc, err := c.genDescForType(rt.Elem())
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"type":  "array",
			"items": desc,
		}, nil

	case reflect.Struct:
		props := map[string]interface{}{}
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			tags := strings.Split(field.Tag.Get("json"), ",")
			name := tags[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = field.Name
			}

			fieldDesc, err := c.genDescForType(field.Type)
			if err != nil {
				return nil, err
			}
			fieldDesc["name"] = field.Name
			props[name] = fieldDesc
		}
		return map[string]interface{}{
			"type":       "object",
			"properties": props,
		}, nil

	case reflect.Pointer:
		elem := rt.Elem()
		desc, err := c.genDescForType(elem)
		if err != nil {
			return desc, err
		}
		desc["nullable"] = true

		return desc, nil
	}

	return nil, errors.New("unsupported type: " + rt.String())
}
