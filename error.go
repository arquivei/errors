package errors

var _ error = Error{}

// Error is an error that wraps another error and adds a key-value pair.
type Error struct {
	err    error
	keyval KeyValuer
}

// Error returns the error message formatted by the formatter associated with the error.
// If no formatter is set, it uses the default formatter.
func (e Error) Error() string {
	return GetFormatter(e)(e)
}

// Unwrap returns the underlying error wrapped by this Error.
func (e Error) Unwrap() error {
	return e.err
}
