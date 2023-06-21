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

func FromPtr[T any](ptr *T) *Opt[T] {
	if ptr == nil {
		return New[T]()
	}
	return FromVal(*ptr)
}

func (x *Opt[T]) Value() (T, bool) {
	return x.Val, x.HasVal
}

// If the optional has a value, return it;
// otherwise, return the specified val.
func (x *Opt[T]) Or(val T) T {
	if x.HasVal {
		return x.Val
	}
	return val
}
