package generateFromMatrix

import (
	"reflect"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := map[string]struct {
		in                Constraint
		w, h, patternSize int
		out               []int
	}{
		"one black pixel": {
			in: Constraint{
				PatternSize: 1,
				Index:       []string{"[0,0,0]"},
				Patterns: map[string]int{
					"[0]": 1,
				},
			},
			w: 2, h: 2, patternSize: 1,
			out: []int{0, 0, 0, 0},
		},
	}

	for title, tt := range tests {
		t.Run(title, func(t *testing.T) {
			HelperGenerate(t, tt.in, tt.patternSize, tt.w, tt.h, func() float64 { return 1.0 }, tt.out)
		})
	}
}

func HelperGenerate(t *testing.T, c Constraint, patternSize, w, h int, rand func() float64, exp []int) {
	t.Helper()
	result, err := generate(c, patternSize, w, h, rand)
	if err != nil {
		t.Fatalf("failed : %v", err)
	}
	if !reflect.DeepEqual(result, exp) {
		t.Fatalf("exp \n%#v\ngot \n%#v", exp, result)
	}
}
