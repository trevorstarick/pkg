package db

import (
	json "encoding/json"
	jsonv2 "encoding/json/v2"
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func UnmarshalJSONPointer[V any](b []byte, p *V) error {
	if len(b) == 0 {
		return nil
	}

	switch b[0] {
	case '"': // string
		var s string

		err := jsonv2.Unmarshal(b, &s)
		if err != nil {
			panic(fmt.Errorf("%w: %s", err, string(b)))
		}

		// Check if this is a primary_key field
		vPtr := reflect.ValueOf(p)
		if vPtr.Kind() == reflect.Pointer {
			vElem := vPtr.Elem()
			vType := vElem.Type()

			for i := range vType.NumField() {
				field := vType.Field(i)
				tagParts := strings.Split(field.Tag.Get("db"), ",")

				if len(tagParts) < 2 {
					continue
				}

				isPrimaryKey := slices.Contains(tagParts[1:], "primary_key")

				if isPrimaryKey && vElem.Field(i).CanSet() {
					fieldPtr := reflect.New(field.Type)

					err := jsonv2.Unmarshal(b, fieldPtr.Interface())
					if err != nil {
						return err
					}

					vElem.Field(i).Set(fieldPtr.Elem())
					*p = vElem.Interface().(V)

					return nil
				}
			}

			panic("no primary_key field found")
		}

		return nil
	case 'n': // null
		return nil
	case '[':
		return jsonv2.Unmarshal(b, p)
	case '{':
		if _, ok := any(p).(jsonv2.Unmarshaler); !ok {
			return jsonv2.Unmarshal(b, p)
		}

		// Unmarshal object manually
		var raw map[string]json.RawMessage

		err := jsonv2.Unmarshal(b, &raw)
		if err != nil {
			return err
		}

		// Iterate through fields of V
		v := reflect.ValueOf(p).Elem()
		t := v.Type()

		for i := range t.NumField() {
			field := t.Field(i)
			fieldValue := v.Field(i)

			// Get the json tag
			tagName := strings.Split(field.Tag.Get("db"), ",")[0]
			if tagName == "-" {
				continue
			}

			if tagName == "" {
				tagName = field.Name
			}

			if rawValue, ok := raw[tagName]; ok && fieldValue.CanSet() {
				fieldPtr := reflect.New(field.Type)

				err := jsonv2.Unmarshal(rawValue, fieldPtr.Interface())
				if err != nil {
					return err
				}

				fieldValue.Set(fieldPtr.Elem())
			}
		}
	default:
		return jsonv2.Unmarshal(b, p)
	}

	return nil
}
