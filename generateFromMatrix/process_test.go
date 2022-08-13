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
		"one black pixel on 2x2": {
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
		"2x2 pattern on 3x3": {
			in: Constraint{
				PatternSize: 2,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0,1,1,0]": 1,
					"[1,0,0,1]": 1,
				},
			},
			w: 3, h: 3, patternSize: 2,
			out: []int{
				1, 0, 1,
				0, 1, 0,
				1, 0, 1,
			},
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

func TestCanBePlaced(t *testing.T) {
	tests := map[string]struct {
		written              []bool
		indexes, patterns    []int
		x, y, w, patternSize int
		exp                  bool
	}{
		"1x1 not written": {
			written:     []bool{false},
			indexes:     []int{0},
			patterns:    []int{0},
			x:           0,
			y:           0,
			w:           1,
			patternSize: 1,
			exp:         true,
		},
		"1x1 already written": {
			written:     []bool{true},
			indexes:     []int{0},
			patterns:    []int{0},
			x:           0,
			y:           0,
			w:           1,
			patternSize: 1,
			exp:         true,
		},
		"1x1 already written but not compatible": {
			written:     []bool{true},
			indexes:     []int{1},
			patterns:    []int{0},
			x:           0,
			y:           0,
			w:           1,
			patternSize: 1,
			exp:         false,
		},
		"2x2 already written ": {
			written:     []bool{true, true, false, false},
			indexes:     []int{1, 0, 0, 0},
			patterns:    []int{1, 0, 0, 1},
			x:           0,
			y:           0,
			w:           2,
			patternSize: 2,
			exp:         true,
		},
	}

	for title, tt := range tests {
		t.Run(title, func(t *testing.T) {
			result := canBePlaced(tt.written, tt.indexes, tt.patterns, tt.x, tt.y, tt.w, tt.patternSize)
			if result != tt.exp {
				t.Fatalf("failed got %v want %v", result, tt.exp)
			}
		})
	}
}
