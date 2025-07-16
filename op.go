package errors

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	_ KeyValuer = Op("")

	// VerboseOpOnAnonymousFunctions determines whether the Op should include file and line information
	// for anonymous functions.
	// If set to true, the Op will include the line number where the anonymous function was defined.
	VerboseOpOnAnonymousFunctions = true

	// NoOp is a special Op that indicates no operation is associated with the error.
	// It can e used to disable the automatic addition of Op to errors for a specific call to With.
	NoOp = Op("no-op")
)

// Op represents an operation in the error stack. An Op is a string that describes an operation,
// such as a function call or a method invocation, that is associated with an error.
// It implements the KeyValuer interface, allowing it to be used as a key-value pair in errors.
type Op string

func (op Op) Key() any {
	return opKey{}
}

func (op Op) Value() any {
	return op
}

func (op Op) String() string {
	return string(op)
}

type opKey struct{}

// GetOpStack retrieves the operation stack from an error.
// It returns a string representation of the operations in the stack,
// formatted as "op1: op2: ...", where each operation is separated by ": ".
// If no operations are found, it returns an empty string.
func GetOpStack(err error) string {
	ops := Values(err, opKey{})
	if len(ops) == 0 {
		return ""
	}

	sb := strings.Builder{}
	sb.Grow(32)

	sb.WriteString(stringify(ops[0]))
	for _, op := range ops[1:] {
		sb.WriteString(": ")
		sb.WriteString(stringify(op))
	}

	return sb.String()
}

// opUnknownFunction is used when the function name cannot be determined.
const opUnknownFunction Op = "<unknown function>"

// getCallerOp retrieves the operation for the caller at the given program counter (pc).
func getCallerOp(pc uintptr, alwaysIncludeLocation bool) Op {
	funcForPc := runtime.FuncForPC(pc)
	if funcForPc == nil {
		return opUnknownFunction
	}

	funcName := funcForPc.Name()
	funcName = discardPackagePath(funcName)

	if alwaysIncludeLocation || (VerboseOpOnAnonymousFunctions && isAnonymousFunction(funcName)) {
		file, line := funcForPc.FileLine(pc)
		return Op(funcNameWithLocation(funcName, file, line))
	}

	return Op(funcName)
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

// funcNameWithLocation formats the function name with its file and line number.
func funcNameWithLocation(funcName, file string, line int) string {
	const lineNumberChars = 32 // " (<file>:<line>)"
	var sb strings.Builder
	sb.Grow(len(funcName) + lineNumberChars)
	sb.WriteString(funcName)
	sb.WriteString(" (")
	sb.WriteString(filepath.Base(file))
	sb.WriteString(":")
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
