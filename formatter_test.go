package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestFormatter(t *testing.T) {
	formatter := errors.Formatter(func(err error) string {
		return "formatted error: " + err.Error()
	})

	if formatter == nil {
		t.Error("expected non-nil formatter")
	}
	if formatter.String() != "<ErrorFormatter>" {
		t.Errorf("expected '<ErrorFormatter>', got '%s'", formatter.String())
	}
	if formatter.Key() == nil {
		t.Error("expected non-nil key")
	}
	if formatter.Value() == nil {
		t.Error("expected non-nil value")
	}
}

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
	expectedErrorMsg := `op 2: op 1: [input] (BAD_REQUEST) some error {map=map[key 1:1 key 2:2], slice=[a b c], int=2, str=value}`
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

func TestRootErrorFormatter(t *testing.T) {
	rootErr := errors.New("root error")
	err := errors.With(
		rootErr,
		errors.RootErrorFormatter,
		errors.SeverityInput,
		errors.Code("BAD_REQUEST"),
		errors.KV("key1", "value1"),
	)

	if err.Error() != "root error" {
		t.Errorf("expected 'root error', got '%s'", err.Error())
	}
}
func TestRootErrorKVFormatter(t *testing.T) {
	rootErr := errors.New("root error")
	err := errors.With(
		rootErr,
		errors.RootErrorKVFormatter,
		errors.SeverityInput,
		errors.Code("BAD_REQUEST"),
		errors.KV("key1", "value1"),
	)

	expected := "root error {key1=value1}"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}
