package types

import "time"

type Type string
type Operator string

const (
	ADD Operator = "+"
	SUB          = "-"
	MUL          = "*"
	DIV          = "/"
	MOD          = "%"
)
const (
	NUMERIC Type = "numeric"
	TIME         = "time"
	STRING       = "string"
	NAN          = "nan"
	MULTI        = "multiple"
)

type C interface {
	Add(C) C
	Sub(C) C
	Mul(C) C
	Div(C) C
	Mod(C) C
	Great(C) bool
	Less(C) bool
	Equal(C) bool
	NotEqual(C) bool
	Type() Type
}

func NewC(i interface{}) C {
	var ret C
	switch i.(type) {
	case C:
		ret = i.(C)
	case float64:
		ret = Numeric(i.(float64))
	case int:
		ret = Numeric(float64(i.(int)))
	case time.Time:
		ret = Time(i.(time.Time))
	case string:
		ret = String(i.(string))
	case Numeric:
		ret = i.(Numeric)
	case String:
		ret = i.(String)
	default:
		ret = NewNan()
	}
	return ret
}
