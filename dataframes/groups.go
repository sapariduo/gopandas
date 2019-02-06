package dataframes

import (
	"fmt"
	"gopandas/indices"
	"gopandas/series"
	"gopandas/types"
	"reflect"
	"strings"
)

type Key map[string]types.C

type Keys struct {
	Key   Key
	Index []indices.Indices
}

// Groups Type
//type Group map[string]types.C
type Groups struct {
	Keys    []Keys
	Columns []string //list of columns subject to groping operation
	Grouper []string //list of columns used for grouping parameter
	Group   map[types.C][]indices.Index
	// Group map[interface{}][]indices.Index
	Df *DataFrame
}

//NewGroup Create New Group based on column name
//TODO: Make multi column grouping
func NewGroup(dataframe *DataFrame, columns ...string) *Groups {
	// ret := &Groups{Columns: []string{}, Grouper: columns, Group: make(map[types.C][]indices.Index), Df: dataframe}
	ret := &Groups{Keys: []Keys{}, Columns: []string{}, Grouper: columns, Group: make(map[types.C][]indices.Index), Df: dataframe}

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

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), compositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
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

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), compositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
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

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), compositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
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

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), compositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
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

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), compositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}
