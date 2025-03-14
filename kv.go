package errors

import "reflect"

type KeyValuer interface {
	Key() any
	Value() any
	String() string
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

func (kv KeyValue) String() string {
	return stringify(kv.value)
}

func KV(key any, value any) KeyValuer {
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
	case stringer:
		return s.String()
	case string:
		return s
	case nil:
		return "<nil>"
	}
	return reflect.TypeOf(v).String()
}

type stringer interface {
	String() string
}
