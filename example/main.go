package main

import (
	"fmt"

	"github.com/arquivei/errors"
)

func doStuff() error {
	const op = errors.Op("doStuff")
	err := fmt.Errorf("some error")

	return errors.With(err,
		errors.SeverityRuntime, op,
		errors.Code("RUNTIME_ERROR"),
		errors.KV("context1", "value1"),
		errors.KV("context2", "value2"),
	)
}

func doMoreStuff() error {
	err := doStuff()
	return errors.With(err,
		errors.SeverityFatal,
		errors.Op("doMoreStuff"),
		errors.KV("context3", "value3"),
	)
}

func main() {

	err := doMoreStuff()
	fmt.Println(err)

	err = errors.With(err, errors.Formatter(func(err error) string {
		return fmt.Sprintf("formatted error: %s", errors.GetRootError(err))
	}))

	fmt.Println(err)
}
