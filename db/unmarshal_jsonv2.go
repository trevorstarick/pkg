//go:build goexperiment.jsonv2

package db

import (
	json "encoding/json/v2"
)

type UnmarshalerFrom interface {
	json.UnmarshalerFrom
}

func (t *Table[V]) unmarshal(bytes []byte, v *V) error {
	return json.Unmarshal(bytes, v)
}
