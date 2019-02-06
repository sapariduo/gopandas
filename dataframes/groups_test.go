package dataframes

import (
	"fmt"
	"gopandas/series"
	"testing"
)

func TestGroups_Max(t *testing.T) {
	df := New([]string{"working", "person", "unit_deps", "city"}, []*series.Series{workHour, person, department, city})
	fmt.Println(df)
	dfg := df.GroupBy("unit_deps", "city")
	fmt.Println(dfg)
	fmt.Println(dfg.Max())
	fmt.Println(dfg)
}

func TestGroups_Sum(t *testing.T) {
	df := New([]string{"working", "person", "unit", "city"}, []*series.Series{workHour, person, department, city})
	dfg := df.GroupBy("unit", "city")
	fmt.Println(dfg.Sum())
}

func TestGroups_Info(t *testing.T) {
	df := New([]string{"working", "person", "unit", "city"}, []*series.Series{workHour, person, department, city})
	dfg := df.GroupBy("unit", "city")
	fmt.Println(dfg.Info())
}

func TestGroups_Min(t *testing.T) {
	df := New([]string{"working", "person", "unit", "city"}, []*series.Series{workHour, person, department, city})
	dfg := df.GroupBy("unit", "city")
	fmt.Println(dfg.Min())
}

func TestGroups_Count(t *testing.T) {
	df := New([]string{"working", "person", "unit", "city"}, []*series.Series{workHour, person, department, city})
	dfg := df.GroupBy("unit", "city")
	fmt.Println(dfg.Count())
}

func TestGroups_Mean(t *testing.T) {
	df := New([]string{"working", "person", "unit", "city"}, []*series.Series{workHour, person, department, city})
	dfg := df.GroupBy("unit", "city")
	fmt.Println(dfg.Mean())
}
