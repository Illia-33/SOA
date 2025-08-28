package types

import "encoding/json"

type Optional[T any] struct {
	Value    T
	HasValue bool
}

func (o *Optional[T]) UnmarshalJSON(b []byte) error {
	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	*o = Optional[T]{
		Value:    value,
		HasValue: true,
	}
	return nil
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *Optional[T]) ToPointer() *T {
	if !o.HasValue {
		return nil
	}

	val := o.Value
	return &val
}
