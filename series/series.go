package series

import (
	"fmt"
	"gopandas/indices"
	"gopandas/types"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

// Series Type
//type Series map[indices.Index]types.C
type Series struct {
	Series  map[indices.Index]types.C
	Indices indices.Indices
}

func (s Series) Get(i indices.Index) (v types.C, ok bool) {
	v, ok = s.Series[i]
	return
}

func (s *Series) Set(i indices.Index, v interface{}) error {
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

func (s *Series) Replace(i indices.Index, v interface{}) error {
	if _, ok := s.Get(i); !ok {
		return fmt.Errorf("Index %v doesn't exist, impossible to replace it's value", i)
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
func (s *Series) getIndex(i indices.Index) (int, error) {
	for j, index := range s.Indices {
		if i == index {
			return j, nil
		}
	}
	return -1, fmt.Errorf("Index %v doesn't exist", i)
}

func (s *Series) Del(i indices.Index) error {
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

func (s *Series) ReIndex(idx indices.Indices) {
	if len(idx) != s.Len() {
		fmt.Println("Error: mismatch lengths of series and new indices")
	}
	newIndices := make(indices.Indices, len(idx))
	newValues := map[indices.Index]types.C{}

	for i, index := range idx {
		newIndices[i] = index
		newValues[index] = s.Series[s.Indices[i]]
	}
	s.Series = newValues
	s.Indices = newIndices
}

func (s *Series) ReArrange(idx indices.Indices) {
	if len(idx) != s.Len() {
		fmt.Println("Error: mismatch lengths of series and new indices")
	}
	newIndices := make(indices.Indices, len(idx))
	newValues := map[indices.Index]types.C{}

	for i, index := range idx {
		newIndices[i] = index
		newValues[index] = s.Series[index]
	}
	s.Series = newValues
	s.Indices = newIndices
}

//NewEmpty creates empty Series
func NewEmpty() *Series {
	return &Series{Series: map[indices.Index]types.C{}, Indices: indices.Indices{}}
}

//New Creates a new serie by passing map or slice
func New(values interface{}) *Series {
	ret := &Series{Series: map[indices.Index]types.C{}, Indices: indices.Indices{}}

	switch values.(type) {
	case map[indices.Index]types.C:
		for k, v := range values.(map[indices.Index]types.C) {
			ret.Set(k, v)
		}
	case map[indices.Index]interface{}:
		for k, v := range values.(map[indices.Index]interface{}) {
			ret.Set(k, types.NewC(v))
		}
	case map[indices.Index]int:
		for k, v := range values.(map[indices.Index]int) {
			ret.Set(k, types.Numeric(v))
		}
	case map[indices.Index]float64:
		for k, v := range values.(map[indices.Index]float64) {
			ret.Set(k, types.Numeric(v))
		}
	case map[indices.Index]string:
		for k, v := range values.(map[indices.Index]string) {
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

//Types methods Returns the summary of each type inside the series
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

//Type Return the type of series
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
	ret := &Series{Series: map[indices.Index]types.C{}, Indices: indices.Indices{}}
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

//Equal method Compare if two series are equal
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

//GetIndices Returns a slice of series's indices
func (s *Series) GetIndices() indices.Indices {
	return s.Indices
}

//GetValues Returns a slice of series's values
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
	ret := &Series{Series: map[indices.Index]types.C{}, Indices: indices.Indices{}}

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

func (s *Series) filter(i interface{}, op types.Operator) indices.Indices {
	c := types.NewC(i)
	ret := indices.Indices{}
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

func (s *Series) FilterLT(i interface{}) indices.Indices {
	return s.filter(i, types.LESS)
}

func (s *Series) FilterLTEQ(i interface{}) indices.Indices {
	return s.filter(i, types.LESSEQ)
}

func (s *Series) FilterGT(i interface{}) indices.Indices {
	return s.filter(i, types.GREATER)
}

func (s *Series) FilterGTEQ(i interface{}) indices.Indices {
	return s.filter(i, types.GREATEREQ)
}

func (s *Series) FilterEQ(i interface{}) indices.Indices {
	return s.filter(i, types.EQUAL)
}
func (s *Series) FilterNEQ(i interface{}) indices.Indices {
	return s.filter(i, types.NOTEQUAL)
}

func (s *Series) Swap(i, j int) {
	s.Indices[i], s.Indices[j] = s.Indices[j], s.Indices[i]
}

func (s *Series) Less(i, j int) bool {
	indexI, indexJ := s.Indices[i], s.Indices[j]
	return s.Series[indexI].Less(s.Series[indexJ])
}

//Sort is Sort by values in ascending order
func (s *Series) Sort() {
	sort.Sort(s)
}

//Reverse is Sort by values in descending order
func (s *Series) Reverse() {
	sort.Sort(sort.Reverse(s))
}

//Max Basic implementation of max
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

//Min Basic implementation of min
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

//Mean Returns the mean of values
func (s Series) Mean() types.C {
	sum := s.Sum()
	return sum.Div(types.Numeric(s.Len()))
}

func (s *Series) Select(indices indices.Indices) *Series {
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

//Median return median value of series
//It will return Nan for other types than NUMERIC
func (s Series) Median() types.C {
	// median := types.NewC(nil)
	s.Sort()
	vals := s.GetValues()
	median := median(vals)
	return median
}

func (s Series) Quantile(Q string) types.C {
	Qx := types.NewC(nil)
	s.Sort()
	vals := s.GetValues()
	il := len(vals)
	if il == 0 {
		return types.NewNan()
	}

	// Find the cutoff places depeding on if
	// the input slice length is even or odd
	var c1 int
	var c2 int
	if il%2 == 0 {
		c1 = il / 2
		c2 = il / 2
	} else {
		c1 = (il - 1) / 2
		c2 = c1 + 1
	}

	switch Q {
	case "Q1":
		Qx = median(vals[:c1])
	case "Q2":
		Qx = median(vals)
	case "Q3":
		Qx = median(vals[c2:])
	default:
		Qx = types.NewNan()
	}

	return Qx
}

func (s Series) StdDev() types.C {
	if s.Len() == 0 {
		return types.NewNan()
	}

	pv := variance(&s)
	if pv.Type() != types.NUMERIC {
		return types.NewNan()
	}

	sd := math.Pow(float64(pv.(types.Numeric)), 0.5)

	return types.Numeric(sd)

}

func median(a []types.C) types.C {
	il := len(a)
	if il == 0 {
		return types.NewNan()
	}

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

func variance(input *Series) (variance types.C) {
	if input.Len() == 0 {
		return types.NewNan()
	}

	m := input.Mean()
	if m.Type() != types.NUMERIC {
		return types.NewNan()
	}
	variance = types.Numeric(0)
	for _, v := range input.Series {
		variance = variance.Add(v.Sub(m).Mul(v.Sub(m)))
	}

	return variance.Div(types.Numeric(input.Len()))

}
