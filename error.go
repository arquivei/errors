package errors

import (
	"fmt"
	"log/slog"
)

var _ error = Error{}

// Error is an error that wraps another error and adds a key-value pair.
type Error struct {
	err    error
	keyval KeyValuer
}

// Error returns the error message of the wrapped error.
func (e Error) Error() string {
	if e.err == nil {
		return "<root error is nil>"
	}

	return e.err.Error()
}

// Unwrap returns the underlying error wrapped by this Error.
func (e Error) Unwrap() error {
	return e.err
}

// Format implements the fmt.Formatter interface allowing for standard formatting of errors (e.g. fmt.Printf("%+v", err)).
// When formatted with the + flag, it injects the full rich context string via FullFormater,
// bypassing any custom formatters in the error chain to ensure all details are available for debugging.
func (e Error) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		_, _ = fmt.Fprint(s, FullFormater(e))
		return
	}
	_, _ = fmt.Fprintf(s, fmt.FormatString(s, verb), e.Error())
}

// LogValue implements the slog.LogValuer interface, allowing standard slog loggers
// to unpack the context stored in this error (Key-Values, Severity, Code, etc)
// into structured log attributes without requiring manual unwrapping.
func (e Error) LogValue() slog.Value {
	attrs := make([]slog.Attr, 0, 8)
	processed := make(map[any]struct{})

	// Extract standard fields for structured logging
	if code := GetCode(e); code != CodeUnset {
		attrs = append(attrs, slog.String("code", code.String()))
		processed[codeKey{}] = struct{}{}
	}
	if severity := GetSeverity(e); severity != SeverityUnset {
		attrs = append(attrs, slog.String("severity", severity.String()))
		processed[severityKey{}] = struct{}{}
	}
	if ops := GetOpStack(e); ops != "" {
		attrs = append(attrs, slog.String("op", ops))
		processed[opKey{}] = struct{}{}
	}

	// Add root error message
	attrs = append(attrs, slog.String("cause", e.Error()))

	// Iterate the error chain to collect other KVs
	for curr := error(e); curr != nil; curr = Unwrap(curr) {
		err, ok := curr.(Error)
		if !ok || err.keyval == nil {
			continue
		}
		key := err.keyval.Key()
		if _, exists := processed[key]; !exists {
			if key != (formatterKey{}) {
				attrs = append(attrs, slog.Any(stringify(key), err.keyval.Value()))
			}
			processed[key] = struct{}{}
		}
	}

	return slog.GroupValue(attrs...)
}
