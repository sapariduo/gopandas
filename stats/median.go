package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
)

func _median(a []types.C) types.C {
	il := len(a)
	// if il == 0 {
	// 	return types.NewNan()
	// }

	if il%2 == 0 {
		sum := a[il/2-1].(types.C).Add(a[il/2].(types.C))
		median := sum.Div(types.Numeric(2))
		return median
	} else {
		switch a[il/2].Type() {
		case types.NUMERIC:
			median := a[il/2].(types.C)
			return median
		default:
			median := types.NewNan()
			return median
		}
	}
}

//Median get median value of series
func Median(input *series.Series) (median types.C, err error) {
	// median := types.NewC(nil)
	if input.Len() == 0 {
		return nil, ErrEmptyInput
	}
	input.Sort()
	vals := input.GetValues()
	median = _median(vals)
	return median, nil
}
