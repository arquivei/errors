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
	funcName = discardPackagePath(funcName)

	if VerboseOpOnAnonymousFunctions && isAnonymousFunction(funcName) {
		return funcNameWithLineNumber(funcName, line)
	}

	return funcName
}

// isAnonymousFunction checks if the function name indicates an anonymous function.
// It does this by checking if the function name contains a dot (.) and starts with "func" after the last dot.
// If there isn't a dot in the function name, it checks if the function name starts with "func".
func isAnonymousFunction(funcName string) bool {
	idx := strings.LastIndex(funcName, ".")
	if idx > 0 {
		return strings.HasPrefix(funcName[idx+1:], "func")
	}
	return strings.HasPrefix(funcName, "func")
}

func funcNameWithLineNumber(funcName string, line int) string {
	const lineNumberChars = 12 // " (line XXXX)" has 12 characters
	var sb strings.Builder
	sb.Grow(len(funcName) + lineNumberChars)
	sb.WriteString(funcName)
	sb.WriteString(" (line ")
	sb.WriteString(strconv.Itoa(line))
	sb.WriteString(")")
	return sb.String()
}

func discardPackagePath(s string) string {
	if idx := strings.LastIndex(s, "/"); idx >= 0 {
		return s[idx+1:]
	}
	return s
}
