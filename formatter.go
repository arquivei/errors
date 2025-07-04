package errors

import "strings"

var (
	_ KeyValuer = Formatter(nil)

	// DefaultFormatter is the default formatter used when no custom formatter is set.
	// It defaults to `FullFormatter`, which provides a comprehensive view of the error.
	DefaultFormatter = FullFormater
)

type formatterKey struct{}

// Formatter is a function type that formats an error into a string representation.
// It can be used to customize how errors are displayed, including their severity, code, and context.
// It is typically used in conjunction with the `errors.With` function to attach a custom formatter to an error.
type Formatter func(err error) string

func (Formatter) Key() any {
	return formatterKey{}
}

func (f Formatter) Value() any {
	return f
}

func (Formatter) String() string {
	return "<ErrorFormatter>"
}

// Format formats the error using the custom formatter associated with it.
// If no custom formatter is set, it uses the default formatter, which is `FullFormatter`.
// If the error is nil, it returns an empty string.
func Format(err error) string {
	// If the error is nil, return an empty string.
	if err == nil {
		return ""
	}

	return GetFormatter(err)(err)
}

// GetFormatter retrieves the custom formatter associated with the error.
// If no custom formatter is set, it returns the default formatter, which is `FullFormatter`.
func GetFormatter(err error) Formatter {
	if formatter, ok := Value(err, formatterKey{}).(Formatter); ok {
		return formatter
	}

	return DefaultFormatter
}

// FullFormater formats the error with its operation stack, severity, code, and key-value pairs.
// It provides a comprehensive view of the error, including its context and any additional information that has been attached to it.
// The format is as follows:
// operation2: ... operation1: [severity] (code) root error message {key1: value1, key2: value2, ...}
var FullFormater Formatter = func(err error) string {
	sb := strings.Builder{}
	sb.Grow(32)

	writeOpStack(&sb, GetOpStack(err))
	writeSeverity(&sb, GetSeverity(err))
	writeCode(&sb, GetCode(err))

	sb.WriteString(err.Error())
	writeKV(&sb, ValueAllSlice(err))

	return sb.String()
}

// KVFormatter formats the error's message along with its key-value pairs.
var KVFormatter Formatter = func(err error) string {
	sb := strings.Builder{}
	sb.Grow(32)

	sb.WriteString(err.Error())
	writeKV(&sb, ValueAllSlice(err))

	return sb.String()
}

func writeCode(sb *strings.Builder, code Code) {
	if code == CodeUnset {
		return
	}
	sb.WriteString("(")
	sb.WriteString(code.String())
	sb.WriteString(") ")
}

func writeSeverity(sb *strings.Builder, severity Severity) {
	if severity == SeverityUnset {
		return
	}
	sb.WriteString("[")
	sb.WriteString(severity.String())
	sb.WriteString("] ")
}

func writeOpStack(sb *strings.Builder, ops string) {
	if ops == "" {
		return
	}
	sb.WriteString(ops)
	sb.WriteString(": ")
}

func writeKV(sb *strings.Builder, kvs []KeyValuer) {
	if len(kvs) <= 0 {
		return
	}
	sb.WriteString(" {")
	shouldAddComma := false
	for _, kv := range kvs {
		if shouldAddComma {
			sb.WriteString(", ")
		}
		sb.WriteString(stringify(kv.Key()))
		sb.WriteString("=")
		sb.WriteString(stringify(kv.Value()))

		shouldAddComma = true
	}

	sb.WriteString("}")
}
