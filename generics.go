package generics

import "encoding/json"

// Stack ...
type Stack[T any] struct {
	data []T
}

// NewStack ...
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push ...
func (s *Stack[T]) Push(e T) {
	s.data = append(s.data, e)
}

func (s *Stack[T]) Pop() T {
	n := len(s.data)
	result := s.data[n-1]
	s.data = s.data[:n-1]
	return result
}

// Null ...
type Null[T any] struct {
	Valid bool
	Data  T
}

// NullEmpty ...
func NullEmpty[T any]() Null[T] {
	return Null[T]{}
}

func NullValue[T any](d T) Null[T] {
	return Null[T]{
		Valid: true,
		Data:  d,
	}
}

// NullMap ...
func NullMap[X, Y any](a Null[X], fn func(X) Y) Null[Y] {
	if !a.Valid {
		return Null[Y]{}
	}
	return Null[Y]{
		Valid: true,
		Data:  fn(a.Data),
	}
}

// NullAndThen ...
func NullAndThen[X, Y any](a Null[X], fn func(X) Null[Y]) Null[Y] {
	if !a.Valid {
		return Null[Y]{}
	}
	return fn(a.Data)
}

// SliceMap ...
func SliceMap[X, Y any](data []X, fn func(a X) Y) []Y {
	result := make([]Y, 0, len(data))
	for _, e := range data {
		result = append(result, fn(e))
	}
	return result
}

// GoMapMap ...
func GoMapMap[K comparable, X, Y any](data map[K]X, fn func(a X) Y) map[K]Y {
	result := map[K]Y{}
	for k, v := range data {
		result[k] = fn(v)
	}
	return result
}

// MarshalJSON ...
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Data)
}

func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*n = Null[T]{}
		return nil
	}
	n.Valid = true
	return json.Unmarshal(data, &n.Data)
}
