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
		errors.KV("str", "value"),
		errors.KV("int", 2),
		errors.KV("slice", []string{"a", "b", "c"}),
		errors.KV("map", map[string]int{"key 1": 1, "key 2": 2}),
	)
	expectedErrorMsg := `op 2: op 1: [input] (BAD_REQUEST) some error {map: map[key 1:1 key 2:2], slice: [a b c], int: 2, str: value}`
	if err.Error() != expectedErrorMsg {
		t.Errorf("expected '%s', got '%s'", expectedErrorMsg, err.Error())
	}

	err = errors.With(err, errors.Formatter(func(err error) string {
		return "custom formatter"
	}))

	if err.Error() != "custom formatter" {
		t.Error("expected custom formatter, got", err.Error())
	}
}
