package stats

import (
	"github.com/sapariduo/gopandas/series"
	"github.com/sapariduo/gopandas/types"
	"reflect"
	"testing"
)

func TestVarP(t *testing.T) {

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
		{name: "Population Variance",
			args:     args{input: s1},
			wantSdev: types.Numeric(28.692430555555557),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSdev, err := VarP(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdev, tt.wantSdev) {
				t.Errorf("VarP() = %v, want %v", gotSdev, tt.wantSdev)
			}
		})
	}
}

func TestVarS(t *testing.T) {
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
		{name: "Sample Variance",
			args:     args{input: s1},
			wantSdev: types.Numeric(31.300833333333333),
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSdev, err := VarS(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("VarS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSdev, tt.wantSdev) {
				t.Errorf("VarS() = %v, want %v", gotSdev, tt.wantSdev)
			}
		})
	}
}
