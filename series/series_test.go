package series

import (
	"gopandas/types"
	"testing"
)

func TestNewSeries(t *testing.T) {
	s := NewSeries(1)
	if s != nil {
		t.Error("Nop")
	}
	s = NewSeries([]int{1, 2, 3})
	if s == nil {
		t.Error("Nop")
	}
}

func TestSeriesType(t *testing.T) {
	s := Series{
		0:      types.Numeric(1),
		1:      types.String("un"),
		"deux": types.Nan("Nan"),
		3:      types.Numeric(2),
	}
	st := s.Type()

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

func TestEqual(t *testing.T) {
	s1 := NewSeries([]int{1, 2, 3})
	s2 := NewSeries([]int{1, 2})

	if s1.Equal(s2) != s2.Equal(s1) {
		t.Error("Bug")
	}
	if s1.Equal(s2) {
		t.Error("Nop")
	}
}

func TestApply(t *testing.T) {
	s1 := NewSeries([]int{1, 2, 3, 4})
	s2 := s1.Apply(func(c types.C) types.C {
		return c.Add(types.Numeric(1))
	})
	if !s2.Equal(NewSeries([]int{2, 3, 4, 5})) {
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
		{c: types.Nan("Nan"), value: 1},
	}
	s := Series{
		0:      types.Numeric(1),
		5:      types.Numeric(1),
		1:      types.String("un"),
		"deux": types.Nan("Nan"),
		3:      types.Numeric(2),
	}
	counts := s.ValuesCount()
	for _, test := range tests {
		if counts[test.c] != test.value {
			t.Errorf("Error: %d vs %d", counts[test.c], test.value)

		}
	}
}
