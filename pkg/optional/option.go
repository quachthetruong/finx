package optional

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Optional is a generic type that can be used to represent a value of type T that may or may not be present.
type Optional[T any] struct {
	isPresent bool
	value     T
}

// Some returns an Optional[T] with isPresent set to true and value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		isPresent: true,
		value:     value,
	}
}

// None returns an Optional[T] with isPresent set to false, can be used to crete a new optional value.
func None[T any]() Optional[T] {
	return Optional[T]{
		isPresent: false,
	}
}

// IsPresent returns true if the value is present, otherwise false.
func (o *Optional[T]) IsPresent() bool {
	return o.isPresent
}

// Get returns the value of type T, if the value is not present, it will return zero value of type T.
func (o *Optional[T]) Get() T {
	return o.value
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if o.isPresent {
		return json.Marshal(o.value)
	}

	// currently do not support `omitempty` struct tag
	return json.Marshal(nil)
}

func (o *Optional[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		o.isPresent = false
		return nil
	}

	err := json.Unmarshal(b, &o.value)
	if err != nil {
		return err
	}

	o.isPresent = true
	return nil
}

// empty returns the zero value of type T.
func empty[T any]() (t T) {
	return
}

// Scan implements the SQL driver.Scanner interface.
func (o *Optional[T]) Scan(src any) error {
	if src == nil {
		o.isPresent = false
		o.value = empty[T]()
		return nil
	}

	if av, err := driver.DefaultParameterConverter.ConvertValue(src); err == nil {
		if v, ok := av.(T); ok {
			o.isPresent = true
			o.value = v
			return nil
		}
	}

	return fmt.Errorf("failed to scan Option[T]")
}

// Value implements the driver Valuer interface.
func (o Optional[T]) Value() (driver.Value, error) {
	if !o.isPresent {
		return nil, nil
	}

	return o.value, nil
}

// FromValueNonZero returns an Optional[T] from a value of type T, the value will present if it is not the zero value of T.
func FromValueNonZero[T comparable](t T) Optional[T] {
	if t == empty[T]() {
		return None[T]()
	}
	return Some[T](t)
}

// FromTime returns an Optional[time.Time] from a time.Time, the value will present if it is not the zero value of time.Time.
func FromTime(t time.Time) Optional[time.Time] {
	if t.IsZero() {
		return None[time.Time]()
	}
	return Some[time.Time](t)
}
