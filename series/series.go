package series

import (
	"fmt"
	"gopandas/types"
	"strings"
	"time"
)

// Index Type
type Index interface{}

// Series Type
type Series map[Index]types.C

// Creates a new serie by passing map or slice
func NewSeries(values interface{}) Series {
	ret := Series{}

	switch values.(type) {
	case map[Index]interface{}:
		for k, v := range values.(map[Index]interface{}) {
			ret[k] = types.NewC(v)
		}
		return ret
	case []interface{}:
		for k, v := range values.([]interface{}) {
			ret[k] = types.NewC(v)
		}
		return ret
	case []int:
		for k, v := range values.([]int) {
			ret[k] = types.Numeric(v)
		}
		return ret
	case []float64:
		for k, v := range values.([]float64) {
			ret[k] = types.Numeric(v)
		}
		return ret
	case []string:
		for k, v := range values.([]string) {
			ret[k] = types.String(v)
		}
		return ret
	case []time.Time:
		for k, v := range values.([]time.Time) {
			ret[k] = types.Time(v)
		}
		return ret
	default:
		fmt.Println("format of series not recognized: use a map or a slice")
		return nil

	}
}

// Returns the summary of each type inside the series
func (s Series) Type() map[types.Type]int {
	ret := map[types.Type]int{}
	for _, v := range s {
		if _, ok := ret[v.Type()]; !ok {
			ret[v.Type()] = 0
		}
		ret[v.Type()]++
	}
	return ret
}

// Returns the length of the series
func (s Series) Len() int {
	return len(s)
}

// Apply a function on a series and returns a new one
func (s Series) Apply(f func(c types.C) types.C) Series {
	ret := Series{}
	for k, v := range s {
		ret[k] = f(v)
	}
	return ret
}

// Returns the number of occurences for each values inside a series
func (s Series) ValuesCount() map[types.C]int {
	ret := map[types.C]int{}

	for _, c := range s {
		if _, ok := ret[c]; !ok {
			ret[c] = 0
		}
		ret[c]++
	}
	return ret
}

func (s Series) String() string {
	ret := "Series:{"
	elements := []string{}
	for k, v := range s {
		elements = append(elements, fmt.Sprintf("%v:%v", k, v))
	}
	ret += strings.Join(elements, ", ")
	ret += "}\n"
	return ret
}

// Compare if two series are equal
func (s1 Series) Equal(s2 Series) bool {
	if s1.Len() != s2.Len() {
		return false
	}
	for k, v1 := range s1 {
		v2, ok := s2[k]
		if !ok || (v1 != v2) {
			return false
		}
	}
	for k, v2 := range s2 {
		v1, ok := s1[k]
		if !ok || (v1 != v2) {
			return false
		}
	}
	return true
}

// Returns a slice of series's indices
func (s Series) Indices() []Index {
	ret := make([]Index, s.Len())
	i := 0
	for k := range s {
		ret[i] = k
		i++
	}
	return ret
}

// Returns a slice of series's values
func (s Series) Values() []types.C {
	ret := make([]types.C, s.Len())
	i := 0
	for _, v := range s {
		ret[i] = v
		i++
	}
	return ret
}
