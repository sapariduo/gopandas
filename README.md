Gopandas is an experimental golang package that aims to use "R/Python DataFrames".
This package is full inspirated by the pandas package in python and try to reproduce some of his functionalities as:

- Have a structure of data organized in columns (called a DataFrame)
- Create a DataFrame from a CSV or JSON file
- Convert a DataFrame in CSV or JSON file
- Selection by columns
- Filtering by logical operator as AND, OR, GREAT, LESS, EQUAL 
- Group by (TO DO)
- Join (TO DO)
- Apply functions
- Statistical functions (TO DO)

This package has been develop with the intention to be usable by machine learning packages in go.

**Output representation of a dataframe use the [github.com/olekukonko/tablewriter](https://github.com/olekukonko/tablewriter) package**

**/!\ THIS PACKAGE IS STILL UNDER DEVELOPMENT. A LOT OF FUNCTIONALITIES IS MISSING AND/OR COULD BE CHANGED /!\ **

# Installation

```
go get github.com/olekukonko/tablewriter
go get github.com/fmarmol/gopandas
```

# Initialize a DataFrame

The easiest way to create a new dataframe from scratch is to use the method SetMatrix, for example:
```go
package main

import (
	pd "github.com/fmarmol/gopandas"
	"fmt"
)

func main() {
	df := &pd.DataFrame{}
	df.SetMatrix([][]interface{}{{1,2,3},{"foo","bar","baz"}},"A","B")
	fmt.Println(df)
}
```
```
+---------+---------+--------+
|  Index  |    A    |   B    |
+---------+---------+--------+
|       0 |       1 |    foo |
|       1 |       2 |    bar |
|       2 |       3 |    baz |
+---------+---------+--------+
| COUNT:3 | numeric | string |
+---------+---------+--------+

```
**Note**: The output representation comes with an additional "Index" column. This column is purely aesthetic and doesn't exist
in the DataFrame structure.

The last line contains two informations:

1. The number of lines contained in the DataFrame
2. The "Type" of each column (see Types section)


# Read From a CSV file

Example, with a CSV file:
```
$> cat example.csv

A,B,C
1,foo,2016-03-24 15:50:05
2,bar,2000-01-01 12:12:12
3,baz,2020-12-20 23:43:54

```

``` go
import(
	 pd "github.com/fmarmol/gopandas"
	"fmt"
)

c := &pd.ConfigDataFrame{
	File: "example.csv",
	Header: true, // Only used for CSV file. 
	Sep: ',', // Rune Separator used only for CSV file
}

df := NewDataFrameCSV(c)
fmt.Println(df)
```

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       2 |    bar | 2000-01-01T12:12:12Z |
|       2 |       3 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:3 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```
**Note**: You can see that a Time Type exists in the gopandas package. For now the layout for a date to be parsed is %Y-%m-%d %H:%M:%S.

# Selection by column's names

You can easily select you columns thanks to their name:

```go
df2 := df.Select("A")
df2 := df.Select("A","B")
// etc...
```

**Important Note**: The select by column functionality doesn't make a copy of the dataframe. Also any changes in the dataframe selected
will be apply also on the original dataframe.
If you want make a copy you can do:

```go
newdf := df.Select("A").Copy()
```

If you make a selection with a column's name that doesn'tt exist a simple warning will be printed and this selection will be just skipped
and go to the next selection.

# Selection by index row

You can select specifics rows of a DataFrame with SelectByIndex method. This method was initially created to be used in combination with
filtering methods.
With a DataFrame as:

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       2 |    bar | 2000-01-01T12:12:12Z |
|       2 |       3 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:3 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

```go
df = df.SelectByIndex([]int{0,2})
fmt.Println(df)
//produces
```

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       3 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:2 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

**Important Note**: The select by index functionality makes a full copy of the dataframe

# Add Columns

 You can add a column in a dataframe, you have to pass in first argument
a slice of something (could be interface, int, float, string or whatever) and in second argument the name of the new column. 
The construction of the new column will try to convert
element of the slice in one special type of gopandas package.


```go
df.SetList([]int{1,2,3},"A")
//or
df.SetList([]interface{}{"a",1,time.Time{},"B"}
```

The elements of the slice doesn't need to be the same type.

You can also set several new colums with the SetMatrix method in one shot:

```go
//With a initial dataframe as:
+---------+---------+--------+
|  Index  |    A    |   B    |
+---------+---------+--------+
|       0 |       1 |    one |
|       1 |       2 |    two |
|       2 |       3 |  three |
+---------+---------+--------+
| COUNT:3 | numeric | string |
+---------+---------+--------+

//You can set the matrix list:
df.SetMatrix([][]interface{}{{time.Now(), time.Time{}, time.Now()}, {1.1, 2, 3}}, "C", "D")
//And obtain:
+---------+---------+--------+---------------------------+---------+
|  Index  |    A    |   B    |             C             |    D    |
+---------+---------+--------+---------------------------+---------+
|       0 |       1 |    one | 2016-03-30T00:29:30+02:00 |     1.1 |
|       1 |       2 |    two |      0001-01-01T00:00:00Z |       2 |
|       2 |       3 |  three | 2016-03-30T00:29:30+02:00 |       3 |
+---------+---------+--------+---------------------------+---------+
| COUNT:3 | numeric | string |           time            | numeric |
+---------+---------+--------+---------------------------+---------+

```

# Filtering

You can apply filters on a DataFrame thanks to Filter* Methods.

- FilterGT (Greater than ≡ >)
- FilterLT (Lower than ≡ <)
- FilterEQ (Equal ≡ ==)

This methods returns the a slice on index corresponding to the rows number which are true with the condition.
This methods can be used in combinaison with AND/OR methods to construct more complex conditions

With a DataFrame as:

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       2 |    bar | 2000-01-01T12:12:12Z |
|       2 |       3 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:3 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

``` go
df = df.SelectByIndex(pd.OR(df.FilterEQ("A",1),df.FilterEQ("B","bar")))
fmt.Println(df)
//and get:
```
```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       2 |    bar | 2000-01-01T12:12:12Z |
+---------+---------+--------+----------------------+
| COUNT:2 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

If you try the same thing with AND method, you have an impossible condition in this case and you'll get something like:

```
Error: No indexes availables
<nil>

```

# Apply functions

You can apply functions on a DataFrame with the  Apply method. This method takes in argument a function on a pd.C interface an return a pd.C interface. 
Gopandas will try to apply the function of all elements of the DataFrame.

All gopandas Types implements the C interface

``` go
type C interface {
    Add(C) C
    Mul(C) C
    Div(C) C
    Mod(C) C
    Great(C) bool
    Less(C) bool
    Equal(C) bool
    NotEqual(C) bool
}
```

With a DataFrame as:

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       1 |    foo | 2016-03-24T15:50:05Z |
|       1 |       2 |    bar | 2000-01-01T12:12:12Z |
|       2 |       3 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:3 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

```go
df.Apply(func(c C) C {
	return c.Add(Numeric(1))
})
fmt.Println(df)
//And get:
```

```
+---------+---------+--------+----------------------+
|  Index  |    A    |   B    |          C           |
+---------+---------+--------+----------------------+
|       0 |       2 |    foo | 2016-03-24T15:50:05Z |
|       1 |       3 |    bar | 2000-01-01T12:12:12Z |
|       2 |       4 |    baz | 2020-12-20T23:43:54Z |
+---------+---------+--------+----------------------+
| COUNT:3 | numeric | string |         time         |
+---------+---------+--------+----------------------+
```

**Note**: It lacks a lot of behaviors on apply functions for each Type of gopandas

# Types

Gopandas manages few types representing the most common types:

- Time for time (constructed over time.Time)
- Numeric for numbers (constructed over float64)
- String for string (constructed over string)
- Nan for Not A Number (constructed over string)

**Note**: All Types of Gopandas implement the pd.C interface
