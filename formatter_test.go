package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestGetFormatter(t *testing.T) {
	err := errors.New("some error")
	if errors.GetFormatter(err) == nil {
		t.Error("expected formatter, got nil")
	}

	err = errors.With(err,
		errors.Op("op 1"),
		errors.Op("op 2"),
		errors.SeverityInput,
		errors.Code("BAD_REQUEST"),
		errors.KV("key 1", "value 1"),
		errors.KV("key 2", "value 2"),
	)
	if err.Error() != "op 2: op 1: [input] (BAD_REQUEST) some error {key 2: value 2, key 1: value 1}" {
		t.Error("expected some error, got", err.Error())
	}

	err = errors.With(err, errors.Formatter(func(err error) string {
		return "custom formatter"
	}))

	if err.Error() != "custom formatter" {
		t.Error("expected custom formatter, got", err.Error())
	}
}
