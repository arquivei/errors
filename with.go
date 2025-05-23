package errors

import (
	"reflect"
)

func With(err error, keyvalues ...KeyValuer) error {
	if err == nil {
		return nil
	}

	for _, keyval := range keyvalues {
		if !reflect.TypeOf(keyval.Key()).Comparable() {
			panic("key is not comparable")
		}
		err = Error{err: err, keyval: keyval}
	}

	return err
}
