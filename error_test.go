package errors

import "testing"

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

}
