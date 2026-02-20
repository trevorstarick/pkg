package db

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
)

func MarshalJSONPointer[T any](p T) ([]byte, error) {
	res := map[string]any{}

	t := reflect.TypeFor[T]()
	v := reflect.ValueOf(p)

	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return json.Marshal(nil)
		}

		_, pkey, err := GetPrimaryKey(v.Interface())
		if err == nil {
			return json.Marshal(pkey)
		}

		t = t.Elem()
		v = v.Elem()
	}

	for i := range t.NumField() {
		field := t.Field(i)
		fieldValue := v.Field(i)

		if field.Name[0] >= 'a' && field.Name[0] <= 'z' {
			continue
		}

		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}

		// Parse json tag to extract name and options
		parts := strings.Split(name, ",")
		name = parts[0]

		if name == "-" {
			continue
		}

		if len(parts) == 1 {
			parts = append(parts, "") // ensure at least two parts
		}

		hasOmitEmpty := slices.Contains(parts[1:], "omitempty")

		//nolint:exhaustive // todo but we only care about slices, maps, pointers, and basic types for now
		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Map, reflect.Pointer:
			if fieldValue.IsNil() {
				if !hasOmitEmpty {
					res[name] = nil
				}

				continue
			}
		case reflect.String:
			if hasOmitEmpty && fieldValue.Len() == 0 {
				continue
			}
		case reflect.Bool:
			if hasOmitEmpty && !fieldValue.Bool() {
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if hasOmitEmpty && fieldValue.Int() == 0 {
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if hasOmitEmpty && fieldValue.Uint() == 0 {
				continue
			}
		case reflect.Float32, reflect.Float64:
			if hasOmitEmpty && fieldValue.Float() == 0 {
				continue
			}
		}

		hasOmitZero := slices.Contains(parts[1:], "omitzero")
		if hasOmitZero && fieldValue.Elem().IsZero() {
			continue
		}

		//nolint:exhaustive // todo but we only care about slices, maps, pointers, and basic types for now
		switch fieldValue.Kind() {
		case reflect.Slice, reflect.Array:
			arr := make([]any, 0, fieldValue.Len())

			for i := range fieldValue.Len() {
				elem := fieldValue.Index(i)
				if elem.Kind() == reflect.Pointer {
					_, pkey, err := GetPrimaryKey(elem.Interface())
					if err == nil {
						arr = append(arr, pkey)
					}
				}
			}

			if len(arr) > 0 {
				res[name] = arr
			} else {
				res[name] = fieldValue.Interface()
			}
		case reflect.Map:
			m := map[any]any{}

			for i := range fieldValue.Len() {
				mapKey := fieldValue.MapKeys()[i]

				mapValue := fieldValue.MapIndex(mapKey)
				if mapValue.Kind() == reflect.Pointer {
					_, pkey, err := GetPrimaryKey(mapValue.Interface())
					if err == nil {
						m[mapKey.Interface()] = pkey
					}
				}
			}

			if len(m) > 0 {
				res[name] = m
			} else {
				res[name] = fieldValue.Interface()
			}
		case reflect.Pointer:
			_, pkey, err := GetPrimaryKey(fieldValue.Interface())
			if err == nil {
				res[name] = pkey
			} else {
				res[name] = fieldValue.Interface()
			}
		default:
			res[name] = fieldValue.Interface()
		}
	}

	return json.Marshal(res)
}
