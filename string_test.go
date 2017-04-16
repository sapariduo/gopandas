package gopandas

import "testing"

func TestStringType(t *testing.T) {
	var s String

	switch s.Type() {
	case STRING:
		return
	default:
		t.Error("Wrong type")
	}
}

func TestStringAdd(t *testing.T) {
	s1 := String("foo")
	s2 := String("bar")
	s3 := String("baz")
	if s1.Add(s2.Add(s3)) != s1.Add(s2).Add(s3) {
		t.Error("Error String Add String")
	}
}

func TestStringMul(t *testing.T) {
	s1 := String("foo")
	s2 := String("bar")
	if s1.Mul(s2) != s1 {
		t.Error("Error String Mul String")
	}
}

func TestStringDiv(t *testing.T) {
	s1 := String("foo")
	s2 := String("bar")
	if s1.Div(s2) != s1 {
		t.Error("Error String Div String")
	}
}

func TestStringMod(t *testing.T) {
	s1 := String("foo")
	s2 := String("bar")
	if s1.Mod(s2) != s1 {
		t.Error("Error String Mod String")
	}
}

func TestStringLess(t *testing.T) {
	if String("a").Less(String("b")) == false {
		t.Error("Error String Less String")
	}
	if String("aa").Less(String("b")) == false {
		t.Error("Error String Less String")
	}
	if String("").Less(String("b")) == false {
		t.Error("Error String Less String")
	}
	if String("").Less(String("")) == true {
		t.Error("Error String Less String")
	}
}

func TestStringGreat(t *testing.T) {
	if String("a").Great(String("b")) == true {
		t.Error("Error String Great String")
	}
	if String("aa").Great(String("b")) == true {
		t.Error("Error String Great String")
	}
	if String("").Great(String("b")) == true {
		t.Error("Error String Great String")
	}
	if String("").Great(String("")) == true {
		t.Error("Error String Great String")
	}
}

func TestStringEqual(t *testing.T) {
	if String("").Equal(String("")) == false {
		t.Error("Error String Equal String")
	}
}

func TestStringNotEqual(t *testing.T) {
	if String("").NotEqual(String("")) == true {
		t.Error("Error String Equal String")
	}
}
