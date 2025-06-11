package errors

import (
	"errors"
	"reflect"
)

// Value returns the last (more recent) value associated with the given key from the error chain.
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

// ValueT returns the value associated with the given key from the error chain, cast to type T.
func ValueT[T any](err error, key any) T {
	valT, _ := Value(err, key).(T)
	return valT
}

// Values returns a slice of values associated with the given key from the error chain.
// It traverses the error chain and collects all values that match the specified key.
// If there are multiple values for the same key, all of them are included in the slice.
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

// ValuesT returns a slice of values associated with the given key from the error chain, cast to type T.
// It traverses the error chain and collects all values that match the specified key.
// If there are multiple values for the same key, all of them are included in the slice.
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

// ValueAllSlice returns a slice of all values from the error chain.
// It skips built-in key-value pairs like code, severity, operation, and formatter.
// If there are multiple values for the same key, only the first occurrence (last added) is included.
func ValueAllSlice(err error) []KeyValuer {
	var values []KeyValuer
	processed := make(map[any]struct{})

	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if isBuiltInKeyValuer(e.keyval.Key()) {
				continue
			}
			if _, exists := processed[e.keyval.Key()]; !exists {
				values = append(values, e.keyval)
				processed[e.keyval.Key()] = struct{}{}
			}
		}
	}

	return values
}

// ValuesMapOf returns a map of key-value pairs from the error chain, filtered by the specified key type.
// It collects all values associated with the same key, allowing multiple values for the same key.
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

// ValueMap returns a map of key-value pairs from the error chain.
// It skips built-in key-value pairs like code, severity, operation, and formatter.
// If there are multiple values for the same key, only the first occurrence (last added) is included.
func ValueMap(err error) map[any]any {
	m := make(map[any]any)
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if isBuiltInKeyValuer(e.keyval.Key()) {
				continue
			}
			if _, ok := m[e.keyval.Key()]; !ok {
				m[e.keyval.Key()] = e.keyval.Value()
			}
		}
	}

	return m
}

// ValueMapOf returns a map of key-value pairs from the error chain, filtered by the specified key type.
// If there are multiple values for the same key, only the first occurrence (last added) is included.
func ValueMapOf(err error, keyType any) map[any]any {
	m := make(map[any]any)
	for ; err != nil; err = errors.Unwrap(err) {
		if e, ok := err.(Error); ok {
			if reflect.TypeOf(e.keyval.Key()) == reflect.TypeOf(keyType) {
				if _, ok := m[e.keyval.Key()]; !ok {
					m[e.keyval.Key()] = e.keyval.Value()
				}
			}
		}
	}

	return m
}

func isBuiltInKeyValuer(key any) bool {
	switch key {
	case codeKey{}, severityKey{}, opKey{}, formatterKey{}:
		return true
	default:
		return false
	}
}
