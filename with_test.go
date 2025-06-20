package errors

import (
	"runtime"
	"testing"
)

func TestWithNoKeyValues(t *testing.T) {
	rootErr := New("some error")

	err := With(rootErr)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	var e Error
	if !As(err, &e) {
		t.Fatalf("expected error to be of type Error, got %T", err)
	}
}

func TestWithNoError(t *testing.T) {
	err := With(nil)
	if err != nil {
		t.Error("expected nil, got", err)
	}
}

func TestWith(t *testing.T) {
	// With should return an Error type
	err := With(New("some error"), KV("key", "value"))
	if _, ok := err.(Error); !ok {
		t.Error("expected Error, got", err)
	}

	// Check if the error message is formatted correctly
	err = With(New("some error"), KV("key", "value"))
	if err.Error() != "errors.TestWith: some error {key=value}" {
		t.Error("expected 'errors.TestWith: some error {key: value}', got", err.Error())
	}

	err = func() error {
		return With(New("some error"))
	}()
	expected := "errors.TestWith.func1 (with_test.go:42): some error"
	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}

	// Force an anonymous function name
	err = With(New("some error"), Op("customOp"))
	if err.Error() != "customOp: some error" {
		t.Error("expected 'customOp: some error', got", err.Error())
	}

	t.Run("UncomparableKey", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic, got nil")
				t.FailNow()
			}
			if err, ok := r.(error); ok {
				if err.Error() != "key is not comparable" {
					t.Errorf("expected 'key is not comparable', got %s", err.Error())
				}
			} else {
				t.Errorf("expected panic with error, got %v", r)
			}
		}()
		_ = With(New("some error"), KV(func() {}, "value"))
	})
}

func Test_isAnonymousFunction(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
		{"func1", true},
		{"package.func1", true},
		{"package.func2", true},
		{"package.SayHello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAnonymousFunction(tt.name); got != tt.want {
				t.Errorf("isAnonymousFunction(%q) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_discardPackagePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"func1", "func1"},
		{"package.func1", "package.func1"},
		{"path/package.SayHello", "package.SayHello"},
		{"path/SayHello", "SayHello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := discardPackagePath(tt.name); got != tt.want {
				t.Errorf("discardPackagePath(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func Test_getCallerOp(t *testing.T) {
	// Caso normal
	pc, _, _, _ := runtime.Caller(0)
	result := getCallerOp(pc, false)
	if result == opUnknownFunction {
		t.Error("Expected function name, got unknown")
	}

	// Caso !ok
	done := make(chan Op)
	go func() {
		pc, _, _, _ := runtime.Caller(3)
		done <- getCallerOp(pc, false)
	}()
	if result := <-done; result != "<unknown function>" {
		t.Errorf("Expected <unknown function>, got %s", result)
	}
}
