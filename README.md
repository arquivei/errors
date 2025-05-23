# errors

A golang package to create useful errors.

## Using the `errors` package

Import the package:

``` go
import (
	"github.com/arquivei/errors"
)
```

Use `errors.With()`:

``` go
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
```

