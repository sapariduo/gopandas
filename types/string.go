package types

import "fmt"

type String string

func (String) Type() Type {
	return STRING
}

func (s String) Add(c C) C {
	var ret C
	switch c.(type) {
	case String:
		ret = s + c.(String)
	default:
		ret = s
	}
	return ret
}

func (s String) Mul(c C) C {
	return s
}
func (s String) Div(c C) C {
	return s
}
func (s String) Mod(c C) C {
	return s
}

func (s String) Great(c C) bool {
	var ret bool

	switch c.(type) {
	case String:
		ret = s > c.(String)
	default:
		ret = false
	}
	return ret
}

func (s String) Less(c C) bool {
	var ret bool

	switch c.(type) {
	case String:
		ret = s < c.(String)
	default:
		ret = false
	}
	return ret
}

func (s String) Equal(c C) bool {
	var ret bool

	switch c.(type) {
	case String:
		ret = s == c.(String)
	default:
		ret = false
	}
	return ret
}

func (s String) NotEqual(c C) bool {
	var ret bool

	switch c.(type) {
	case String:
		ret = s != c.(String)
	default:
		ret = false
	}
	return ret
}

func (s String) String() string {
	return fmt.Sprintf("\"%s\"", string(s))
}
