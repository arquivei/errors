package main

import (
	"flag"
	"fmt"

	"github.com/arquivei/errors"
)

type name string

func (n name) String() string {
	return string(n)
}

type greeter[T fmt.Stringer] struct {
	name T
}

func (g greeter[T]) sayHello() error {
	if g.name.String() != "" {
		fmt.Printf("Hello, %s!\n", g.name)
		return nil
	}
	return errors.With(fmt.Errorf("name cannot be empty"),
		errors.SeverityInput,
		errors.Code("BAD_REQUEST"),
		errors.KV("context1", "value1"),
		errors.KV("context2", "value2"),
	)
}

func greetings(n name) error {
	err := greeter[name]{name: n}.sayHello()
	if err != nil {
		// This is to show that the Op will not be added twice and context1 will be overridden
		err = errors.With(err, errors.KV("context1", "this will be overridden"))

		return errors.With(err,
			errors.SeverityRuntime,
			errors.Code("RUNTIME_ERROR"),
			errors.KV("context1", "value1"),
			errors.KV("context2", "value2"),
		)
	}
	fmt.Println("Greetings sent successfully!")
	return nil
}

func doGreetings(n name) (err error) {
	err = greetings(n)
	if err != nil {
		err = errors.With(err)
	}
	// Wrapping in an anonymous function to simulate a call stack without caller name
	defer func() {
		if err != nil {
			err = errors.With(err,
				errors.SeverityFatal,
				errors.KV("context3", "value3"),
			)
		}
	}()

	return
}

func thisWillPanic() {
	panic("this will panic")
}

func main() {
	flag.Parse()
	n := ""
	if flag.NArg() > 0 {
		n = flag.Arg(0)
	}
	err := doGreetings(name(n))
	err = errors.With(err, errors.Op("customOpExample"))

	fmt.Println("Default formatter ==>", err)

	fmt.Println("Root error formatter ==>", errors.With(err, errors.RootErrorFormatter))
	fmt.Println("Root error formatter with KV ==>", errors.With(err, errors.RootErrorKVFormatter))

	err = errors.With(err, errors.Formatter(func(err error) string {
		return fmt.Sprintf("formatted error: %s", errors.GetRootError(err))
	}))

	fmt.Println("Custom formatter ==>", err)

	err = errors.DontPanic(func() {
		panic("hello anonymous function panic")
	})

	fmt.Println("DontPanic (anonymous) ==>", err)

	err = errors.DontPanic(thisWillPanic)
	fmt.Println("DontPanic (named) ==>", err)
}
