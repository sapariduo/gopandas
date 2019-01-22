package series

import (
	"fmt"
	"gopandas/indices"
	"gopandas/types"
	"testing"
)

func TestNew(t *testing.T) {
	s := New(1)
	if s != nil {
		t.Error("Nop")
	}
	s = New([]int{1, 2, 3})
	if s == nil {
		t.Error("Nop")
	}
}

func TestDel(t *testing.T) {
	s := New(map[indices.Index]int{"one": 0, "two": 1, "three": 2})
	err := s.Del(3)
	if err == nil {
		t.Error("Nop")
	}
	err = s.Del("two")
	if err != nil || !s.Equal(New(map[indices.Index]int{"three": 2, "one": 0})) {
		t.Error("Nop")
	}
	s = New(map[indices.Index]int{"one": 0, "two": 1, "three": 2})
	err = s.Del("one")
	if err != nil {
		t.Error("Nop")
	}
	if !s.Equal(New(map[indices.Index]int{"three": 2, "two": 1})) {
		t.Error("Nop")
	}
}

func TestSeriesTypes(t *testing.T) {
	s := New(map[indices.Index]interface{}{
		0:     1,
		1:     "one",
		"two": types.Nan("Nan"),
		3:     2,
	})
	st := s.Types()

	if st[types.NUMERIC] != 2 {
		t.Error("NUMERIC type should be 2 occurences")
	}
	if st[types.STRING] != 1 {
		t.Error("STRING type should be 1 occurence")
	}
	if st[types.NAN] != 1 {
		t.Error("NAN type should be 1 occurence")
	}
}

func TestSeriesType(t *testing.T) {
	s := New([]int{1, 2, 3})
	if s.Type() != types.NUMERIC {
		t.Errorf("Should be NUMERIC vs %v, details:%v", s.Type(), s.Types())
	}
	s = New([]interface{}{1, "one"})
	if s.Type() != types.MULTI {
		t.Errorf("Should be MULTI vs %v, details:%v", s.Type(), s.Types())
	}
}

func TestEqual(t *testing.T) {
	s1 := New([]int{1, 2, 3})
	s2 := New([]int{1, 2})

	if s1.Equal(s2) != s2.Equal(s1) {
		t.Error("Bug")
	}
	if s1.Equal(s2) {
		t.Error("Nop")
	}
}

func TestApply(t *testing.T) {
	s1 := New([]int{1, 2, 3, 4})
	s2 := s1.Apply(func(c types.C) types.C {
		return c.Add(types.Numeric(1))
	})
	if !s2.Equal(New([]int{2, 3, 4, 5})) {
		t.Error("Not equal")
	}
}

func TestSeriesValuesCount(t *testing.T) {
	tests := []struct {
		c     types.C
		value int
	}{
		{c: types.String("un"), value: 1},
		{c: types.Numeric(1), value: 2},
		{c: types.Numeric(2), value: 1},
		{c: types.NewNan(), value: 1},
	}
	s := New(map[indices.Index]interface{}{
		0:      1,
		5:      1,
		1:      "un",
		"deux": types.NewNan(),
		3:      2,
	})
	counts := s.ValuesCount()
	for _, test := range tests {
		if counts[test.c] != test.value {
			t.Errorf("Error: %v:%d vs %v:%d", test.c, counts[test.c], test.c, test.value)

		}
	}
}

func TestAddSub(t *testing.T) {
	s1 := New([]int{1, 2, 3})
	s2 := New([]int{-1, -2, -3})
	s3 := New([]int{0, 0, 0})

	if !s1.Add(s2).Equal(s3) {
		t.Error("Error Add")
	}
	if s := New([]string{"1", "2", "3"}).Add(New(map[indices.Index]int{1: 1, 2: 2, 3: 3})); s != nil {
		t.Error("Error Add", s)
	}
	if !s1.Sub(s3).Equal(s1) {
		t.Error("Error Sub")
	}

}
func TestMulDivMod(t *testing.T) {
	s1 := New([]int{1, 1, 1})
	s2 := New([]int{0, 0, 0})
	s3 := New([]int{1, 2, 3})

	if !s1.Add(s1).Div(s1).Equal(New([]int{2, 2, 2})) {
		t.Error("Error Div")
	}
	if !s1.Mul(s2).Equal(s2.Mul(s1)) {
		t.Error("Error mul")
	}
	if !s3.Mul(s3).Div(s3).Equal(s3) {
		t.Error("Error mul, div")
	}
}

func TestFilter(t *testing.T) {
	s := New([]int{6, 7, 8})

	for _, test := range []struct {
		idxtest, idxtrue indices.Indices
	}{
		{idxtest: s.FilterEQ(7), idxtrue: indices.Indices{1}},
		{idxtest: s.FilterNEQ(7), idxtrue: indices.Indices{0, 2}},
		{idxtest: s.FilterGT(7), idxtrue: indices.Indices{2}},
		{idxtest: s.FilterGTEQ(7), idxtrue: indices.Indices{1, 2}},
		{idxtest: s.FilterLT(7), idxtrue: indices.Indices{0}},
		{idxtest: s.FilterLTEQ(7), idxtrue: indices.Indices{0, 1}},
	} {
		if !test.idxtest.Equal(test.idxtrue) {
			t.Error(test.idxtest, "vs", test.idxtrue)
		}
	}
}

func TestSelect(t *testing.T) {
	s := New([]int{6, 7, 8})

	if !s.Select(indices.Indices{0, 2}).Equal(New(map[indices.Index]int{0: 6, 2: 8})) {
		t.Error("Nop")
	}
}

/* Test sort not valid
func TestSort(t *testing.T) {
	s := New(map[Index]int{0: 3, 1: 1, 2: 2})
	s.Sort()
	fmt.Print(s)
	if !s.Equal(New(map[Index]int{1: 1, 2: 2, 0: 3})) {
		t.Errorf("Nop, %v", s)
	}
	s.Reverse()
	fmt.Print(s)
	if !s.Equal(New(map[Index]int{0: 3, 2: 2, 1: 1})) {
		t.Errorf("Nop, %v", s)
	}

}*/

func TestMinMax(t *testing.T) {
	s := New([]float64{1.1, 2, 3, 4, -1, -2})

	if s.Max() != types.Numeric(4) {
		t.Error("Error max")
	}
	if s.Min() != types.Numeric(-2) {
		t.Error("Error min")
	}
}

func TestSumMean(t *testing.T) {
	s := New([]int{1, 2, 3})

	if s.Sum().NotEqual(types.Numeric(6)) {
		t.Error("Error Sum")
	}

	if s.Mean().NotEqual(types.Numeric(2)) {
		t.Error("Error Mean")
	}
}

func TestSeries_Median(t *testing.T) {
	s := New([]float64{1.1, 2, 3, 4, -1})
	fmt.Println(s.Median())

	s.Sort()
	fmt.Println(s)
}
