package errors

type Code string

var _ KeyValuer = Code("")

const NoCode Code = ""

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

func GetCode(err error) Code {
	val := Value(err, codeKey{})
	if code, ok := val.(Code); ok {
		return code
	}

	return NoCode
}
