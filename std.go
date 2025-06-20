// Description: This file contains the standard error functions so that it's not necessary
// import the errors package if this package is already imported.
package errors

import (
	"errors"
	"fmt"
)

// Unwrap returns the next error in the chain, if any.
// This calls the standard library's errors.Unwrap function
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Is checks if the target error is present in the error chain.
// This calls the standard library's errors.Is function.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As checks if the error can be cast to the target type.
// This calls the standard library's errors.As function.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// New creates a new error with the given string message.
// This calls the standard library's errors.New function.
func New(str string) error {
	return errors.New(str)
}

// Errorf formats an error message according to a format specifier and returns it as an error.
// This calls the standard library's fmt.Errorf function.
func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// Join combines multiple errors into a single error.
// This calls the standard library's errors.Join function.
func Join(errs ...error) error {
	return errors.Join(errs...)
}
