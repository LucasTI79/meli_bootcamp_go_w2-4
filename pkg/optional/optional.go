package optional

type Optional[T any] struct {
	Val    T
	HasVal bool
}

func NewOptional[T any]() *Optional[T] {
	var zero T
	return &Optional[T]{zero, false}
}

func NewOptionalFromVal[T any](val T) *Optional[T] {
	return &Optional[T]{val, true}
}

func (x *Optional[T]) Value() (T, bool) {
	return x.Val, x.HasVal
}
