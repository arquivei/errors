package errors

import (
	"strings"
)

// Op represents an operation in the error stack.
type Op string

var _ KeyValuer = Op("")

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
