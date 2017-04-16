package gopandas

import "time"

type Time time.Time

func (Time) Type() Type {
	return TIME
}

func (t Time) Add(c C) C {
	return t
}

func (t Time) Mul(c C) C {
	return t
}
func (t Time) Div(c C) C {
	return t
}
func (t Time) Mod(c C) C {
	return t
}

func (t Time) Great(c C) bool {
	var ret bool

	switch c.(type) {
	case Time:
		ret = time.Time(t).Unix() > time.Time(c.(Time)).Unix()
	default:
		ret = false
	}
	return ret
}

func (t Time) Less(c C) bool {
	var ret bool

	switch c.(type) {
	case Time:
		ret = time.Time(t).Unix() < time.Time(c.(Time)).Unix()
	default:
		ret = false
	}
	return ret
}

func (t Time) Equal(c C) bool {
	var ret bool

	switch c.(type) {
	case Time:
		ret = time.Time(t).Unix() == time.Time(c.(Time)).Unix()
	default:
		ret = false
	}
	return ret
}

func (t Time) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func (t Time) NotEqual(c C) bool {
	var ret bool

	switch c.(type) {
	case Time:
		ret = time.Time(t).Unix() != time.Time(c.(Time)).Unix()
	default:
		ret = false
	}
	return ret
}
