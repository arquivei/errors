package errors

import "strings"

type formatterKey struct{}

type Formatter func(err error) string

var _ KeyValuer = Formatter(nil)

func (Formatter) Key() any {
	return formatterKey{}
}

func (f Formatter) Value() any {
	return f
}

func (Formatter) String() string {
	return "<ErrorFormatter>"
}

func GetFormatter(err error) Formatter {
	if formatter, ok := Value(err, formatterKey{}).(Formatter); ok {
		return formatter
	}

	return defaultFormater
}

func defaultFormater(err error) string {
	sb := strings.Builder{}
	sb.Grow(32)
	ops := GetOpStack(err)
	if len(ops) > 0 {
		sb.WriteString(ops)
		sb.WriteString(": ")
	}

	severity := GetSeverity(err)
	if severity != SeverityUnset {
		sb.WriteString("[")
		sb.WriteString(severity.String())
		sb.WriteString("] ")
	}

	code := GetCode(err)
	if code != NoCode {
		sb.WriteString("(")
		sb.WriteString(code.String())
		sb.WriteString(") ")
	}

	sb.WriteString(GetRootError(err).Error())
	context := ValueMapOf(err, "")
	if len(context) > 0 {
		sb.WriteString(" {")
		shouldAddComma := false
		for k, v := range context {
			if shouldAddComma {
				sb.WriteString(", ")
			}
			sb.WriteString(k.(string))
			sb.WriteString(": ")
			sb.WriteString(v.(string))
			shouldAddComma = true
		}
		sb.WriteString("}")
	}

	return sb.String()
}
