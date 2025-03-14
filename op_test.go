package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestGetOpStack(t *testing.T) {
	err := errors.New("some error")
	if errors.GetOpStack(err) != "" {
		t.Error("expected empty string, got", errors.GetOpStack(err))
	}

	// Showing that op can be passed in various ways
	err = errors.With(err, errors.Op("op 1"))
	err = errors.With(err, errors.Op("op 2"))
	const op3 errors.Op = "op 3"
	err = errors.With(err, op3.Key(), op3.Value())
	op4 := errors.Op("op 4")
	err = errors.With(err, op4)

	expected := "op 4: op 3: op 2: op 1"
	actual := errors.GetOpStack(err)

	if expected != actual {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}
