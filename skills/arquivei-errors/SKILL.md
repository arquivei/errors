---
name: arquivei-errors
description: Instructions for creating, handling, and extracting errors in Go projects using the github.com/arquivei/errors package. Use this skill whenever writing or modifying Go code that imports github.com/arquivei/errors, or when asked to fix error handling in a project using this library.
---

# Arquivei Errors Package

This skill provides instructions on how to correctly use the `github.com/arquivei/errors` package to create meaningful errors with injected key-value pairs in Go.

## Core Rule
ALWAYS use `errors.With(err, ...)` or `errors.Errorf(...)` from `github.com/arquivei/errors` instead of standard `fmt.Errorf` or standard `errors` package when manipulating and wrapping errors.

## Key Concepts

The package combines an `error` with a Key-Value pair (conceptually similar to `context.Context`).

### 1. Wrapping Errors
Wrap existing errors with contextual information using `errors.With`. Multiple properties can be chained together.

```go
import "github.com/arquivei/errors"

func doStuff() error {
    err := someOtherFunction()
    if err != nil {
        return errors.With(err,
            errors.Code("DO_STUFF_FAILED"),
            errors.SeverityRuntime,
            errors.KV("extra_context", "my_value"),
        )
    }
    return nil
}
```

### 2. Built-in KeyValuers

Always inject these built-in values when they add context to an error:
- **`errors.Op(string)`**: The operation or function running. `errors.With()` adds the function name automatically, but you can override it explicitly using `errors.Op("customOp")` or disable it with `errors.NoOp`.
- **`errors.Severity*`**: Use to indicate error severity.
  - `errors.SeverityInput` (Bad input, do not retry, similar to HTTP 400).
  - `errors.SeverityRuntime` (External failure, DB timeout, etc. Potentially retryable).
  - `errors.SeverityFatal` (Code not ready to handle this, notify developers).
- **`errors.Code(string)`**: A simple string code to differentiate errors (e.g., `errors.Code("BAD_REQUEST")`). Often used in `switch` statements by the caller.
- **`errors.KV(key, value)`**: Arbitrary key-value pair for extra context. Helps with logging (prints as `{key=value, ...}`).

#### Performance Tip
For performance-critical code (hot paths), you can avoid the overhead of runtime stack capturing by:
1. Passing `errors.NoOp` to individual `errors.With` calls.
2. Setting `errors.AutomaticallyAddOp = false` globally in your application's `init()` or `main()`.

### 3. Extracting Values
Use these helper functions to retrieve injected values from an error:

```go
code := errors.GetCode(err)                 // Returns string
severity := errors.GetSeverity(err)         // Returns errors.Severity
val := errors.Value(err, mykey)             // Returns any
valStr := errors.ValueT[string](err, mykey) // Returns specific type via Generics
```

### 4. Option Pattern for Complex Errors
If you need to build the error context conditionally, use a slice of `KeyValuer`:

```go
errorOpts := []errors.KeyValuer{
    errors.Code("RUNTIME_ERROR"),
}
if condition {
    errorOpts = append(errorOpts, errors.KV("key1", "value1"))
    errorOpts = append(errorOpts, errors.SeverityInput)
}
return errors.With(err, errorOpts...)
```

### 5. Error Formatting & Interoperability
The package is designed to work seamlessly with Go's standard ecosystem.

#### Standard Formatting (fmt.Printf)
The error implements `fmt.Formatter`. 
- Use `%+v` to print the **full error context** (Operation stack, Severity, Code, and KV pairs).
- Use `%v` or `.Error()` to print only the **root error message**.

```go
fmt.Printf("%+v", err) // Outputs: op2: op1: [runtime] (DB_ERR) db timeout {table=users}
```

#### Structured Logging (slog)
The error natively implements `slog.LogValuer`. Passing the error to an `slog` logger will automatically extract all context into structured JSON/text fields:

```go
import "log/slog"

// Fields like 'code', 'severity', 'op' and all custom KVs will be extracted automatically.
logger.Error("failed", "error", err) 
```

#### Manual Formatting
If you need the full string representation manually, use `errors.Format(err)`:
```go
fullMsg := errors.Format(err)
```

### 6. Panic Recovery
To safely catch panics and transform them into `arquivei/errors`, use `DontPanic`:
```go
err = errors.DontPanic(func() {
    panic("something went wrong")
})
```