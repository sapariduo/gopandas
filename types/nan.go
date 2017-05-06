package types

type Nan string

func NewNan() C {
	return Nan("Nan")
}
func (Nan) Type() Type {
	return NAN
}

func (n Nan) Add(c C) C {
	return n
}

func (n Nan) Sub(c C) C {
	return n
}

func (n Nan) Mul(c C) C {
	return n
}
func (n Nan) Div(c C) C {
	return n
}
func (n Nan) Mod(c C) C {
	return n
}

func (n Nan) Great(c C) bool {
	return false
}

func (n Nan) Less(c C) bool {
	return false
}

func (n Nan) Equal(c C) bool {
	switch c.(type) {
	case Nan:
		return true
	default:
		return false
	}
}

func (n Nan) NotEqual(c C) bool {
	return true
}
