package errors

import "errors"

// GetRootError returns the root error of the error chain.
func GetRootError(err error) error {
	for unwrapped := errors.Unwrap(err); unwrapped != nil; unwrapped = errors.Unwrap(unwrapped) {
		err = unwrapped
	}

	return err
}
