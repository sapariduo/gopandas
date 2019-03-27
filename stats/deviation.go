package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
	"math"
)

// StandardDeviation the amount of variation in the dataset
func StandardDeviation(input *series.Series) (sdev types.C, err error) {
	return StandardDeviationPopulation(input)
}

// StandardDeviationPopulation finds the amount of variation from the population
func StandardDeviationPopulation(input *series.Series) (sdev types.C, err error) {

	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}

	// Get the population variance
	vp, _ := PopulationVariance(input)
	vpC := math.Pow(float64(vp.(types.Numeric)), 0.5)
	// Return the population standard deviation
	return types.Numeric(vpC), nil
}

// StandardDeviationSample finds the amount of variation from a sample
func StandardDeviationSample(input *series.Series) (sdev types.C, err error) {

	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}

	// Get the sample variance
	vs, _ := SampleVariance(input)
	vsC := math.Pow(float64(vs.(types.Numeric)), 0.5)

	// Return the sample standard deviation
	return types.Numeric(vsC), nil
}
