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
	got := errors.GetFormatter(err)(err)
	if got != expectedErrorMsg {
		t.Errorf("expected '%s', got '%s'", expectedErrorMsg, got)
	}

	err = errors.With(err, errors.Formatter(func(err error) string {
		return "custom formatter"
	}))

	got = errors.GetFormatter(err)(err)
	if got != "custom formatter" {
		t.Error("expected custom formatter, got", err.Error())
	}
}

func TestRootErrorKVFormatter(t *testing.T) {
	rootErr := errors.New("root error")
	err := errors.With(
		rootErr,
		errors.KVFormatter,
		errors.SeverityInput,
		errors.Code("BAD_REQUEST"),
		errors.KV("key1", "value1"),
	)

	expected := "root error {key1=value1}"
	got := errors.Format(err)
	if got != expected {
		t.Errorf("expected '%s', got '%s'", expected, got)
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		err  error
		want string
	}{
		{
			name: "nil error",
			err:  nil,
			want: "",
		},
		{
			name: "simple error",
			err:  errors.New("simple error"),
			want: "simple error",
		},
		{
			name: "error with operation and severity",
			err: errors.With(
				errors.New("operation error"),
				errors.Op("operation"),
				errors.SeverityInput,
			),
			want: "operation: [input] operation error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Format(tt.err)
			if got != tt.want {
				t.Errorf("Format(err) = %v, want %v", got, tt.want)
			}
		})
	}
}
