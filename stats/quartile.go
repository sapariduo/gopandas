package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
)

//Quartiles collection
type Quartiles struct {
	Q1 types.C
	Q2 types.C
	Q3 types.C
}

//Quartile provide value of Q1, Q2, Q3 of series
func Quartile(input *series.Series) (Quartiles, error) {
	input.Sort()
	vals := input.GetValues()
	il := len(vals)
	if il == 0 {
		return Quartiles{}, ErrEmptyInput
	}

	// Find the cutoff places depeding on if
	// the input slice length is even or odd
	var c1 int
	var c2 int
	if il%2 == 0 {
		c1 = il / 2
		c2 = il / 2
	} else {
		c1 = (il - 1) / 2
		c2 = c1 + 1
	}

	Q1 := _median(vals[:c1])
	Q2 := _median(vals)
	Q3 := _median(vals[c2:])

	return Quartiles{Q1, Q2, Q3}, nil
}

// InterQuartileRange finds the range between Q1 and Q3
func InterQuartileRange(input *series.Series) (types.C, error) {
	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}
	qs, _ := Quartile(input)
	iqr := qs.Q3.Sub(qs.Q1)
	return iqr, nil
}

// Midhinge finds the average of the first and third quartiles
func Midhinge(input *series.Series) (types.C, error) {
	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}
	qs, _ := Quartile(input)
	mh := (qs.Q1.Add(qs.Q3)).Div(types.Numeric(2))
	return mh, nil
}

// Trimean finds the average of the median and the midhinge
func Trimean(input *series.Series) (types.C, error) {
	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}

	// c := sortedCopy(input)
	q, _ := Quartile(input)

	return (q.Q1.Add(q.Q2.Mul(types.Numeric(2))).Add(q.Q3)).Div(types.Numeric(4)), nil
}
