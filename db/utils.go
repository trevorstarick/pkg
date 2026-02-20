package db

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

func GetPrimaryKey[V any](v V) (string, any, error) {
	if pk, ok := any(v).(interface{ PrimaryKey() any }); ok {
		value := pk.PrimaryKey()

		return "func()", value, nil
	}

	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)

	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	//nolint:exhaustive // we only care about structs for primary keys, so we can ignore other types
	switch k := rt.Kind(); k {
	case reflect.Struct:
		for i := range rt.NumField() {
			field := rt.Field(i)
			if slices.Contains(strings.Split(field.Tag.Get("db"), ","), "primary_key") {
				value := rv.Field(i).Interface()

				return field.Name, value, nil
			}
		}
	default:
		return "", nil, fmt.Errorf("unsupported type: %s", k)
	}

	return "", nil, fmt.Errorf("primary key not found in struct %s", rt.Name())
}

func GetUniqueValues[V any](v V) (map[string]any, error) {
	uniqueKeys := make(map[string]any)

	reflectType := reflect.TypeFor[V]().Elem()
	for i := range reflectType.NumField() {
		field := reflectType.Field(i)
		if field.Tag.Get("db") == "unique" || field.Tag.Get("db") == "primary_key" {
			value := reflect.ValueOf(v).Elem().Field(i).Interface()
			uniqueKeys[field.Name] = value
		}
	}

	return uniqueKeys, nil
}
