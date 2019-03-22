package dataframes

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/sapariduo/gopandas/indices"
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
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

type chanSeries struct {
	column string
	data   *series.Series
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
func (grouper *Groups) Max(columns ...string) *DataFrame {
	ret := NewEmpty()
	var wg sync.WaitGroup
	var grp []string
	if len(columns) == 0 {
		grp = grouper.Columns
	} else {
		grp = columns
	}
	srs := make(chan *chanSeries, len(grp))
	wg.Add(len(grp))
	for _, x := range grp {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for k, v := range grouper.Group {
				indices := v
				dfs := grouper.Df.SelectByIndex(indices).Df[x].Max()
				sr.Set(k, dfs)
			}
			chvalue.data = sr

			srs <- chvalue
		}(x)
	}

	wg.Wait()

	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		ret.AddSeries(chans.column, chans.data)
	}

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), CompositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}

//Min create Max value of Each column based on grouper
func (grouper *Groups) Min(columns ...string) *DataFrame {
	ret := NewEmpty()
	var wg sync.WaitGroup
	var grp []string
	if len(columns) == 0 {
		grp = grouper.Columns
	} else {
		grp = columns
	}
	srs := make(chan *chanSeries, len(grp))
	wg.Add(len(grp))
	for _, x := range grp {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for k, v := range grouper.Group {
				indices := v
				dfs := grouper.Df.SelectByIndex(indices).Df[x].Min()
				sr.Set(k, dfs)
			}
			chvalue.data = sr

			srs <- chvalue
		}(x)
	}

	wg.Wait()

	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		ret.AddSeries(chans.column, chans.data)
	}

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), CompositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}

//Count create Count value of Each column based on grouper
func (grouper *Groups) Count(columns ...string) *DataFrame {
	ret := NewEmpty()
	var wg sync.WaitGroup
	var grp []string
	if len(columns) == 0 {
		grp = grouper.Columns
	} else {
		grp = columns
	}
	srs := make(chan *chanSeries, len(grp))
	wg.Add(len(grp))
	for _, x := range grp {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for k, v := range grouper.Group {
				indices := v
				dfs := grouper.Df.SelectByIndex(indices).Df[x].Len()
				sr.Set(k, dfs)
			}
			chvalue.data = sr

			srs <- chvalue
		}(x)
	}

	wg.Wait()

	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		ret.AddSeries(chans.column, chans.data)
	}

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), CompositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}

//Sum create Max value of Each column based on grouper
func (grouper *Groups) Sum(columns ...string) *DataFrame {
	ret := NewEmpty()
	var wg sync.WaitGroup
	var grp []string
	if len(columns) == 0 {
		grp = grouper.Columns
	} else {
		grp = columns
	}
	srs := make(chan *chanSeries, len(grp))
	wg.Add(len(grp))
	for _, x := range grp {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for k, v := range grouper.Group {
				indices := v
				dfs := grouper.Df.SelectByIndex(indices).Df[x].Sum()
				sr.Set(k, dfs)
			}
			chvalue.data = sr

			srs <- chvalue
		}(x)
	}

	wg.Wait()

	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		ret.AddSeries(chans.column, chans.data)
	}

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), CompositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}

//Mean create Mean value of Each column based on grouper
func (grouper *Groups) Mean(columns ...string) *DataFrame {
	ret := NewEmpty()
	var wg sync.WaitGroup
	var grp []string
	if len(columns) == 0 {
		grp = grouper.Columns
	} else {
		grp = columns
	}
	srs := make(chan *chanSeries, len(grp))
	wg.Add(len(grp))
	for _, x := range grp {
		go func(x string) {
			defer wg.Done()
			chvalue := &chanSeries{column: x}
			sr := series.NewEmpty()
			for k, v := range grouper.Group {
				indices := v
				dfs := grouper.Df.SelectByIndex(indices).Df[x].Mean()
				sr.Set(k, dfs)
			}
			chvalue.data = sr

			srs <- chvalue
		}(x)
	}

	wg.Wait()

	for y := 0; y < cap(srs); y++ {
		chans := <-srs
		ret.AddSeries(chans.column, chans.data)
	}

	for i, y := range grouper.Grouper {
		sr := series.NewEmpty()
		for k := range grouper.Group {
			iv := strings.Split(string(k.(types.String)), CompositeChar)
			sr.Set(k, iv[i])
		}
		ret.AddSeries(y, sr)
	}

	return ret
}
