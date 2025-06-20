package errors

import "testing"

type stringer struct{}

func (s stringer) String() string {
	return "stringer"
}

func Test_stringify(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"nil", nil, "<nil>"},
		{"string", "test", "test"},
		{"int", 42, "42"},
		{"float64", 3.14, "3.14"},
		{"bool", true, "true"},
		{"stringer", stringer{}, "stringer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringify(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
