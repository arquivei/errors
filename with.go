package errors

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

var (
	// AutomaticallyAddOp determines whether the Op should be automatically added
	// to errors when using the With function.
	// If set to true, the Op will be automatically added to errors that do not already have an Op.
	AutomaticallyAddOp = true
	// VerboseOpOnAnonymousFunctions determines whether the Op should include file and line information
	// for anonymous functions.
	// If set to true, the Op will include the line number where the anonymous function was defined.
	VerboseOpOnAnonymousFunctions = true
)

func With(err error, keyvalues ...KeyValuer) error {
	if err == nil {
		return nil
	}

	shouldAddAutomaticOp := AutomaticallyAddOp

	for _, keyval := range keyvalues {
		if !reflect.TypeOf(keyval.Key()).Comparable() {
			panic("key is not comparable")
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
	op := Op(getWithCaller())
	if ValueT[Op](err, opKey{}) == op {
		return err // Op already exists, no need to add it again
	}
	return With(err, Op(getWithCaller()))
}

func getWithCaller() string {
	pc, _, line, ok := runtime.Caller(3)
	if !ok {
		return "<unknown function>"
	}

	funcName := runtime.FuncForPC(pc).Name()
	if VerboseOpOnAnonymousFunctions && isAnonymousFunction(funcName) {
		funcName += " (line " + strconv.Itoa(line) + ")"
	}

	return funcName
}

// isAnonymousFunction checks if the function name indicates an anonymous function.
func isAnonymousFunction(funcName string) bool {
	parts := strings.Split(funcName, ".")
	return strings.HasPrefix(parts[len(parts)-1], "func")
}
