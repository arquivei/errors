package errors

import (
	"fmt"
	"strconv"
)

// KeyValuer is an interface for key-value pairs that can be used in errors.
// It provides methods to retrieve the key, value, and a string representation of the value.
type KeyValuer interface {
	Key() any
	Value() any
}

type KeyValue struct {
	key   any
	value any
}

var _ KeyValuer = KeyValue{}

func (kv KeyValue) Key() any {
	return kv.key
}

func (kv KeyValue) Value() any {
	return kv.value
}

// KV is a constructor for KeyValuer types.
// It ensures at compile-time that the provided key is of a comparable type.
func KV[K comparable](key K, value any) KeyValuer {
	return KeyValue{
		key:   key,
		value: value,
	}
}

// stringify tries a bit to stringify v, without using fmt, since we don't
// want context depending on the unicode tables. This is only used by
// *valueCtx.String().
// NOTE: Extracted from the context package.
func stringify(v any) string {
	switch s := v.(type) {
	case nil:
		return "<nil>"
	case string:
		return s
	case int:
		return strconv.Itoa(s)
	case bool:
		if s {
			return "true"
		}
		return "false"
	case error:
		return s.Error()
	case fmt.Stringer:
		return s.String()
	}
	return fmt.Sprintf("%v", v)
}
