package errors_test

import (
	"fmt"
	"testing"

	"github.com/arquivei/errors"
)

func TestError_Format_Compliance(t *testing.T) {
	err := errors.Errorf("base error")
	// Add some context to check %+v later
	err = errors.With(err, errors.Op("MyOp"), errors.Code("TEST_CODE"))

	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "verb %s",
			format:   "%s",
			expected: "base error",
		},
		{
			name:     "verb %q",
			format:   "%q",
			expected: `"base error"`,
		},
		{
			name:     "verb %v",
			format:   "%v",
			expected: "base error",
		},
		{
			name:     "verb %+v",
			format:   "%+v",
			expected: "MyOp: (TEST_CODE) base error",
		},
		{
			name:     "verb %s with width",
			format:   "%15s",
			expected: "     base error",
		},
		{
			name:     "verb %q with width",
			format:   "%15q",
			expected: `   "base error"`,
		},
		{
			name:     "verb %s with precision",
			format:   "%.4s",
			expected: "base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fmt.Sprintf(tt.format, err)
			if got != tt.expected {
				t.Errorf("fmt.Sprintf(%q, err) = %v, want %v", tt.format, got, tt.expected)
			}
		})
	}
}
