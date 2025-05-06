package errors

import (
	"reflect"
)

func With(err error, args ...any) error {
	if err == nil {
		return nil
	}

	keyvalues := argsToKeyValues(args)
	for _, keyval := range keyvalues {
		if !reflect.TypeOf(keyval.Key()).Comparable() {
			panic("key is not comparable")
		}
		// err = Error{err: err, KeyValue: KeyValue{key: keyval.Key(), value: keyval.Value()}}
		err = Error{err: err, keyval: keyval}
	}

	return err
}

func argsToKeyValues(args []any) []KeyValuer {
	var (
		keyvalue  KeyValuer
		keyvalues []KeyValuer
	)
	for len(args) > 0 {
		keyvalue, args = cutKV(args)
		keyvalues = append(keyvalues, keyvalue)
	}

	return keyvalues
}

func cutKV(args []any) (KeyValuer, []any) {
	switch head := args[0].(type) {
	case KeyValuer:
		return head, args[1:]
	default:
		if len(args) < 2 {
			panic("invalid number of args")
		}
		return KV(head, args[1]), args[2:]
	}
}
