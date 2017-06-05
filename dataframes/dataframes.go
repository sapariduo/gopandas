package dataframes

import (
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
	Indices series.Indices
	Df      map[string]*series.Series
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
// Can be a float or a time in RFC3339 format {
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

func (df *DataFrame) ReIndex(indices series.Indices) {
	for _, col := range df.Columns {
		df.Df[col].ReIndex(indices)
	}
	df.Indices = indices
}

func (df *DataFrame) Len() int {
	if df == nil {
		fmt.Println("DataFrame is nil")
		return -1
	}
	return len(df.Indices)
}

func (df *DataFrame) isEmpty() bool {
	if len(df.Columns) == 0 && df.Len() == 0 {
		return true
	}
	return false
}

func (df *DataFrame) AddSeries(col string, s *series.Series) error {
	if df == nil {
		return fmt.Errorf("DataFrame is nil")
	}
	if df.isEmpty() {
		df.Columns = append(df.Columns, col)
		df.Df[col] = s
		df.Indices = s.GetIndices()
		df.NbLines = s.Len()
		return nil
	}
	for _, c := range df.Columns {
		if c == col {
			return fmt.Errorf("Error: Column %v already exists", col)
		}
	}
	if s.Len() != df.Len() {
		return fmt.Errorf("Error: lengths are not compatible")
	}
	for _, index := range df.Indices {
		if _, ok := s.Get(index); !ok {
			return fmt.Errorf("Error: Index: %v doesn't exist in series", index)
		}
	}
	df.Columns = append(df.Columns, col)
	df.Df[col] = s
	return nil
}

// Create a empty dataframe
func NewEmpty() *DataFrame {
	return &DataFrame{Columns: []string{}, Indices: series.Indices{}, Df: map[string]*series.Series{}, NbLines: 0}
}

func New(columns []string, ss []*series.Series) *DataFrame {
	if len(columns) != len(ss) {
		fmt.Println("Error: lenght of columns is not equal to length of series")
		return nil
	}
	df := NewEmpty()
	for i, c := range columns {
		if err := df.AddSeries(c, ss[i]); err != nil {
			log.Println(err)
			return nil
		}
	}
	return df
}

func (df *DataFrame) Describe() *DataFrame {
	if df == nil {
		return nil
	}
	ret := NewEmpty()
	for _, c := range df.Columns {
		ret.AddSeries(c, series.New(map[series.Index]interface{}{
			"min":   df.Df[c].Min(),
			"max":   df.Df[c].Max(),
			"mean":  df.Df[c].Mean(),
			"count": df.Len(),
		}))
	}
	return ret
}

// NewFromCSV is a function to create of a dataframe from a CSV file
func NewFromCSV(c *ConfigDataFrame) *DataFrame {
	fd, err := os.Open(c.File)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer fd.Close()

	df := DataFrame{}
	df.Df = map[string]*series.Series{}
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
			df.Df[col] = series.NewEmpty()
		}
	} else {
		if c.Index {
			df.Columns = make([]string, len(firstline[1:]))
			df.Indices = append(df.Indices, firstline[0])
			for i, c := range firstline[1:] {
				col := fmt.Sprintf("%v", i)
				df.Columns[i] = col
				df.Df[col] = series.NewEmpty()
				df.Df[col].Set(firstline[0], types.NewC(convertTo(c)))
			}
		} else {
			df.Columns = make([]string, len(firstline))
			df.Indices = append(df.Indices, index)
			for i, c := range firstline {
				col := fmt.Sprintf("%v", i)
				df.Columns[i] = col
				df.Df[col] = series.NewEmpty()
				df.Df[col].Set(index, types.NewC(convertTo(c)))
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
				df.Df[col].Set(line[0], v)
			} else {
				v := types.NewC(convertTo(line[icol]))
				df.Df[col].Set(index, v)
			}
		}
		index++
		df.NbLines++
	}
	return &df
}

func (df *DataFrame) String() string {
	if df == nil {
		return "Nil dataFrame"
	}
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
			v, _ := df.Df[col].Get(index)
			l = append(l, fmt.Sprintf("%v", v))
		}
		table.Append(l)
	}
	footer := make([]string, len(header))
	footer[0] = fmt.Sprintf("COUNT:%v", df.NbLines)
	for i, col := range df.Columns {
		m := df.Df[col].Type()
		footer[i+1] = string(m)
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
	if df == nil {
		fmt.Println("DataFrame is nil")
		return nil
	}
	dfs := &DataFrame{}
	dfs.NbLines = df.NbLines
	dfs.Df = map[string]*series.Series{}
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
}*/

// FilterGT is function to filter if data in the specified column are greater than i argument
// Return of the function is the indexes of data wich are greater than i
func (df *DataFrame) FilterGT(col string, i interface{}) series.Indices {
	return df.Df[col].FilterGT(i)
}

func (df *DataFrame) FilterGTEQ(col string, i interface{}) series.Indices {
	return df.Df[col].FilterGTEQ(i)
}

// FilterLT is a function similar to FilterGT for the lower than condition
func (df *DataFrame) FilterLT(col string, i interface{}) series.Indices {
	return df.Df[col].FilterLT(i)
}

func (df *DataFrame) FilterLTEQ(col string, i interface{}) series.Indices {
	return df.Df[col].FilterLTEQ(i)
}

// FilterEQ is a function similar to FilterGT for the equal condition
func (df *DataFrame) FilterEQ(col string, i interface{}) series.Indices {
	return df.Df[col].FilterEQ(i)
}

func (df *DataFrame) FilterNEQ(col string, i interface{}) series.Indices {
	return df.Df[col].FilterNEQ(i)
}

// SelectByIndex make a full copy of dataframe for the rows indexes specified
func (df *DataFrame) SelectByIndex(indices series.Indices) *DataFrame {
	if len(indices) == 0 {
		fmt.Println("Error: No indices available")
		return nil
	}
	ret := NewEmpty()
	for _, col := range df.Columns {
		ret.AddSeries(col, df.Df[col].Select(indices))
	}
	return ret
}

/*
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
}*/

// Apply function apply a function on dataframe
// To Do
func (df *DataFrame) Apply(f func(types.C) types.C) {
	for _, c := range df.Columns {
		df.Df[c] = df.Df[c].Apply(f)
	}
}

func checkerr(e error) {
	if e != nil {
		log.Panic(e)
	}
}
