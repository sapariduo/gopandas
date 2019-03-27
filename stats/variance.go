package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
)

func _variance(input *series.Series, sample int) (variance types.C, err error) {
	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}

	m := input.Mean()
	if m.Type() != types.NUMERIC {
		return types.NewNan(), ErrNaN
	}
	vals := input.GetValues()
	variance = types.Numeric(0)
	for _, v := range vals {
		variance = variance.Add(v.Sub(m).Mul(v.Sub(m)))
	}
	division := float64(types.Numeric(input.Len())) - float64(1*sample)

	return variance.Div(types.Numeric(division)), nil
}

// PopulationVariance finds the amount of variance within a population
func PopulationVariance(input *series.Series) (pvar types.C, err error) {

	v, err := _variance(input, 0)
	if err != nil {
		return types.NewNan(), err
	}

	return v, nil
}

// SampleVariance finds the amount of variance within a sample
func SampleVariance(input *series.Series) (svar types.C, err error) {

	v, err := _variance(input, 1)
	if err != nil {
		return types.NewNan(), err
	}

	return v, nil
}
