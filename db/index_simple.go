package db

import (
	"log/slog"
	"sync"
)

var _ Index[Tableable] = (*SimpleIndex[Tableable])(nil)

type SimpleIndex[V Tableable] struct {
	tex sync.RWMutex

	data   Map[any, *V]
	key    func(*V) []any
	unique bool
}

func NewUniqueSimpleIndex[V Tableable](key func(*V) []any) *SimpleIndex[V] {
	return &SimpleIndex[V]{
		tex:    sync.RWMutex{},
		data:   Map[any, *V]{},
		key:    key,
		unique: true,
	}
}

func (p *SimpleIndex[V]) Index(v *V) error {
	slog.Debug("indexing", "v", v)

	for _, k := range p.key(v) {
		_v, exists := p.data.LoadOrStore(k, v)
		if p.unique && exists && _v != nil {
			continue
		}
	}

	return nil
}

func (p *SimpleIndex[V]) Lookup(v V) *V {
	vv := p.key(&v)

	for _, key := range vv {
		if v, ok := p.data.Load(key); ok {
			return v
		}
	}

	slog.Debug("lookup failed", "v", v, "vv", vv)

	return nil
}
