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

type Index interface{}
type Series map[Index]C

func (s Series) Type() map[Type]int {
	ret := map[Type]int{}
	for _, v := range s {
		if _, ok := ret[v.Type()]; !ok {
			ret[v.Type()] = 0
		}
		ret[v.Type()]++
	}
	return ret
}
