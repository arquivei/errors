package errors

import (
	"errors"
	"reflect"
)

var _ error = Error{}

type Error struct {
	err    error
	keyval KeyValuer
}

func (e Error) Error() string {
	return GetFormatter(e)(e)
}

func (e Error) Unwrap() error {
	return e.err
}

func Value(err error, key any) any {
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if e.keyval.Key() == key {
				return e.keyval.Value()
			}
		}
	}

	return nil
}

func ValueKV[T KeyValuer](err error) T {
	var val T

	return ValueT[T](err, val.Key())
}

func ValueT[T any](err error, key any) T {
	valT, _ := Value(err, key).(T)
	return valT
}

func Values(err error, key any) []any {
	var values []any

	var e Error
	for ; err != nil; err = errors.Unwrap(err) {
		if errors.As(err, &e) {
			if e.keyval.Key() == key {
				values = append(values, e.keyval.Value())
			}
		}
	}

	return values
}

func ValuesT[T any](err error, key any) []T {
	values := Values(err, key)
	if len(values) == 0 {
		return nil
	}

	tValues := make([]T, 0, len(values))
	for _, v := range values {
		if v == nil {
			continue
		}

		if t, ok := v.(T); ok {
			tValues = append(tValues, t)
		}

	}
	return tValues
}

func ValuesKV[T KeyValuer](err error) []T {
	var val T

	return ValuesT[T](err, val.Key())
}

func ValuesMapOf(err error, keyType any) map[any][]any {
	m := make(map[any][]any)
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if reflect.TypeOf(e.keyval.Key()) == reflect.TypeOf(keyType) {
				m[e.keyval.Key()] = append(m[e.keyval.Key()], e.keyval.Value())
			}
		}
	}

	return m
}

func ValueMap(err error) map[any]any {
	m := make(map[any]any)
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if _, ok := m[e.keyval.Key()]; !ok {
				m[e.keyval.Key()] = e.keyval.Value()
			}
		}
	}

	return m
}

func ValueMapOf(err error, keyType any) map[any]any {
	m := make(map[any]any)
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if reflect.TypeOf(e.keyval.Key()) == reflect.TypeOf(keyType) {
				m[e.keyval.Key()] = e.keyval.Value()
			}
		}
	}

	return m
}
