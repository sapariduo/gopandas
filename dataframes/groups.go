package dataframes

import (
	"fmt"
	"gopandas/indices"
	"gopandas/series"
	"gopandas/types"
	"reflect"
)

// Groups Type
//type Group map[string]types.C
type Groups struct {
	Columns []string
	Group   map[types.C][]indices.Index
	Df      *DataFrame
}

//NewGroup Create New Group based on column name
//TODO: Make multi column grouping
func NewGroup(column string, dataframe *DataFrame) *Groups {
	ret := &Groups{Columns: []string{}, Group: make(map[types.C][]indices.Index), Df: dataframe}
	return ret
}

//Info create Info value of grouper index
func (grouper *Groups) Info() *DataFrame {
	ret := NewEmpty()

	sr := series.NewEmpty()
	for k, v := range grouper.Group {
		indices := reflect.TypeOf(v).String()
		vals := reflect.ValueOf(v)
		indicesVals := vals
		idx := fmt.Sprintf("%v%v", indices, indicesVals)
		sr.Set(k, idx)
	}
	ret.AddSeries("cols", sr)

	return ret
}

//Max create Max value of Each column based on grouper
func (grouper *Groups) Max() *DataFrame {
	ret := NewEmpty()

	for _, x := range grouper.Columns {
		sr := series.NewEmpty()
		for k, v := range grouper.Group {
			indices := v
			dfs := grouper.Df.SelectByIndex(indices).Df[x].Max()
			sr.Set(k, dfs)
		}
		ret.AddSeries(x, sr)
	}

	return ret
}

//Min create Max value of Each column based on grouper
func (grouper *Groups) Min() *DataFrame {
	ret := NewEmpty()

	for _, x := range grouper.Columns {
		sr := series.NewEmpty()
		for k, v := range grouper.Group {
			indices := v
			dfs := grouper.Df.SelectByIndex(indices).Df[x].Min()
			sr.Set(k, dfs)
		}
		ret.AddSeries(x, sr)
	}

	return ret
}

//Count create Count value of Each column based on grouper
func (grouper *Groups) Count() *DataFrame {
	ret := NewEmpty()

	for _, x := range grouper.Columns {
		sr := series.NewEmpty()
		for k, v := range grouper.Group {
			indices := v
			dfs := grouper.Df.SelectByIndex(indices).Len()
			sr.Set(k, dfs)
		}
		ret.AddSeries(x, sr)
	}

	return ret
}

//Sum create Max value of Each column based on grouper
func (grouper *Groups) Sum() *DataFrame {
	ret := NewEmpty()

	for _, x := range grouper.Columns {
		sr := series.NewEmpty()
		for k, v := range grouper.Group {
			indices := v
			dfs := grouper.Df.SelectByIndex(indices).Df[x].Sum()
			sr.Set(k, dfs)
		}
		ret.AddSeries(x, sr)
	}

	return ret
}

//Mean create Mean value of Each column based on grouper
func (grouper *Groups) Mean() *DataFrame {
	ret := NewEmpty()

	for _, x := range grouper.Columns {
		sr := series.NewEmpty()
		for k, v := range grouper.Group {
			indices := v
			dfs := grouper.Df.SelectByIndex(indices).Df[x].Mean()
			sr.Set(k, dfs)
		}
		ret.AddSeries(x, sr)
	}

	return ret
}
