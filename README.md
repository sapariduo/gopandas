Gopandas is an experimental golang package that aims to use "R/Python DataFrames".
This package is full inspirated by the pandas package in python and try to reproduce some of his functionalities as:

- Have a structure of data organized in columns (called a DataFrame)
- [x] Create a DataFrame from a csv file (done)
- [ ] Convert a DataFrame into csv and json file (to do)
- [x] Selection by columns
- [x] Selection by indices
- [x] Filtering by logical operators
- [ ] Group by
- [ ] Join
- [x] Apply functions
- [ ] Statistical functions

This package has been develop with the intention to be usable by machine learning packages in go.

**Output representation of a dataframe use the [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) package**

**/!\ THIS PACKAGE IS STILL UNDER HEAVY DEVELOPMENT /!\ **

# Installation

```
go get github.com/olekukonko/tablewriter
go get github.com/fmarmol/gopandas
```

# C interface
All gopandas types implements the C interface

```
type C interface {
    Add(C) C
    Sub(C) C
    Mul(C) C
    Div(C) C
    Mod(C) C
    Great(C) bool
    Less(C) bool
    Equal(C) bool
    NotEqual(C) bool
}
```

# Series

Series are basically simple ordered maps of values with indices as keys

**With series you can use :**

- Min/Max
- Sum/Mean
- Sort/Reverse
- Add/Sub/Mul/Div/Mod between 2 series
- Get the values
- Get the indices
- Get the number of occurences for each value
- Types/Type
- ReIndex

## Initialize series
Series can be constructed with maps or slices of values. Maps are used as input if there is a need to specify the indices of the values. Indices are empty interfaces so feel free to use what you want. If slices are used the indices will be the positions in the slices of the values.
Gopandas will convert automatically all types that are compatibles with the gopandas types (see section types). 

```go
import
    (
    "fmt"
    "github.com/fmarmol/gopandas/series"
)

func main() {
    s1 := series.New([]int{1,2,3})
    s2 := series.New(map[series.Index]interface{}{"one":1, 2:"two"})
    fmt.Print(s1)
    fmt.Print(s2)
}

//output:
Series:{0:1, 1:2, 2:3}
Series:{one:1, 2:two}
```

## Usage of series

- You can check types insides a series with Types method, this will give you a simple summary

```go
fmt.Println(s1.Types())
fmt.Println(s2.Types())
//output:
map[numeric:3]
map[string:1 numeric:1]
```

If you want just the unique result of the type of the series, use the Type method:
```go
fmt.Println(s1.Type())
fmt.Println(s2.Type())
//outpout:
numeric
multiple
```

- You can have the number of occurences by values inside a series with ValuesCount method
```go
fmt.Println(series.New([]float64{1.1,1.2,1.3,1.4,1.3,1.2}).ValuesCount())
//output:
map[1.1:1 1.3:2 1.4:1 1.2:2]
```

- You can apply function on series with method Apply. You need to pass a func(types.C) types.C, and it returns
a new series with values modified by the function:

```go
fmt.Print(
    series.New([]float64{1.1, 1.2}).Apply(
        func(c types.C) types.C {
            return c.Add(types.Numeric(1))
        }),
)
//output:
Series:{0:2.1, 1:2.2}
```

- You can also do basic operations between 2 series as for example:
```go
fmt.Print(series.New([]int{1,2,3}).Sub(series.New([]int{1,2,3})))

//ouput:
Series:{0:0, 1:0, 2:0}
```

- You can calculate Sum and Mean of a series:
```go
sum := s1.Sum()
mean := s1.Mean()
```

- You can sort the series either by ascending or descending order. The Sort and Reverse methods modifies series instead of returning a new one
```go
s := series.New([]float64{1.1, 1.2, -1.0})
s.Sort()
fmt.Print(s)
s.Reverse()
fmt.Print(s)
//output
Series:{2:-1, 0:1.1, 1:1.2}
Series:{1:1.2, 0:1.1, 2:-1}
```

## Usage of DataFrames

- You can create a dataFrame from scratch with the New method:

```go
df := dataframes.New([]string{"A", "B"}, []*series.Series{
    series.New([]int{1, 2, 3, 4}),
    series.New([]string{"one", "two", "three", "four"}),
})
fmt.Print(df)
//output
+---------+---------+--------+
|  Index  |    A    |   B    |
+---------+---------+--------+
|       0 |       1 |    one |
|       1 |       2 |    two |
|       2 |       3 |  three |
|       3 |       4 |   four |
+---------+---------+--------+
| COUNT:4 | numeric | string |
+---------+---------+--------+
```

- You can get a simple statistical description by using the Describe method:

```go
fmt.Print(df.Describe())
//ouput
+---------+---------+----------+
|  Index  |    A    |    B     |
+---------+---------+----------+
|     min |       1 |     four |
|     max |       4 |      two |
|    mean |     2.5 |      Nan |
|   count |       4 |        4 |
+---------+---------+----------+
| COUNT:4 | numeric | multiple |
+---------+---------+----------+
```

- It's possible to select only a subset of columns with the Select method:
```go
df := dataframes.New([]string{"A", "B", "C"}, []*series.Series{
    series.New([]int{1, 2, 3, 4}),
    series.New([]string{"one", "two", "three", "four"}),
    series.New([]interface{}{time.Now(), 2, "3", 1.1}),
})
fmt.Print(df.Select("A", "B"))
//output
+---------+---------+--------+
|  Index  |    A    |   B    |
+---------+---------+--------+
|       0 |       1 |    one |
|       1 |       2 |    two |
|       2 |       3 |  three |
|       3 |       4 |   four |
+---------+---------+--------+
| COUNT:4 | numeric | string |
+---------+---------+--------+
```

- It's easy to add series with the AddSeries method:

```go
df.AddSeries("D", series.New([]float64{4.4, 3.3, 2.2, 1.1}))
//output
+---------+---------+--------+---------------------------+---------+
|  Index  |    A    |   B    |             C             |    D    |
+---------+---------+--------+---------------------------+---------+
|       0 |       1 |    one | 2017-06-08T12:26:44+02:00 |     4.4 |
|       1 |       2 |    two |                         2 |     3.3 |
|       2 |       3 |  three |                         3 |     2.2 |
|       3 |       4 |   four |                       1.1 |     1.1 |
+---------+---------+--------+---------------------------+---------+
| COUNT:4 | numeric | string |         multiple          | numeric |
+---------+---------+--------+---------------------------+---------+
```
- Filters are easy to use:

```go
mask := df.FilterLT("A", 2).Or(df.FilterGT("D", 2))
fmt.Print(df.SelectByIndex(mask))
//output
+---------+---------+--------+---------------------------+---------+
|  Index  |    A    |   B    |             C             |    D    |
+---------+---------+--------+---------------------------+---------+
|       0 |       1 |    one | 2017-06-08T12:30:28+02:00 |     4.4 |
|       1 |       2 |    two |                         2 |     3.3 |
|       2 |       3 |  three |                         3 |     2.2 |
+---------+---------+--------+---------------------------+---------+
| COUNT:3 | numeric | string |         multiple          | numeric |
+---------+---------+--------+---------------------------+---------+
```

- You can also apply functions on a whole dataFrame with the Apply method:

```go
df.Apply(func(c types.C) types.C {
    return c.Div(types.Numeric(2))
})
fmt.Print(df)
//output
+---------+---------+-----+---------------------------+---------+
|  Index  |    A    |  B  |             C             |    D    |
+---------+---------+-----+---------------------------+---------+
|       0 |     0.5 | Nan | 2017-06-08T12:36:17+02:00 |     2.2 |
|       1 |       1 | Nan |                         1 |    1.65 |
|       2 |     1.5 | Nan |                       Nan |     1.1 |
+---------+---------+-----+---------------------------+---------+
| COUNT:3 | numeric | nan |         multiple          |+---------+---------+-----+---------------------------+---------+
|  Index  |    A    |  B  |             C             |    D    |
+---------+---------+-----+---------------------------+---------+
|       0 |     0.5 | Nan | 2017-06-08T12:36:17+02:00 |     2.2 |
|       1 |       1 | Nan |                         1 |    1.65 |
|       2 |     1.5 | Nan |                       Nan |     1.1 |
+---------+---------+-----+---------------------------+---------+
| COUNT:3 | numeric | nan |         multiple          | numeric |
+---------+---------+-----+---------------------------+---------+
 numeric |
+---------+---------+-----+---------------------------+---------+

```


# Types

Gopandas manages few types representing the most common types:

- Time for time (constructed over time.Time)
- Numeric for numbers (constructed over float64)
- String for string (constructed over string)
- Nan for Not A Number (constructed over string)

**Note**: All Types of Gopandas implement the gopandas types.C interface
