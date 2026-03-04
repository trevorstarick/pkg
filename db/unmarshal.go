//go:build !goexperiment.jsonv2

package db

import "encoding/json"

func (t *Table[V]) unmarshal(bytes []byte, v *V) error {
	return json.Unmarshal(bytes, v)
}
