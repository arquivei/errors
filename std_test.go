package errors

import "testing"

func TestUnwrap(t *testing.T) {
	rootErr := New("root error")
	wrappedErr := With(rootErr)

	if _, ok := wrappedErr.(Error); !ok {
		t.Fatal("expected wrappedErr to be of type Error")
	}

	if unwrapped := Unwrap(wrappedErr); unwrapped != rootErr {
		t.Errorf("expected %v, got %v", rootErr, unwrapped)
	}
}

func TestIs(t *testing.T) {
	rootErr := New("root error")
	err := With(rootErr, KV("key", "value"))
	err = With(err, KV("key2", "value2"))

	if _, ok := err.(Error); !ok {
		t.Fatal("expected wrappedErr to be of type Error")
	}

	if !Is(err, rootErr) {
		t.Errorf("expected error to match rootErr, got %v", err)
	}
}

func TestAs(t *testing.T) {
	rootErr := New("root error")
	err := With(rootErr, KV("key", "value"))

	var e Error
	if !As(err, &e) {
		t.Fatal("expected error to be of type Error")
	}
}

func TestNew(t *testing.T) {
	err := New("test error")
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if err.Error() != "test error" {
		t.Errorf("expected 'test error', got %s", err.Error())
	}
}

func TestErrorf(t *testing.T) {
	err := Errorf("test error: %s", "details")
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}

	if err.Error() != "test error: details" {
		t.Errorf("expected 'test error: details', got %s", err.Error())
	}
}

func TestJoin(t *testing.T) {
	err1 := New("first error")
	err2 := New("second error")

	err := Join(err1, err2)
	if err == nil {
		t.Fatal("expected non-nil error, got nil")
	}
	if err.Error() != "first error\nsecond error" {
		t.Errorf("expected 'first error\nsecond error', got %s", err.Error())
	}
}
