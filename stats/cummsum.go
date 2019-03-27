package stats

import (
	"github.com/sapariduo/gopandas/dataframes"
	"github.com/sapariduo/gopandas/series"
)

// CumulativeSum calculates the cumulative sum of the input slice
func CumulativeSum(data *series.Series, datetime *series.Series) (*dataframes.DataFrame, error) {
	df := dataframes.NewEmpty()
	if data.Len() == 0 {
		return df, ErrEmptyInput
	}
	newSeries := series.NewEmpty()
	sortedDate := series.NewEmpty()
	datetime.Sort()

	for i, x := range datetime.Indices {
		if i == 0 {
			newSeries.Indices = append(newSeries.Indices, i)
			newSeries.Series[i] = data.Series[x]
			sortedDate.Indices = append(sortedDate.Indices, i)
			sortedDate.Series[i] = datetime.Series[x]
		} else {
			newSeries.Indices = append(newSeries.Indices, i)
			newSeries.Series[i] = newSeries.Series[i-1].Add(data.Series[x])
			sortedDate.Indices = append(sortedDate.Indices, i)
			sortedDate.Series[i] = datetime.Series[x]
		}

	}
	df.AddSeries("cumsum", newSeries)
	df.AddSeries("datetime", sortedDate)
	return df, nil
}
