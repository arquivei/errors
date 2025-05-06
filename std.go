// Description: This file contains the standard error functions so that it's not necessary
// import the errors package if this package is already imported.
package errors

import (
	"errors"
	"fmt"
)

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func New(str string) error {
	return errors.New(str)
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
