package errors

import (
	"runtime"
	"strings"
)

// CodePanic is set when a panic occurs in the DontPanic function.
var CodePanic Code = "PANIC"

// DontPanic executes the provided function and recovers from any panic that occurs.
// It returns an error containing the panic information if a panic occurs,
// else it returns nil or the error returned by func() error.
// If a panic occurs, it sets:
// - Code to CodePanic
// - Op to the operation where the panic occurred
// - Severity to SeverityFatal
func DontPanic[F func() | func() error](fn F) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = With(Errorf("panic: %v", r), getPanicOp(), CodePanic, SeverityFatal)
		}
	}()
	switch f := any(fn).(type) {
	case func():
		f()
	case func() error:
		err = f()
	}
	return
}

func getPanicOp() (op Op) {
	// find the panic operation in the stack trace
	callers := make([]uintptr, 10)
	runtime.Callers(1, callers)
	pcIdx := findPcAfterPanic(callers)

	// If we can't find the panic operation, return an unknown operation
	if pcIdx == -1 {
		return opUnknownFunction
	}

	return getCallerOp(callers[pcIdx], true)
}

func findPcAfterPanic(callers []uintptr) int {
	found := -1

	// Look for the runtime.gopanic function in the stack trace
	for i, pc := range callers {
		if runtime.FuncForPC(pc).Name() == "runtime.gopanic" {
			found = i + 1
		}
	}
	if found == -1 {
		return -1 // No panic found
	}

	// Skip runtime functions in the stack trace to find the actual panic location
	for i := found; i < len(callers); i++ {
		if !strings.HasPrefix(runtime.FuncForPC(callers[i]).Name(), "runtime.") {
			break
		}
		found = i + 1
	}

	return found
}
