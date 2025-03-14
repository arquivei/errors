package errors

type severityKey struct{}

// Severity is the error severity. It's used to classify errors in groups to be easily handled by the code. For example,
// a retry layer should be only checking for Runtime errors to retry. Or in an HTTP layer, errors of input type are always
// returned a 400 status.
type Severity string

var _ KeyValuer = Severity("")

const (
	// SeverityUnset indicates the severity was not set
	SeverityUnset = Severity("")
	// SeverityRuntime indicates the error is returned for an operation that should/could be executed again. For example, timeouts.
	SeverityRuntime = Severity("runtime")
	// SeverityFatal indicates the error is unrecoverable and the execution should stop, or not being retried.
	SeverityFatal = Severity("fatal")
	// SeverityInput indicates  an expected, like a bad user input/request. For example, an invalid email.
	SeverityInput = Severity("input")
)

func (s Severity) String() string {
	return string(s)
}

func (s Severity) Key() any {
	return severityKey{}
}

func (s Severity) Value() any {
	return s
}

// Value returns the severity of the error. If there is not severity, Unset is returned.
func GetSeverity(err error) Severity {
	val := Value(err, severityKey{})
	if severity, ok := val.(Severity); ok {
		return severity
	}

	return SeverityUnset
}
