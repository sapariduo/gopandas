package dataframes

import (
	//"bufio"

	"bytes"
	"encoding/csv"
	"fmt"
	"gopandas/series"
	"gopandas/types"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

// DataFrame is the structure of a dataframe, Data are avaibable with Df attribute
// Final goal is to have all attributes in private.
type DataFrame struct {
	Columns []string
	Indices []series.Index
	Df      map[string]series.Series
	NbLines int
}

// ConfigDataFrame is a structure of configuration to create a dataframe
type ConfigDataFrame struct {
	File       string
	Header     bool
	Index      bool
	Sep        rune
	TimeLayout string
}

// Basic Parser of string in a interface{}
// Can be a float or a time in RFC3339 format
// To be completed...
func convertTo(s string) interface{} {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err == nil {
		return t
	}
	return s
}

// NewDataFrameCSV is a function to create of a dataframe from a CSV file
func NewDataFrameCSV(c *ConfigDataFrame) *DataFrame {
	fd, err := os.Open(c.File)
	if err != nil {
		return nil
	}
	defer fd.Close()

	df := DataFrame{}
	df.Df = map[string]series.Series{}
	reader := csv.NewReader(fd)
	reader.Comma = c.Sep
	firstline, err := reader.Read()
	checkerr(err)
	lines, err := reader.ReadAll()
	checkerr(err)

	index := 0
	if c.Header {
		if c.Index {
			df.Columns = firstline[1:]
		} else {
			df.Columns = firstline
		}
		for _, col := range df.Columns {
			df.Df[col] = series.Series{}
		}
	} else {
		if c.Index {
			df.Columns = make([]string, len(firstline[1:]))
			df.Indices = append(df.Indices, firstline[0])
			for i, c := range firstline[1:] {
				col := fmt.Sprintf("%v", i)
				df.Columns[i] = col
				df.Df[col] = series.Series{}
				df.Df[col][firstline[0]] = types.NewC(convertTo(c))
			}
		} else {
			df.Columns = make([]string, len(firstline))
			df.Indices = append(df.Indices, index)
			for i, c := range firstline {
				col := fmt.Sprintf("%v", i)
				df.Columns[i] = col
				df.Df[col] = series.Series{}
				df.Df[col][index] = types.NewC(convertTo(c))
			}
		}
		df.NbLines++
		index++
	}
	for _, line := range lines {
		if c.Index {
			df.Indices = append(df.Indices, line[0])
		} else {
			df.Indices = append(df.Indices, index)
		}
		for icol, col := range df.Columns {
			if c.Index {
				v := types.NewC(convertTo(line[icol+1]))
				df.Df[col][line[0]] = v
			} else {
				v := types.NewC(convertTo(line[icol]))
				df.Df[col][index] = v
			}
		}
		index++
		df.NbLines++
	}
	return &df
}

func (df *DataFrame) String() string {
	b := &bytes.Buffer{}
	table := tablewriter.NewWriter(b)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetAutoFormatHeaders(false)
	header := []string{"Index"}
	header = append(header, df.Columns...)
	table.SetHeader(header)
	for _, index := range df.Indices {
		l := []string{fmt.Sprintf("%v", index)}
		for _, col := range df.Columns {
			l = append(l, fmt.Sprintf("%v", df.Df[col][index]))
		}
		table.Append(l)
	}
	footer := make([]string, len(header))
	footer[0] = fmt.Sprintf("COUNT:%v", df.NbLines)
	for i, col := range df.Columns {
		m := df.Df[col].Type()
		if len(m) == 0 {
			log.Panicf("Error: Types are not defined for [%v] column\n", col)
		} else if len(m) > 1 {
			footer[i+1] = "MULTI"
		} else {
			for k := range m {
				footer[i+1] = string(k)
			}
		}
	}
	table.SetFooter(footer)
	table.Render()
	raw, err := ioutil.ReadAll(b)
	checkerr(err)
	return string(raw)
}

// Select function is used to select colums and return a dataframe with selected columns
// Warning Select function doesn't make a full copy
// To change Data you shoud use the Copy function before or after
func (df *DataFrame) Select(cols ...string) *DataFrame {
	dfs := &DataFrame{}
	dfs.NbLines = df.NbLines
	dfs.Df = map[string]series.Series{}
	dfs.Indices = df.Indices
	for _, col := range cols {
		if _, ok := df.Df[col]; !ok {
			fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		} else {
			dfs.Columns = append(dfs.Columns, col)
			dfs.Df[col] = df.Df[col]
		}
	}
	if len(dfs.Columns) != 0 {
		return dfs
	}
	return nil

}

/*
// Copy function makes a full copy of a dataframe
func (df *DataFrame) Copy() *DataFrame {
	dfs := &DataFrame{}
	dfs.NbLines = df.NbLines
	dfs.Df = map[string][]C{}
	dfs.Types = map[string]map[Type]int{}
	for _, col := range df.Columns {
		_, ok := df.Df[col]
		if ok == false {
			fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		} else {
			dfs.Columns = append(dfs.Columns, col)
			dfs.Df[col] = make([]C, dfs.NbLines)
			dfs.Types[col] = map[Type]int{}
			for k, v := range df.Types[col] {
				dfs.Types[col][k] = v
			}
			for i, v := range df.Df[col] {
				dfs.Df[col][i] = v
			}
		}
	}
	return dfs
}

// SetList function is used to add a new column of data.
// l argument should be a slice of something
// Example: df.SetList("foo",[]int{1,2,3,...})
func (df *DataFrame) SetList(l interface{}, col string) {
	v := reflect.ValueOf(l)
	t := v.Kind()
	if t != reflect.Slice {
		fmt.Printf("Error: Need to Use a Slice in 2nd arg instead of [%v]\n", t)
		return
	}
	length := v.Len()
	if length != df.NbLines && df.NbLines != 0 {
		fmt.Printf("Error: Tried to insert a list with len [%v] in DataFrame with len [%v]\n", length, df.NbLines)
		return
	}
	if df.Df == nil {
		df.Df = map[string][]C{}
	}
	if df.Types == nil {
		df.Types = map[string]map[Type]int{}
	}
	_, ok := df.Df[col]
	if ok == false {
		df.Columns = append(df.Columns, col)
		df.Types[col] = map[Type]int{}
	} else {
		fmt.Printf("Warning: Column [%v] has been replaced\n", col)
	}
	tmp := make([]C, length)
	for i := 0; i < length; i++ {
		vv := v.Index(i)
		tt := vv.Kind()
		switch tt {
		case reflect.String:
			tmp[i] = NewC(convertTo(vv.String()))
		case reflect.Int:
			tmp[i] = NewC(vv.Int())
		case reflect.Float64:
			tmp[i] = NewC(vv.Float())
		default:
			tmp[i] = NewC(vv.Interface())
		}
		tc := tmp[i].Type()
		_, ok = df.Types[col][tc]
		if ok == false {
			df.Types[col][tc] = 0
		}
		df.Types[col][tc]++
	}
	df.Df[col] = tmp
	df.NbLines = length
}

// SetMatrix is a function to add several columns of data
// m argument should be a slice of slice of something
// Example: df.SetMatrix([]string{"foo","bar"}, []interface{}{[]int{1,2,3,...},[]string{"a","b","c",...}})
func (df *DataFrame) SetMatrix(m interface{}, cols ...string) {
	v := reflect.ValueOf(m)
	t := v.Kind()
	if t != reflect.Slice {
		fmt.Printf("Error: Need to Use a Slice in 2nd arg instead of [%v]\n", t)
		return
	}
	length := v.Len()
	for i := 0; i < length; i++ {
		df.SetList(v.Index(i).Interface(), cols[i])
	}
}

// ToFloat function is a function to convert a column of Numeric in a slice of float64
func (df *DataFrame) ToFloat(col string) []float64 {
	ret := []float64{}
	_, ok := df.Df[col]
	if ok == false {
		fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		return nil
	}
	if len(df.Types[col]) != 1 {
		fmt.Printf("Error: Column [%v] is a multiple type can't convert to float64\n", col)
		return nil
	} else if len(df.Types[col]) == 1 {
		var t Type
		for k := range df.Types[col] {
			t = k
		}
		switch t {
		case NUMERIC:
			for i := 0; i < df.NbLines; i++ {
				ret = append(ret, float64(df.Df[col][i].(Numeric)))
			}
		default:
			fmt.Printf("Error: Column [%v] is a [%v] type can't convert to float64\n", col, t)
			return nil
		}
	}
	return ret
}

// ToMatrix is function to convert columns of Numeric in a matrix of float64
func (df *DataFrame) ToMatrix(cols ...string) ([]string, [][]float64) {
	ret := [][]float64{}
	columns := []string{}
	for _, col := range cols {
		r := df.ToFloat(col)
		if r == nil {
			continue
		}
		columns = append(columns, col)
		ret = append(ret, r)
	}
	return columns, ret
}

// FilterGT is function to filter if data in the specified column are greater than i argument
// Return of the function is the indexes of data wich are greater than i
func (df *DataFrame) FilterGT(col string, i interface{}) []int {
	c := NewC(i)
	_, ok := df.Df[col]
	if ok == false {
		fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		return nil
	}
	ret := []int{}
	for i := 0; i < df.NbLines; i++ {
		if df.Df[col][i].Great(c) {
			ret = append(ret, i)
		}
	}
	return ret
}

// FilterLT is a function similar to FilterGT for the lower than condition
func (df *DataFrame) FilterLT(col string, i interface{}) []int {
	c := NewC(i)
	_, ok := df.Df[col]
	if ok == false {
		fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		return nil
	}
	ret := []int{}
	for i := 0; i < df.NbLines; i++ {
		if df.Df[col][i].Less(c) {
			ret = append(ret, i)
		}
	}
	return ret
}

// FilterEQ is a function similar to FilterGT for the equal condition
func (df *DataFrame) FilterEQ(col string, i interface{}) []int {
	c := NewC(i)
	_, ok := df.Df[col]
	if ok == false {
		fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		return nil
	}
	ret := []int{}
	for i := 0; i < df.NbLines; i++ {
		if df.Df[col][i].Equal(c) {
			ret = append(ret, i)
		}
	}
	return ret
}

// SelectByIndex make a full copy of dataframe for the rows indexes specified
func (df *DataFrame) SelectByIndex(l []int) *DataFrame {
	if len(l) == 0 {
		fmt.Println("Error: No indexes availables")
		return nil
	}

	dfs := &DataFrame{}
	dfs.Df = map[string][]C{}
	dfs.Types = map[string]map[Type]int{}

	for _, col := range df.Columns {
		_, ok := df.Df[col]
		if ok == false {
			fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		} else {
			dfs.Columns = append(dfs.Columns, col)
			dfs.Df[col] = []C{}
			dfs.Types[col] = map[Type]int{}
		}
	}
	for _, v := range l {
		if v > df.NbLines-1 {
			fmt.Printf("Warning: Index [%v] out of range\n", v)
			continue
		}
		for _, col := range dfs.Columns {
			c := df.Df[col][v]
			dfs.Df[col] = append(dfs.Df[col], c) //To do: Check if not out of range
			switch c.(type) {
			case Numeric:
				_, ok := dfs.Types[col][NUMERIC]
				if ok == false {
					dfs.Types[col][NUMERIC] = 0
				}
				dfs.Types[col][NUMERIC]++
			case String:
				_, ok := dfs.Types[col][STRING]
				if ok == false {
					dfs.Types[col][STRING] = 0
				}
				dfs.Types[col][STRING]++
			case Time:
				_, ok := dfs.Types[col][TIME]
				if ok == false {
					dfs.Types[col][TIME] = 0
				}
				dfs.Types[col][TIME]++
			default:
				_, ok := dfs.Types[col][NAN]
				if ok == false {
					dfs.Types[col][NAN] = 0
				}
				dfs.Types[col][NAN]++
			}
		}
	}
	dfs.NbLines = len(dfs.Df[dfs.Columns[0]])
	if len(dfs.Columns) != 0 {
		return dfs
	}
	return nil
}

// AND function looks for common indexes between the two arguments of indexes
func AND(l1 []int, l2 []int) []int {
	d1 := map[int]bool{}
	d2 := map[int]bool{}
	for _, v := range l1 {
		d1[v] = true
	}
	for _, v := range l2 {
		d2[v] = true
	}
	ret := []int{}
	for k := range d1 {
		_, ok := d2[k]
		if ok {
			ret = append(ret, k)
			delete(d1, k)
			delete(d2, k)
		}
	}
	for k := range d2 {
		_, ok := d1[k]
		if ok {
			ret = append(ret, k)
			delete(d1, k)
			delete(d2, k)
		}
	}
	sort.Ints(ret)
	return ret
}

// OR function looks for indexes that are in first argument or second aregument of indexes
func OR(l1 []int, l2 []int) []int {
	d1 := map[int]bool{}
	d2 := map[int]bool{}
	for _, v := range l1 {
		d1[v] = true
	}
	for _, v := range l2 {
		d2[v] = true
	}
	ret := []int{}
	for k := range d1 {
		ret = append(ret, k)
		delete(d1, k)
		delete(d2, k)
	}
	for k := range d2 {
		ret = append(ret, k)
		delete(d1, k)
		delete(d2, k)
	}
	sort.Ints(ret)
	return ret
}

// Apply function apply a function on dataframe
// To Do
func (df *DataFrame) Apply(f func(c C) C) {
	for i := 0; i < df.NbLines; i++ {
		for _, col := range df.Columns {
			df.Df[col][i] = f(df.Df[col][i])
		}
	}
}
*/
func checkerr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
