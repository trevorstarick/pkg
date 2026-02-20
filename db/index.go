package db

type Index[V Tableable] interface {
	Index(v *V) error
	Lookup(v V) *V
}
