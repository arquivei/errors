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

func TestWith(t *testing.T) {
	err := errors.With(nil, "key 1", "value 1")
	if err != nil {
		t.Error("expected nil, got", err)
	}

	rootErr := errors.New("some error")
	err = errors.With(rootErr)

	if err != rootErr {
		t.Error("expected", rootErr, "got", err)
	}

	err = errors.With(rootErr, "key 1", "value 1")

	if _, ok := err.(errors.Error); !ok {
		t.Error("expected errors.Error, got", err)
	}
}
