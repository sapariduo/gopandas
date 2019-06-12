package stats

import (
	"sync"

	"github.com/sapariduo/gopandas/dataframes"
	"github.com/sapariduo/gopandas/series"
)

type chanSeries struct {
	column string
	data   *series.Series
}

// Cummulative calculates the cummulative sum of the input slice
func Cummulative(axis, datetime string, dataVal []string, dataframe *dataframes.DataFrame) (*dataframes.DataFrame, error) {
	dfsource := dataframe.Df
	df := dataframes.NewEmpty()
	if dfsource[axis].Len() == 0 {
		return df, ErrEmptyInput
	}
	var wg sync.WaitGroup
	srs := make(chan *chanSeries, len(dataVal))
	wg.Add(len(dataVal))

	newAxis := series.NewEmpty()
	sortedDate := series.NewEmpty()
	dfsource[datetime].Sort()

	for _, x := range dataVal {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for ix, yx := range dfsource[datetime].Indices {
				if ix == 0 {
					sr.Indices = append(sr.Indices, ix)
					sr.Series[ix] = dfsource[x].Series[yx]

				} else {
					sr.Indices = append(sr.Indices, ix)
					sr.Series[ix] = sr.Series[ix-1].Add(dfsource[x].Series[yx])

				}
			}
			chvalue.data = sr
			srs <- chvalue

		}(x)
	}

	for ai, ax := range dfsource[datetime].Indices {
		newAxis.Indices = append(newAxis.Indices, ai)
		newAxis.Series[ai] = dfsource[axis].Series[ax]
		sortedDate.Indices = append(sortedDate.Indices, ai)
		sortedDate.Series[ai] = dfsource[datetime].Series[ax]

	}

	df.AddSeries(datetime, sortedDate)
	df.AddSeries(axis, newAxis)

	wg.Wait()
	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		df.AddSeries(chans.column, chans.data)
	}

	return df, nil
}
