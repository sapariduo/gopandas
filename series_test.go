package gopandas

import "testing"

func TestSeriesType(t *testing.T) {
	s := Series{0: Numeric(1), 1: String("un"), "deux": Nan("Nan"), 3: Numeric(2)}
	st := s.Type()

	if st[NUMERIC] != 2 {
		t.Error("NUMERIC type should be 2 occurences")
	}
	if st[STRING] != 1 {
		t.Error("STRING type should be 1 occurence")
	}
	if st[NAN] != 1 {
		t.Error("NAN type should be 1 occurence")
	}
}
