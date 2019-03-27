package stats

import (
	"fmt"
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
)

// Percentile finds the relative standing in a slice of floats
func Percentile(input *series.Series, percent float64) (percentile types.C, err error) {

	if input.Len() == 0 {
		return types.NewNan(), ErrEmptyInput
	}

	if percent <= 0 || percent > 100 {
		return types.NewNan(), ErrBounds
	}

	// Start by sorting a copy of the slice
	input.Sort()

	c := input.GetValues()
	// Multiply percent by length of input

	index := (percent / 100) * float64(len(c))

	// Check if the index is a whole number
	if index == float64(int64(index)) {

		// Convert float to int
		i := int(index)

		// Find the value at the index
		percentile = c[i-1].Add(c[i]).Div(types.Numeric(2))

	} else if index > 1 {
		fmt.Println(index)
		// Convert float to int via truncation
		i := int(index)
		// Find the average of the index and following values
		percentile = c[i]

	} else {
		return types.NewNan(), ErrBounds
	}

	return percentile, nil

}
