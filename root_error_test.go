package errors_test

import (
	"fmt"
	"testing"

	"github.com/arquivei/errors"
)

func TestGetRootError(t *testing.T) {
	rootErr := errors.New("root error")

	err := errors.With(rootErr, errors.Op("op1"), errors.SeverityInput, errors.KV("key", "value"))
	err = fmt.Errorf("wrapped by fmt: %w", err)

	got := errors.GetRootError(err)
	if got != rootErr {
		t.Errorf("expected %v, got %v", rootErr, got)
	}
}
