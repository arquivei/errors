package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestWithNoKeyValues(t *testing.T) {
	rootErr := errors.New("some error")
	err := errors.With(rootErr)

	if err != rootErr {
		t.Error("expected", rootErr, "got", err)
	}
}

func TestWithNoError(t *testing.T) {
	err := errors.With(nil)
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestWith(t *testing.T) {
	// receiving a single keyvalue will return a new error
	err := errors.With(errors.New("some error"), errors.KV("key", "value"))
	if _, ok := err.(errors.Error); !ok {
		t.Error("expected errors.Error, got", err)
	}
}
