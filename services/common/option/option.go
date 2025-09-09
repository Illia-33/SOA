package option

type Option[T any] struct {
	Value    T
	HasValue bool
}

func (o *Option[T]) ToPointer() *T {
	if !o.HasValue {
		return nil
	}

	v := o.Value
	return &v
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		Value:    v,
		HasValue: true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{
		HasValue: false,
	}
}

func FromPointer[T any](p *T) Option[T] {
	if p == nil {
		return None[T]()
	}

	return Some(*p)
}
