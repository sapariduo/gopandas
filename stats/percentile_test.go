package stats

import (
	"fmt"
	"github.com/sapariduo/gopandas/series"
	"reflect"
	"testing"
)

func TestPercentile(t *testing.T) {
	s1 := series.New([]float64{3, 4, 4, 4, 7, 10, 11, 12, 14, 16, 17, 18})
	quartiles, _ := Quartile(s1) //quartile will produce percentile 25, 50 and 75 respectively

	type args struct {
		input   *series.Series
		percent float64
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}

	t1 := test{name: "Percentile Test", args: args{s1, 25}, wantErr: false}
	v := reflect.ValueOf(quartiles)
	fmt.Printf("%v with type %T\n", v, v)
	wantValues := make([]interface{}, v.NumField())
	for i := range wantValues {
		t.Run(t1.name, func(t *testing.T) {
			fmt.Printf("value of v of field i %+v, with type %T\n", v.Field(i), v.Field(i))
			got, err := Percentile(t1.args.input, t1.args.percent*(float64(i)+1))
			fmt.Println(got.Type())
			if (err != nil) != t1.wantErr {
				t.Errorf("Percentile() error = %v, wantErr %v", err, t1.wantErr)
				return
			}
			// q := fmt.Sprintf("Q%s", i+1)
			if !reflect.DeepEqual(got, v.Field(i)) {
				t.Errorf("Percentile() = %v, want %v", got, v.Field(i))
			}
		})
	}

}
