package db

import (
	"log/slog"
	"slices"
)

type Sort[V Tableable] struct {
	data []*V
	key  func(a, b *V) int
}

func NewSort[V Tableable](key func(a, b *V) int) *Sort[V] {
	return &Sort[V]{
		data: []*V{},
		key:  key,
	}
}

func (s *Sort[V]) Add(v *V) {
	s.data = append(s.data, v)
}

func (s *Sort[V]) Sort() {
	slices.SortStableFunc(s.data, s.key)
}

func (s *Sort[V]) Iter() func(yield func(any, *V) bool) {
	index := 0

	return func(yield func(any, *V) bool) {
		for index < len(s.data) {
			v := s.data[index]
			index++

			_, pkey, err := GetPrimaryKey(v)
			if err != nil {
				slog.Warn("failed to get primary key for sorting", "error", err)

				continue
			}

			if !yield(pkey, v) {
				break
			}
		}
	}
}
