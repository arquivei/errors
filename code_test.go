package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestGetCode(t *testing.T) {
	err := errors.New("some error")
	if errors.GetCode(err) != errors.CodeUnset {
		t.Error("expected NoCode, got", errors.GetCode(err))
	}

	err = errors.With(err, errors.Code("code 1"))
	if errors.GetCode(err) != errors.Code("code 1") {
		t.Error("expected code 1, got", errors.GetCode(err))
	}

	err = errors.With(err, errors.Code("code 2"))
	if errors.GetCode(err) != errors.Code("code 2") {
		t.Error("expected code 2, got", errors.GetCode(err))
	}
}
