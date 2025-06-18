package errors

type Code string

var _ KeyValuer = Code("")

// CodeUnset is the default value for Code, indicating no specific code is set.
const CodeUnset Code = ""

func (c Code) Key() any {
	return codeKey{}
}

func (c Code) Value() any {
	return c
}

func (c Code) String() string {
	return string(c)
}

type codeKey struct{}

// GetCode retrieves the Code from an error, returning CodeUnset if no code is set.
func GetCode(err error) Code {
	val := Value(err, codeKey{})
	if code, ok := val.(Code); ok {
		return code
	}

	return CodeUnset
}
