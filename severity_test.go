package errors_test

import (
	"testing"

	"github.com/arquivei/errors"
)

func TestGetSeverity(t *testing.T) {
	err := errors.New("some error")
	if errors.GetSeverity(err) != errors.SeverityUnset {
		t.Error("expected Unset, got", errors.GetSeverity(err))
	}

	err = errors.With(err, errors.SeverityRuntime)
	if errors.GetSeverity(err) != errors.SeverityRuntime {
		t.Error("expected Runtime, got", errors.GetSeverity(err))
	}

	err = errors.With(err, errors.Severity("custom"))
	if errors.GetSeverity(err) != errors.Severity("custom") {
		t.Error("expected custom, got", errors.GetSeverity(err))
	}
}
