package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
	"reflect"
	"testing"
)

func TestQuartile(t *testing.T) {

	s1 := series.New([]float64{-5, -1, 1.1, 2, 3, 3, 4, 6, 7, 7, 10, 17})
	s2 := series.New([]float64{3, 4, 4, 4, 7, 10, 11, 12, 14, 16, 17, 18})

	quartiles := Quartiles{types.Numeric(1.55), types.Numeric(3.5), types.Numeric(7)}
	iqr := types.Numeric(5.45)
	midhinge := types.Numeric(4.275)
	median, _ := Median(s2)
	midhinge2, _ := Midhinge(s2)
	trimean := (median.Add(midhinge2)).Div(types.Numeric(2))

	type args struct {
		input *series.Series
	}

	type test struct {
		name    string
		args    args
		wantErr bool
	}

	t1 := test{name: "Quartiles Test", args: args{s1}, wantErr: false}
	t2 := test{name: "InterQuartileRange Test", args: args{s1}, wantErr: false}
	t3 := test{name: "Midhinge Test", args: args{s1}, wantErr: false}
	t4 := test{name: "Trimean Test", args: args{s2}, wantErr: false}

	t.Run(t1.name, func(t *testing.T) {
		got, err := Quartile(t1.args.input)
		if (err != nil) != t1.wantErr {
			t.Errorf("Quartile() error = %v, wantErr %v", err, t1.wantErr)
			return
		}
		if !reflect.DeepEqual(got, quartiles) {
			t.Errorf("Quartile() = %v, want %v", got, quartiles)
		}
	})

	t.Run(t2.name, func(t *testing.T) {
		got, err := InterQuartileRange(t2.args.input)
		if (err != nil) != t2.wantErr {
			t.Errorf("InterQuartileRange() error = %v, wantErr %v", err, t2.wantErr)
			return
		}
		if !reflect.DeepEqual(got, iqr) {
			t.Errorf("InterQuartileRange() = %v, want %v", got, iqr)
		}
	})

	t.Run(t3.name, func(t *testing.T) {
		got, err := Midhinge(t3.args.input)
		if (err != nil) != t3.wantErr {
			t.Errorf("Midhinge() error = %v, wantErr %v", err, t3.wantErr)
			return
		}
		if !reflect.DeepEqual(got, midhinge) {
			t.Errorf("Midhinge() = %v, want %v", got, midhinge)
		}
	})

	t.Run(t4.name, func(t *testing.T) {
		got, err := Trimean(t4.args.input)
		if (err != nil) != t4.wantErr {
			t.Errorf("Trimean() error = %v, wantErr %v", err, t4.wantErr)
			return
		}
		if !reflect.DeepEqual(got, trimean) {
			t.Errorf("Trimean() = %v, want %v", got, trimean)
		}
	})

}
