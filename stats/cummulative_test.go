package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/sapariduo/gopandas/dataframes"
	"github.com/sapariduo/gopandas/series"
)

func TestCummulative(t *testing.T) {
	present := time.Now()
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)
	nextfuture := time.Now().Add(48 * time.Hour)
	df := dataframes.NewEmpty()

	s1 := series.New([]float64{5, 10, 20, 30})
	s2 := series.New([]float64{3, 4, 5, 7})
	saxis := series.New([]string{"a4", "a3", "a2", "a1"})
	stime := series.New([]time.Time{present, past, nextfuture, future})
	fmt.Println(stime)
	df.AddSeries("id", saxis)
	df.AddSeries("sales", s1)
	df.AddSeries("margin", s2)
	df.AddSeries("periode", stime)

	type args struct {
		axis      string
		datetime  string
		dataVal   []string
		dataframe *dataframes.DataFrame
	}
	tests := []struct {
		name    string
		args    args
		want    *dataframes.DataFrame
		wantErr bool
	}{
		{
			name:    "test single",
			args:    args{axis: "id", datetime: "periode", dataVal: []string{"margin", "sales"}, dataframe: df},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fmt.Println(tt.args.dataframe.Df[tt.args.axis])
			fmt.Println(tt.args.dataframe.Df["margin"])
			fmt.Println(tt.args.dataframe.Df[tt.args.datetime])
			got, err := Cummulative(tt.args.axis, tt.args.datetime, tt.args.dataVal, tt.args.dataframe)
			if (err != nil) != tt.wantErr {
				t.Errorf("CumulativeSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
