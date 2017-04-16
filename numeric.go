package gopandas

type Numeric float64

func (n Numeric) Add(c C) C {
	var ret C

	switch c.(type) {
	case Numeric:
		ret = Numeric(n + c.(Numeric))
	default:
		ret = n
	}
	return ret
}

func (n Numeric) Mul(c C) C {
	var ret C

	switch c.(type) {
	case Numeric:
		ret = Numeric(n * c.(Numeric))
	default:
		ret = n
	}
	return ret
}

func (n Numeric) Div(c C) C {
	var ret C

	switch c.(type) {
	case Numeric:
		if c.(Numeric) == 0.0 {
			ret = newNan()
		} else {
			ret = Numeric(n / c.(Numeric))
		}
	default:
		ret = n
	}
	return ret
}

func (n Numeric) Mod(c C) C {
	return n
}

func (n Numeric) Great(c C) bool {
	var ret bool

	switch c.(type) {
	case Numeric:
		ret = n > c.(Numeric)
	default:
		ret = false
	}
	return ret
}

func (n Numeric) Equal(c C) bool {
	var ret bool

	switch c.(type) {
	case Numeric:
		ret = n == c.(Numeric)
	default:
		ret = false
	}
	return ret
}

func (n Numeric) Less(c C) bool {
	var ret bool

	switch c.(type) {
	case Numeric:
		ret = n < c.(Numeric)
	default:
		ret = false
	}
	return ret
}

func (n Numeric) NotEqual(c C) bool {
	var ret bool

	switch c.(type) {
	case Numeric:
		ret = n != c.(Numeric)
	default:
		ret = false
	}
	return ret
}

func (Numeric) Type() Type {
	return NUMERIC
}
