package errors

import (
	"reflect"
	"runtime"
)

var (
	// AutomaticallyAddOp determines whether the Op should be automatically added
	// to errors when using the With function.
	// If set to true, the Op will be automatically added to errors that do not already have an Op.
	AutomaticallyAddOp = true

	// ErrKeyNotComparable defines an error that is returned when a key in With is not comparable.
	ErrKeyNotComparable = New("key is not comparable")
)

// With adds key-value pairs to an error, allowing for additional context.
func With(err error, keyvalues ...KeyValuer) error {
	if err == nil {
		return nil
	}

	shouldAddAutomaticOp := AutomaticallyAddOp

	for _, keyval := range keyvalues {
		if !reflect.TypeOf(keyval.Key()).Comparable() {
			panic(ErrKeyNotComparable)
		}
		err = Error{err: err, keyval: keyval}
		if keyval.Key() == (opKey{}) {
			shouldAddAutomaticOp = false
		}
	}

	if shouldAddAutomaticOp {
		return withAutomaticOp(err)
	}

	return err
}

func withAutomaticOp(err error) error {
	pc, _, _, _ := runtime.Caller(2)
	op := getCallerOp(pc, false)
	if ValueT[Op](err, opKey{}) == op {
		return err // Op already exists, no need to add it again
	}
	return With(err, op)
}
