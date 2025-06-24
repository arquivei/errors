package errors

import (
	"testing"
)

func TestDontPanic(t *testing.T) {
	// This test is designed to ensure that the package does not panic
	// when using the standard error functions.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("unexpected panic: %v", r)
		}
	}()

	fnNoPanic := func() error {
		return nil
	}

	fnNoPanic2 := func() {}

	err := DontPanic(fnNoPanic)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	err = DontPanic(fnNoPanic2)
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	fn := func() error {
		panic(New("my panic"))
	}

	err = DontPanic(fn)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	fn2 := func() {
		panic("another panic")
	}

	err = DontPanic(fn2)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

}

func Test_getPanicOp(t *testing.T) {
	op := getPanicOp()
	if op != opUnknownFunction {
		t.Errorf("expected unknown function, got %v", op)
	}

	err := DontPanic(func() {
		s := []string{}
		s[1] = "this should panic"
	})

	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	expectedOp := Op("errors.Test_getPanicOp.func1 (dont_panic_test.go:59)")
	if op := ValueT[Op](err, opKey{}); op != expectedOp {
		t.Errorf("expected operation '%s', got ;'%s'", expectedOp, op)
	}
}
