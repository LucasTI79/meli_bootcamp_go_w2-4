package optional

type Opt[T any] struct {
	Val    T
	HasVal bool
}

func New[T any]() *Opt[T] {
	var zero T
	return &Opt[T]{zero, false}
}

func FromVal[T any](val T) *Opt[T] {
	return &Opt[T]{val, true}
}

func (x *Opt[T]) Value() (T, bool) {
	return x.Val, x.HasVal
}
