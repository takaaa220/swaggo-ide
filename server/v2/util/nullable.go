package util

import "encoding/json"

type Nullable[T any] struct {
	Value T
	Valid bool
}

func Some[T any](value T) Nullable[T] {
	return Nullable[T]{Value: value, Valid: true}
}

func None[T any]() Nullable[T] {
	return Nullable[T]{}
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*n = None[T]()
		return nil
	}
	return json.Unmarshal(data, &n.Value)
}
