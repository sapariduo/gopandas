package stats

import (
	"github.com/sapariduo/gopandas/dataframes"
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
)

// VarP is a shortcut to PopulationVariance
func VarP(input *series.Series) (sdev types.C, err error) {
	return PopulationVariance(input)
}

// VarS is a shortcut to SampleVariance
func VarS(input *series.Series) (sdev types.C, err error) {
	return SampleVariance(input)
}

// StdDevP is a shortcut to StandardDeviationPopulation
func StdDevP(input *series.Series) (sdev types.C, err error) {
	return StandardDeviationPopulation(input)
}

// StdDevS is a shortcut to StandardDeviationSample
func StdDevS(input *series.Series) (sdev types.C, err error) {
	return StandardDeviationSample(input)
}

// func CummSum(data *series.Series, datetime *series.Series) (*dataframes.DataFrame, error) {
// 	return CumulativeSum(data, datetime)
// }

// CumSum is a shortcut to Cummulative Sum
func CumSum(axis, datetime string, dataVal []string, dataframe *dataframes.DataFrame) (*dataframes.DataFrame, error) {
	return Cummulative(axis, datetime, dataVal, dataframe)
}
