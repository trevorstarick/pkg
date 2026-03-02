package db

import (
	"errors"
	"log/slog"
	"path/filepath"
	"reflect"
)

type DB struct {
	m Map[string, any]
}

func NewDB() *DB {
	return &DB{}
}

//nolint:gochecknoglobals // this is a global variable, but it's a default instance of the DB, so it's fine
var Default = NewDB()

func Register[V Tableable](db *DB, name string, table table[V]) error {
	slog.Debug("registered table", "table", name)

	if db.m.Contains(name) {
		slog.Warn("attempting to re-register a table with the same name", "name", name)

		return errors.New("attempting to re-register a table with the same name")
	}

	var zero V

	t := reflect.TypeOf(zero)

	// Handle pointer types
	if t.Kind() == reflect.Pointer {
		_ = t.Elem()

		panic("pointer types not supported for table registration")
	}

	db.m.Store(name, table)

	return nil
}

func GetTable[V Tableable](db *DB, name string) (*Table[V], error) {
	t, ok := db.m.Load(name)
	if !ok {
		slog.Warn("table does not exists", "fn", "get", "table", name)

		return nil, errors.New("table does not exist")
	}

	//nolint:forcetypeassert // we know this is a *V because of the type of the table
	return t.(*Table[V]), nil
}

func MustGetTable[V Tableable](db *DB, name string) *Table[V] {
	t, err := GetTable[V](db, name)
	if err != nil {
		panic(err)
	}

	return t
}

func LoadTable[V Tableable](db *DB, name string, dir string) error {
	t, err := GetTable[V](db, name)
	if err != nil {
		return err
	}

	err = t.Load(filepath.Join(dir, name+".jsonl"))
	if err != nil {
		return err
	}

	for key, value := range t.Iter() {
		if (*value).IsPlaceholder() {
			slog.Debug("skipping placeholder resolve", "table", name, "key", key)

			continue
		}

		if value == nil {
			slog.Error("nil value encountered during resolve", "table", name, "key", key)

			continue
		}

		v, err := (*value).ResolvePointers(db)
		if err != nil {
			slog.Error("resolving pointers", "table", name, "key", key, "error", err)

			continue
		}

		if v != nil {
			//nolint:forcetypeassert // we know this is a *V because of the type of the table
			*value = v.(V)
		}
	}
	slog.Info("finished resolving pointers for table", "table", name, "values", t.data.Len())

	for _, index := range t.uniqueIndices.Iter() {
		for _, v := range t.data.Iter() {
			err = index.Index(v)
			if err != nil {
				slog.Error("failed to index value", "fn", "load", "table", name, "error", err)
			}
		}
	}
	slog.Info("finished indexing unique indices for table", "table", name, "indices", t.uniqueIndices.Len())

	for _, tree := range t.sortTrees.Iter() {
		for _, v := range t.data.Iter() {
			tree.Set(v)
		}
	}
	slog.Info("finished indexing sort trees for table", "table", name, "trees", t.sortTrees.Len())

	return nil
}

func SaveTable[V Tableable](db *DB, name string, dir string) error {
	t, err := GetTable[V](db, name)
	if err != nil {
		return err
	}

	return t.Save(filepath.Join(dir, name+".jsonl"))
}

func MustGet[V Tableable](db *DB, name string, v V) *V {
	t, err := GetTable[V](db, name)
	if err != nil {
		panic(err)
	}

	r := t.Get(v)
	if r == nil {
		panic("value not found")
	}

	return r
}

func Get[V Tableable](db *DB, name string, v V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return new(V), err
	}

	return t.Get(v), nil
}

func Lookup[V Tableable](db *DB, name string, index string, v V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil, err
	}

	return t.Lookup(index, v), nil
}

func Iter[V Tableable](db *DB, name string) func(yield func(any, *V) bool) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil
	}

	return t.Iter()
}

func GetOrSet[V Tableable](db *DB, name string, value V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil, err
	}

	existing := t.Get(value)
	if existing != nil {
		return existing, nil
	}

	return t.Set(value), nil
}

func Set[V Tableable](db *DB, name string, value V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil, err
	}

	existing := t.Get(value)
	if existing != nil {
		slog.Debug("value already exists, not setting", "fn", "set", "table", name)

		return existing, nil
	}

	vv := t.Set(value)

	t.uniqueIndices.Iter()(func(_ string, i Index[V]) bool {
		err = i.Index(vv)
		if err != nil {
			slog.Error("failed to index value", "fn", "set", "table", name, "error", err)

			return true
		}

		return true
	})

	t.indicies.Iter()(func(_ string, i Index[V]) bool {
		err := i.Index(vv)
		if err != nil {
			slog.Error("failed to index value", "fn", "set", "table", name, "error", err)

			return true
		}

		return true
	})

	return vv, nil
}

func Update[V Tableable](db *DB, name string, value V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil, err
	}

	existing := t.Get(value)
	if existing != nil {
		*existing = value

		return existing, nil
	}

	return t.Get(value), nil
}

func Upsert[V Tableable](db *DB, name string, value V) (*V, error) {
	t, err := GetTable[V](db, name)
	if err != nil {
		return nil, err
	}

	existing := t.Get(value)
	if existing != nil {
		*existing = value

		return existing, nil
	}

	return t.Set(value), nil
}

func Delete[V Tableable](db *DB, name string, value V) error {
	t, err := GetTable[V](db, name)
	if err != nil {
		return err
	}

	t.Delete(value)

	return nil
}

func Contains[V Tableable](db *DB, name string, v V) bool {
	t, err := GetTable[V](db, name)
	if err != nil {
		return false
	}

	return t.Contains(v)
}
