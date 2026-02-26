package db

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/tidwall/btree"
)

type Tableable interface {
	IsPlaceholder() bool
	ResolvePointers(db *DB) (any, error)
}

type Indexable[V Tableable] interface {
	AddIndex(key string, index Index[V])
	Reindex(index string) error
	Lookup(index string, v V) *V
}

type Persistent interface {
	Load(paths ...string) error
	Save(path string) error
}

type table[V Tableable] interface {
	Indexable[V]
	Persistent

	Set(value V) *V
	Get(key V) *V
	Delete(key V) bool

	Contains(v V) bool

	Iter(yield func(any, *V) bool)
	AddSort(name string, fn func(a, b *V) int)
	SortedIter(name string, asc bool, yield func(any, *V) bool)
	Len() int
}

var _ table[Tableable] = (*Table[Tableable])(nil)

type Table[V Tableable] struct {
	changedSinceLastSave atomic.Bool
	data                 Map[any, *V]
	uniqueIndices        Map[string, Index[V]]
	indicies             Map[string, Index[V]]
	sortTrees            Map[string, *btree.BTreeG[*V]]
}

func NewTable[V Tableable]() *Table[V] {
	table := Table[V]{
		changedSinceLastSave: atomic.Bool{},
		data:                 Map[any, *V]{},
		uniqueIndices:        Map[string, Index[V]]{},
		indicies:             Map[string, Index[V]]{},
		sortTrees:            Map[string, *btree.BTreeG[*V]]{},
	}

	t := reflect.TypeFor[V]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := range t.NumField() {
		field := t.Field(i)

		dbTags := strings.Split(field.Tag.Get("db"), ",")

		for _, tag := range dbTags {
			if strings.Contains(tag, "unique") {
				index := table.uniqueAutoIndexer(field, i)
				if index != nil {
					table.uniqueIndices.Store(dbTags[0], index)
				}

				break
			}
		}
	}

	return &table
}

func (t *Table[V]) uniqueAutoIndexer(field reflect.StructField, i int) Index[V] {
	dbTags := strings.Split(field.Tag.Get("db"), ",")

	//nolint:exhaustive // will implement later if there is demand for it, but for now just log a warning and skip it
	switch field.Type.Kind() {
	case reflect.Slice:
		slog.Info("table: adding unique index on slice field", "v", fmt.Sprintf("%T", *new(V)), "field", dbTags[0])

		return NewUniqueSimpleIndex[V](func(v *V) []any {
			slice := reflect.ValueOf(*v).Field(i)

			result := make([]any, slice.Len())
			for j := range slice.Len() {
				result[j] = slice.Index(j).Interface()
			}

			return result
		})
	case reflect.Map:
		slog.Warn("table: unique index on map field is not supported", "v", fmt.Sprintf("%T", *new(V)), "field", dbTags[0])

		return nil
	case reflect.Struct:
		slog.Warn("table: unique index on struct field that does not implement UniqueKeyer is not supported",
			"v", fmt.Sprintf("%T", *new(V)),
			"field", dbTags[0],
		)

		return nil
	default:
		slog.Info("table: adding unique index on field", "v", fmt.Sprintf("%T", *new(V)), "field", dbTags[0])

		return NewUniqueSimpleIndex[V](func(v *V) []any {
			return []any{reflect.ValueOf(*v).Field(i).Interface()}
		})
	}
}

func (t *Table[V]) AddIndex(key string, index Index[V]) {
	slog.Debug("table addindex", "v", fmt.Sprintf("%T", *new(V)), "key", key)

	if t.indicies.Contains(key) {
		slog.Warn("table addindex: index already exists", "v", fmt.Sprintf("%T", *new(V)), "key", key)

		return
	}

	t.indicies.Store(key, index)
}

func (t *Table[V]) AddSort(key string, fn func(a, b *V) int) {
	slog.Debug("table addsort", "v", fmt.Sprintf("%T", *new(V)), "key", key)

	sort := btree.NewBTreeG[*V](func(a, b *V) bool {
		return fn(a, b) < 0
	})

	t.Iter(func(_ any, v *V) bool {
		sort.Set(v)

		return true
	})

	t.sortTrees.Store(key, sort)
}

func (t *Table[V]) checkForUnique(value V) (*V, bool) {
	var v *V

	found := false

	t.uniqueIndices.Range(func(_ string, i Index[V]) bool {
		if vv := i.Lookup(value); vv != nil {
			v = vv
			found = true

			return false
		}

		return true
	})

	return v, found
}

func (t *Table[V]) Set(value V) *V {
	slog.Debug("table set", "v", fmt.Sprintf("%T", *new(V)), "value", value)

	_, pkey, err := GetPrimaryKey(value)
	if err != nil {
		slog.Warn("table set: failed to get primary key", "v", fmt.Sprintf("%T", *new(V)), "value", value, "error", err)

		return nil
	}

	if v, ok := t.data.Load(pkey); ok && !(*v).IsPlaceholder() {
		slog.Warn("table set: primary key already exists", "v", fmt.Sprintf("%T", *new(V)), "value", value)

		return nil
	}

	if v, ok := t.checkForUnique(value); ok {
		slog.Warn("table set: unique key already exists", "v", fmt.Sprintf("%T", *new(V)), "value", value)

		return v
	}

	t.data.Store(pkey, &value)
	t.changedSinceLastSave.Store(true)

	t.index(&value)
	t.sortTrees.Range(func(_ string, tree *btree.BTreeG[*V]) bool {
		tree.Set(&value)

		return true
	})

	return &value
}

func (t *Table[V]) Get(v V) *V {
	slog.Debug("table get", "v", fmt.Sprintf("%T", *new(V)))

	_, pkey, err := GetPrimaryKey(v)
	if err == nil && pkey != "" {
		vv, ok := t.data.Load(pkey)
		if ok {
			return vv
		}
	}

	res := (*V)(nil)

	t.uniqueIndices.Range(func(key string, index Index[V]) bool {
		if vv := index.Lookup(v); vv != nil {
			res = vv

			slog.Debug("table get: found value in unique index", "v", fmt.Sprintf("%T", *new(V)), "key", key, "index", index)

			return false
		}

		slog.Debug("table get: key not found in index", "v", fmt.Sprintf("%T", *new(V)), "key", key, "index", index)

		return true
	})

	if res != nil {
		return res
	}

	slog.Debug("table get: key not found in any unique index", "v", fmt.Sprintf("%T", *new(V)), "key", v)

	t.indicies.Range(func(key string, index Index[V]) bool {
		if vv := index.Lookup(v); vv != nil {
			res = vv

			slog.Debug("table get: key found in index", "v", fmt.Sprintf("%T", *new(V)), "key", key, "index", index)

			return false
		}

		slog.Debug("table get: key not found in index", "v", fmt.Sprintf("%T", *new(V)), "key", key, "index", index)

		return true
	})

	return res
}

func (t *Table[V]) Lookup(name string, v V) *V {
	slog.Debug("table lookup", "v", fmt.Sprintf("%T", *new(V)), "index", name)

	index, ok := t.indicies.Load(name)
	if !ok {
		slog.Warn("table lookup: index not found", "v", fmt.Sprintf("%T", *new(V)), "index", name)

		return nil
	}

	return index.Lookup(v)
}

func (t *Table[V]) Delete(key V) bool {
	slog.Debug("table delete", "v", fmt.Sprintf("%T", *new(V)), "key", key)

	_, pkey, err := GetPrimaryKey(key)
	if err != nil {
		slog.Warn("table delete: failed to get primary key", "v", fmt.Sprintf("%T", *new(V)), "key", key, "error", err)

		return false
	}

	_, ok := t.data.LoadAndDelete(pkey)
	if !ok {
		slog.Debug("table delete: key not found", "v", fmt.Sprintf("%T", *new(V)), "key", key)

		return false
	}

	t.changedSinceLastSave.Store(true)

	return true
}

func (t *Table[V]) Reindex(name string) error {
	slog.Debug("table reindex", "v", fmt.Sprintf("%T", *new(V)), "index", name)

	index, ok := t.indicies.Load(name)
	if !ok {
		slog.Warn("table reindex: index not found", "v", fmt.Sprintf("%T", *new(V)), "index", name)

		return nil
	}

	t.data.Range(func(_ any, v *V) bool {
		err := index.Index(v)
		if err != nil {
			slog.Warn("table reindex: index constraint violation",
				"v", fmt.Sprintf("%T", *new(V)),
				"index", index,
				"value", v,
				"error", err,
			)
		}

		return true
	})

	return nil
}

func (t *Table[V]) Contains(v V) bool {
	return t.Get(v) != nil
}

func (t *Table[V]) Iter(yield func(any, *V) bool) {
	t.data.Range(yield)
}

func (t *Table[V]) SortedIter(name string, asc bool, yield func(any, *V) bool) {
	l, ok := t.sortTrees.Load(name)
	if !ok {
		slog.Warn("table sortediter: sort not found, defaulting to unsorted iteration", "v", fmt.Sprintf("%T", *new(V)), "sort", name)

		t.data.SortedRange(yield)

		return
	}

	if !asc {
		l.Reverse(func(v *V) bool {
			return yield(nil, v)
		})

		return
	}

	l.Scan(func(v *V) bool {
		return yield(nil, v)
	})
}

func (t *Table[V]) Len() int {
	slog.Debug("table len", "v", fmt.Sprintf("%T", *new(V)))

	return t.data.Len()
}

func (t *Table[V]) exists(value V) bool {
	_, pkey, err := GetPrimaryKey(value)
	if err != nil {
		slog.Warn("table exists: failed to get primary key", "v", fmt.Sprintf("%T", *new(V)), "value", value, "error", err)

		return false
	}

	return t.data.Contains(pkey)
}

func (t *Table[V]) append(value V) *V {
	_, pkey, err := GetPrimaryKey(value)
	if err != nil {
		slog.Warn("table append: failed to get primary key", "v", fmt.Sprintf("%T", *new(V)), "value", value, "error", err)

		return nil
	}

	return t.data.StoreAndReturn(pkey, &value)
}

func (t *Table[V]) index(v *V) {
	slog.Debug("table index", "v", fmt.Sprintf("%T", *new(V)))

	t.indicies.Range(func(s string, i Index[V]) bool {
		err := i.Index(v)
		if err != nil {
			slog.Warn("table index constraint violation", "v", fmt.Sprintf("%T", *new(V)), "value", *v, "key", s, "error", err)
		}

		return true
	})
}
