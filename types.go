package gopandas

type Type string

const (
	NUMERIC Type = "numeric"
	TIME         = "time"
	STRING       = "string"
	NAN          = "nan"
	MULTI        = "multiple"
)

type C interface {
	Add(C) C
	Mul(C) C
	Div(C) C
	Mod(C) C
	Great(C) bool
	Less(C) bool
	Equal(C) bool
	NotEqual(C) bool
	Type() Type
}
