package gopandas

import (
	//"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"time"
)

// NewC if function with create a C interface element wich is contained in a dataframe.
// C can be Numeric to perform mathematical operations, Time to use datetime, String or Nan
func NewC(i interface{}) (C, Type) {
	var ret C
	var t Type
	switch i.(type) {
	case float64:
		ret = Numeric(i.(float64))
		t = NUMERIC
	case int:
		ret = Numeric(float64(i.(int)))
		t = NUMERIC
	case int64:
		ret = Numeric(float64(i.(int64)))
		t = NUMERIC
	case time.Time:
		ret = Time(i.(time.Time))
		t = TIME
	case string:
		ret = String(i.(string))
		t = STRING
	case Numeric:
		ret = i.(Numeric)
		t = NUMERIC
	case String:
		ret = i.(String)
		t = STRING
	default:
		ret = Nan("NaN")
		t = NAN
	}
	return ret, t
}

// DataFrame is the structure of a dataframe, Data are avaibable with Df attribute
// Final goal is to have all attributes in private.
type DataFrame struct {
	Columns []string
	Types   map[string]map[Type]int
	Df      map[string][]C
	NbLines int
}

// ConfigDataFrame is a structure of configuration to create a dataframe
type ConfigDataFrame struct {
	File       string
	Header     bool
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

// NewDataFrameJSON is a function to create of a dataframe from a JSON file
// Json data must be in the format of []map[string]interface{}
// No more checks are done
func (c *ConfigDataFrame) NewDataFrameJSON() *DataFrame {
	fd, err := os.Open(c.File)
	if err != nil {
		return nil
	}
	defer fd.Close()
	if err != nil {
		return nil
	}
	decoder := json.NewDecoder(fd)
	var ret []map[string]interface{}
	err = decoder.Decode(&ret)
	if err != nil {
		return nil
	}
	df := DataFrame{}
	df.Df = map[string][]C{}
	df.Types = map[string]map[Type]int{}
	df.NbLines = len(ret)
	for k := range ret[0] {
		df.Df[k] = make([]C, df.NbLines)
		df.Types[k] = map[Type]int{}
		df.Columns = append(df.Columns, k)
	}
	for index, r := range ret {
		for _, col := range df.Columns {
			_, ok := r[col]
			if ok == false {
				fmt.Printf("Error: Columns are not consistent\n")
				return nil
			}
			v, t := NewC(r[col])
			df.Df[col][index] = v
			_, ok = df.Types[col][t]
			if ok == false {
				df.Types[col][t] = 0
			}
			df.Types[col][t]++
		}
	}
	return &df
}

// NewDataFrameCSV is a function to create of a dataframe from a CSV file
func (c *ConfigDataFrame) NewDataFrameCSV() *DataFrame {
	fd, err := os.Open(c.File)
	if err != nil {
		return nil
	}
	df := DataFrame{}
	df.Df = map[string][]C{}
	df.Types = map[string]map[Type]int{}
	reader := csv.NewReader(fd)
	reader.Comma = c.Sep
	firstline, err := reader.Read()
	checkerr(err)
	lines, err := reader.ReadAll()
	checkerr(err)
	if c.Header {
		df.Columns = firstline
		for index, col := range df.Columns {
			df.Df[col] = make([]C, len(lines))
			df.Types[col] = map[Type]int{}
			v, t := NewC(convertTo(lines[0][index]))
			df.Df[col][0] = v
			_, ok := df.Types[col]
			if ok == false {
				df.Types[col][t] = 0
			}
			df.Types[col][t]++
		}
		df.NbLines++
		lines = lines[1:]
	} else {
		df.Columns = make([]string, len(firstline))
		for i := 0; i < len(firstline); i++ {
			col := fmt.Sprintf("%v", i)
			df.Columns[i] = col
			df.Df[col] = make([]C, len(lines)+1)
			df.Types[col] = map[Type]int{}
			v, t := NewC(convertTo(firstline[i]))
			df.Df[col][0] = v
			_, ok := df.Types[col][t]
			if ok == false {
				df.Types[col][t] = 0
			}
			df.Types[col][t]++
		}
		df.NbLines++
	}
	for nb, line := range lines {
		for index, col := range df.Columns {
			v, t := NewC(convertTo(line[index]))
			df.Df[col][nb+1] = v
			_, ok := df.Types[col]
			if ok == false {
				df.Types[col][t] = 0
			}
			df.Types[col][t]++
		}
		df.NbLines++
	}
	err = fd.Close()
	checkerr(err)
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
	for i := 0; i < df.NbLines; i++ {
		l := []string{fmt.Sprintf("%v", i)}
		for _, col := range df.Columns {
			l = append(l, fmt.Sprintf("%v", df.Df[col][i]))
		}
		table.Append(l)
	}
	footer := make([]string, len(header))
	footer[0] = fmt.Sprintf("COUNT:%v", df.NbLines)
	for i, col := range df.Columns {
		m := df.Types[col]
		if len(m) == 0 {
			log.Panicf("Error Types are not defined for [%v] column\n", col)
		} else if len(m) > 1 {
			footer[i+1] = MULTI
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
	dfs.Df = map[string][]C{}
	dfs.Types = df.Types
	for _, col := range cols {
		_, ok := df.Df[col]
		if ok == false {
			fmt.Printf("Warning: Column name [%v] doesn't exist\n", col)
		} else {
			dfs.Columns = append(dfs.Columns, col)
			dfs.Df[col] = df.Df[col]
			dfs.Types[col] = df.Types[col]
		}
	}
	if len(dfs.Columns) != 0 {
		return dfs
	}

	return nil

}

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
		var tc Type
		switch tt {
		case reflect.String:
			tmp[i], tc = NewC(convertTo(vv.String()))
		case reflect.Int:
			tmp[i], tc = NewC(vv.Int())
		case reflect.Float64:
			tmp[i], tc = NewC(vv.Float())
		default:
			tmp[i], tc = NewC(vv.Interface())
		}
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
	c, _ := NewC(i)
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
	c, _ := NewC(i)
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
	c, _ := NewC(i)
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

func checkerr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
