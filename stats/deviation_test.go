package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
	"reflect"
	"testing"
)

func TestStandardDeviation(t *testing.T) {
	s1 := series.New([]float64{-5, -1, 1.1, 2, 3, 3, 4, 6, 7, 7, 10, 17})

	type args struct {
		input *series.Series
	}
	tests := []struct {
		name     string
		args     args
		wantSdev types.C
		wantErr  bool
	}{
		{name: "Population Standard Deviation",
			args:     args{input: s1},
			wantSdev: types.Numeric(5.356531578881577),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSdev, err := StandardDeviation(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StandardDeviation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdev, tt.wantSdev) {
				t.Errorf("StandardDeviation() = %v, want %v", gotSdev, tt.wantSdev)
			}
		})
	}
}

func TestStandardDeviationSample(t *testing.T) {
	s1 := series.New([]float64{-5, -1, 1.1, 2, 3, 3, 4, 6, 7, 7, 10, 17})

	type args struct {
		input *series.Series
	}
	tests := []struct {
		name     string
		args     args
		wantSdev types.C
		wantErr  bool
	}{
		{name: "Sample Standard Deviation",
			args:     args{input: s1},
			wantSdev: types.Numeric(5.594714767826268),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSdev, err := StandardDeviationSample(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StandardDeviationSample() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdev, tt.wantSdev) {
				t.Errorf("StandardDeviationSample() = %v, want %v", gotSdev, tt.wantSdev)
			}
		})
	}
}
