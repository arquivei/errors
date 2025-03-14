package errors

import (
	"strings"
)

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
