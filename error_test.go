package errors

import (
	"fmt"
	"testing"
)

func TestErrorError(t *testing.T) {
	rootErr := New("root error")

	t.Run("Unwrap", func(t *testing.T) {
		var err error = Error{err: rootErr}
		if unwrapped := err.(Error).Unwrap(); unwrapped != rootErr {
			t.Errorf("expected %v, got %v", rootErr, unwrapped)
		}
	})

	t.Run("Error", func(t *testing.T) {
		var err error = Error{err: rootErr}
		expected := "root error"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("Nil Error", func(t *testing.T) {
		var err error = Error{err: nil}
		expected := "<root error is nil>"
		if err.Error() != expected {
			t.Errorf("expected %q, got %q", expected, err.Error())
		}
	})

	t.Run("mixed error chain", func(t *testing.T) {
		err := With(New("root error"), Op("op1"), SeverityInput, KV("key", "value"))
		err = fmt.Errorf("wrapped by fmt: %w", err)
		err = With(err, Op("op2"), SeverityFatal)

		expectedError := "wrapped by fmt: root error"
		if err.Error() != expectedError {
			t.Errorf("expected %q, got %q", expectedError, err.Error())
		}

		expectedFormatted := "op2: op1: [fatal] wrapped by fmt: root error {key=value}"
		if Format(err) != expectedFormatted {
			t.Errorf("expected %q, got %q", expectedFormatted, Format(err))
		}
	})
}
