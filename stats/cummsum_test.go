package stats

import (
	"fmt"
	"testing"
	"time"

	"github.com/sapariduo/gopandas/dataframes"
	"github.com/sapariduo/gopandas/series"
)

func TestCumulativeSum(t *testing.T) {
	present := time.Now()
	past := time.Now().Add(-24 * time.Hour)
	future := time.Now().Add(24 * time.Hour)
	nextfuture := time.Now().Add(48 * time.Hour)
	df := dataframes.NewEmpty()

	s1 := series.New([]float64{5, 10, 20, 30})
	s2 := series.New([]float64{3, 4, 5, 7})
	stime := series.New([]time.Time{present, past, nextfuture, future})
	fmt.Println(stime)
	df.AddSeries("sales", s1)
	df.AddSeries("margin", s2)

	err := df.AddSeries("tdate", stime)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(df)

	type args struct {
		data     *series.Series
		datetime *series.Series
	}
	tests := []struct {
		name    string
		args    args
		want    *dataframes.DataFrame
		wantErr bool
	}{
		{
			name: "sales test",
			args: args{data: df.Df["sales"], datetime: df.Df["tdate"]},
		},
		{
			name: "margin test",
			args: args{data: df.Df["margin"], datetime: df.Df["tdate"]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(tt.args.data)
			fmt.Println(tt.args.datetime)
			got, err := CumulativeSum(tt.args.data, tt.args.datetime)
			if (err != nil) != tt.wantErr {
				t.Errorf("CumulativeSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)

		})
	}
}
