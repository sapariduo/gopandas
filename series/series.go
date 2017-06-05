package series

import (
	"fmt"
	"gopandas/types"
	"log"
	"sort"
	"strings"
	"time"
)

// Index Type
type Index interface{}

type Indices []Index

// Series Type
//type Series map[Index]types.C

type Series struct {
	Series  map[Index]types.C
	Indices Indices
}

func (idx1 Indices) equal(idx2 Indices) bool {
	if len(idx1) != len(idx2) {
		return false
	}
	for i := range idx1 {
		if idx1[i] != idx2[i] {
			return false
		}
	}
	return true
}

func (s Series) Get(i Index) (v types.C, ok bool) {
	v, ok = s.Series[i]
	return
}

func (s *Series) Set(i Index, v interface{}) error {
	if _, ok := s.Get(i); !ok {
		s.Indices = append(s.Indices, i)
		switch v.(type) {
		case types.C:
			s.Series[i] = v.(types.C)
		default:
			s.Series[i] = types.NewC(v)
		}
		return nil
	} else {
		return fmt.Errorf("Index already set, use Replace method of you want to replace the value")
	}
}

func (s *Series) Replace(i Index, v interface{}) error {
	if _, ok := s.Get(i); !ok {
		return fmt.Errorf("Index %v doesn't exist, impossible to replace it's value")
	}
	switch v.(type) {
	case types.C:
		s.Series[i] = v.(types.C)
	default:
		s.Series[i] = types.NewC(v)
	}
	return nil
}

// Returns the position in the slice of the Index i
func (s *Series) getIndex(i Index) (int, error) {
	for j, index := range s.Indices {
		if i == index {
			return j, nil
		}
	}
	return -1, fmt.Errorf("Index %v doesn't exist", i)
}

func (s *Series) Del(i Index) error {
	index_i, err := s.getIndex(i)
	if err != nil {
		return err
	}
	if index_i == s.Len()-1 {
		s.Indices = s.Indices[:index_i]
	} else {
		s.Indices = append(s.Indices[:index_i], s.Indices[index_i+1:]...)
	}
	delete(s.Series, i)
	return nil
}

func (s *Series) ReIndex(indices Indices) {
	if len(indices) != s.Len() {
		fmt.Println("Error: mismatch lengths of series and new indices")
	}
	newIndices := make(Indices, len(indices))
	newValues := map[Index]types.C{}

	for i, index := range indices {
		newIndices[i] = index
		newValues[index] = s.Series[s.Indices[i]]
	}
	s.Series = newValues
	s.Indices = newIndices
}

func NewEmpty() *Series {
	return &Series{Series: map[Index]types.C{}, Indices: Indices{}}
}

// Creates a new serie by passing map or slice
func New(values interface{}) *Series {
	ret := &Series{Series: map[Index]types.C{}, Indices: Indices{}}

	switch values.(type) {
	case map[Index]types.C:
		for k, v := range values.(map[Index]types.C) {
			ret.Set(k, v)
		}
	case map[Index]interface{}:
		for k, v := range values.(map[Index]interface{}) {
			ret.Set(k, types.NewC(v))
		}
	case map[Index]int:
		for k, v := range values.(map[Index]int) {
			ret.Set(k, types.Numeric(v))
		}
	case map[Index]float64:
		for k, v := range values.(map[Index]float64) {
			ret.Set(k, types.Numeric(v))
		}
	case map[Index]string:
		for k, v := range values.(map[Index]string) {
			ret.Set(k, types.String(v))
		}
	case []interface{}:
		for k, v := range values.([]interface{}) {
			ret.Set(k, types.NewC(v))
		}
	case []int:
		for k, v := range values.([]int) {
			ret.Set(k, types.Numeric(v))
		}
	case []float64:
		for k, v := range values.([]float64) {
			ret.Set(k, types.Numeric(v))
		}
	case []string:
		for k, v := range values.([]string) {
			ret.Set(k, types.String(v))
		}
	case []time.Time:
		for k, v := range values.([]time.Time) {
			ret.Set(k, types.Time(v))
		}
	default:
		fmt.Println("format of series not recognized: use a map or a slice")
		return nil
	}
	return ret
}

// Returns the summary of each type inside the series
func (s *Series) Types() map[types.Type]int {
	ret := map[types.Type]int{}
	for _, v := range s.Series {
		if _, ok := ret[v.Type()]; !ok {
			ret[v.Type()] = 0
		}
		ret[v.Type()]++
	}
	return ret
}

// Return the type of series
func (s *Series) Type() types.Type {
	res := s.Types()
	var t types.Type

	if len(res) > 1 {
		t = types.MULTI
	} else {
		for k := range res {
			t = k
		}
	}
	return t
}

// Returns the length of the series
func (s *Series) Len() int {
	return len(s.Series)
}

// Apply a function on a series and returns a new one
func (s *Series) Apply(f func(c types.C) types.C) *Series {
	ret := &Series{Series: map[Index]types.C{}, Indices: Indices{}}
	for k, v := range s.Series {
		ret.Set(k, f(v))
	}
	return ret
}

// Returns the number of occurences for each values inside a series
func (s Series) ValuesCount() map[types.C]int {
	ret := map[types.C]int{}

	for _, c := range s.Series {
		if _, ok := ret[c]; !ok {
			ret[c] = 0
		}
		ret[c]++
	}
	return ret
}

func (s *Series) String() string {
	ret := "Series:{"
	elements := []string{}
	for _, index := range s.Indices {
		if v, ok := s.Get(index); !ok {
			log.Panic(fmt.Sprintf("critical error: values:%+v, indices:%+v", s.Series, s.Indices))
		} else {
			elements = append(elements, fmt.Sprintf("%v:%v", index, v))
		}
	}
	ret += strings.Join(elements, ", ")
	ret += "}\n"
	return ret
}

// Compare if two series are equal
func (s1 *Series) Equal(s2 *Series) bool {
	if s1.Len() != s2.Len() {
		return false
	}
	for k, v1 := range s1.Series {
		v2, ok := s2.Get(k)
		if !ok || (v1 != v2) {
			return false
		}
	}
	for k, v2 := range s2.Series {
		v1, ok := s1.Get(k)
		if !ok || (v1 != v2) {
			return false
		}
	}
	return true
}

// Returns a slice of series's indices
func (s *Series) GetIndices() Indices {
	return s.Indices
}

// Returns a slice of series's values
func (s *Series) GetValues() []types.C {
	ret := make([]types.C, s.Len())
	for i, index := range s.Indices {
		if v, ok := s.Get(index); !ok {
			panic(fmt.Errorf("Critical error"))
		} else {
			ret[i] = v
		}
	}
	return ret
}

func (s1 *Series) op(s2 *Series, op types.Operator) *Series {
	if s1.Len() != s2.Len() {
		return nil
	}
	for k := range s1.Indices {
		if _, ok := s2.Get(k); !ok {
			return nil
		}
	}
	for k := range s2.Indices {
		if _, ok := s1.Get(k); !ok {
			return nil
		}
	}
	ret := &Series{Series: map[Index]types.C{}, Indices: Indices{}}

	for _, index := range s1.Indices {
		v1, _ := s1.Get(index)
		v2, _ := s2.Get(index)
		switch op {
		case types.ADD:
			ret.Set(index, v1.Add(v2))
		case types.MUL:
			ret.Set(index, v1.Mul(v2))
		case types.DIV:
			ret.Set(index, v1.Div(v2))
		case types.MOD:
			ret.Set(index, v1.Mod(v2))
		case types.SUB:
			ret.Set(index, v1.Sub(v2))
		default:
			return nil
		}
	}
	return ret
}

func (s1 *Series) Add(s2 *Series) *Series {
	return s1.op(s2, types.ADD)
}
func (s1 *Series) Sub(s2 *Series) *Series {
	return s1.op(s2, types.SUB)
}
func (s1 *Series) Mul(s2 *Series) *Series {
	return s1.op(s2, types.MUL)
}
func (s1 *Series) Div(s2 *Series) *Series {
	return s1.op(s2, types.DIV)
}
func (s1 *Series) Mod(s2 *Series) *Series {
	return s1.op(s2, types.MOD)
}

func (s *Series) filter(i interface{}, op types.Operator) Indices {
	c := types.NewC(i)
	ret := Indices{}
	for _, index := range s.Indices {
		value := s.Series[index]
		switch op {
		case types.LESS:
			if value.Less(c) {
				ret = append(ret, index)
			}
		case types.LESSEQ:
			if value.Less(c) || value.Equal(c) {
				ret = append(ret, index)
			}
		case types.GREATER:
			if value.Great(c) {
				ret = append(ret, index)
			}
		case types.GREATEREQ:
			if value.Great(c) || value.Equal(c) {
				ret = append(ret, index)
			}
		case types.EQUAL:
			if value.Equal(c) {
				ret = append(ret, index)
			}
		case types.NOTEQUAL:
			if !value.Equal(c) {
				ret = append(ret, index)
			}
		}
	}
	return ret
}

func (s *Series) FilterLT(i interface{}) Indices {
	return s.filter(i, types.LESS)
}

func (s *Series) FilterLTEQ(i interface{}) Indices {
	return s.filter(i, types.LESSEQ)
}

func (s *Series) FilterGT(i interface{}) Indices {
	return s.filter(i, types.GREATER)
}

func (s *Series) FilterGTEQ(i interface{}) Indices {
	return s.filter(i, types.GREATEREQ)
}

func (s *Series) FilterEQ(i interface{}) Indices {
	return s.filter(i, types.EQUAL)
}
func (s *Series) FilterNEQ(i interface{}) Indices {
	return s.filter(i, types.NOTEQUAL)
}

func (s *Series) Swap(i, j int) {
	s.Indices[i], s.Indices[j] = s.Indices[j], s.Indices[i]
}

func (s *Series) Less(i, j int) bool {
	index_i, index_j := s.Indices[i], s.Indices[j]
	return s.Series[index_i].Less(s.Series[index_j])
}

//Sort by values in ascending order
func (s *Series) Sort() {
	sort.Sort(s)
}

//Sort by values in descending order
func (s *Series) Reverse() {
	sort.Sort(sort.Reverse(s))
}

// Basic implementation of max
func (s Series) Max() types.C {
	i := true
	var max types.C
	for _, v := range s.Series {
		if i {
			max = v
			i = false
		} else {
			if max.Less(v) {
				max = v
			}
		}
	}
	return max
}

// Basic implementation of min
func (s Series) Min() types.C {
	i := true
	var min types.C
	for _, v := range s.Series {
		if i {
			min = v
			i = false
		} else {
			if min.Great(v) {
				min = v
			}
		}
	}
	return min
}

// Returns the sum of values
func (s Series) Sum() types.C {
	i := true
	var sum types.C
	for _, v := range s.Series {
		if i {
			sum = v
			i = false
		} else {
			sum = sum.Add(v)
		}
	}
	return sum
}

// Returns the mean of values
func (s Series) Mean() types.C {
	sum := s.Sum()
	return sum.Div(types.Numeric(s.Len()))
}

func (s *Series) Select(indices Indices) *Series {
	values := []interface{}{}
	for _, index := range indices {
		if _, ok := s.Series[index]; ok {
			values = append(values, s.Series[index])
		}
	}
	ret := New(values)
	ret.ReIndex(indices)
	return ret
}
