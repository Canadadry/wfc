package map2matrix

import (
	"app/generateFromMatrix"
	"image/color"
	"reflect"
	"testing"
)

type FakeImg struct {
	Pixels [][]int
	Colors map[int]string
	t      *testing.T
}

func (fi FakeImg) At(x, y int) color.Color {
	fi.t.Helper()
	if y >= len(fi.Pixels) {
		fi.t.Fatalf("fake img height = %d want pixel at %d", len(fi.Pixels), y)
	}
	if x >= len(fi.Pixels[y]) {
		fi.t.Fatalf("fake img with = %d want pixel at %d", len(fi.Pixels[y]), x)
	}
	pix := fi.Pixels[y][x]
	str, ok := fi.Colors[pix]
	if !ok {
		fi.t.Fatalf("color index %d not found", pix)
	}
	c, err := generateFromMatrix.ColorFromString(str)
	if err != nil {
		fi.t.Fatalf("cannot convert '%s' to color", str)
	}
	return c
}

func TestExtract(t *testing.T) {
	tests := map[string]struct {
		in                FakeImg
		w, h, patternSize int
		out               Result
	}{
		"one black pixel": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
				},
			},
			w: 1, h: 1, patternSize: 1,
			out: Result{
				PatternSize: 1,
				Index:       []string{"[0,0,0]"},
				Patterns: map[string]int{
					"[0]": 1,
				},
			},
		},
		"two pixels": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0, 1},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
					1: "[1,1,1]",
				},
			},
			w: 2, h: 1, patternSize: 1,
			out: Result{
				PatternSize: 1,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0]": 1,
					"[1]": 1,
				},
			},
		},
		"a 2x2 alternate pixels": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0, 1},
					[]int{1, 0},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
					1: "[1,1,1]",
				},
			},
			w: 2, h: 2, patternSize: 1,
			out: Result{
				PatternSize: 1,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0]": 2,
					"[1]": 2,
				},
			},
		},
		"a 2x2 alternate pixels with pattern size of 2": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0, 1},
					[]int{1, 0},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
					1: "[1,1,1]",
				},
			},
			w: 2, h: 2, patternSize: 2,
			out: Result{
				PatternSize: 2,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0,1,1,0]": 1,
				},
			},
		},
		"a 4x4 alternate pixels with pattern size of 2": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0, 1, 0, 1},
					[]int{1, 0, 1, 0},
					[]int{0, 1, 0, 1},
					[]int{1, 0, 1, 0},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
					1: "[1,1,1]",
				},
			},
			w: 4, h: 4, patternSize: 2,
			out: Result{
				PatternSize: 2,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0,1,1,0]": 5,
					"[1,0,0,1]": 4,
				},
			},
		},
		"a 4x4 alternate pixels with pattern size of 3": {
			in: FakeImg{
				Pixels: [][]int{
					[]int{0, 1, 0, 1},
					[]int{1, 0, 1, 0},
					[]int{0, 1, 0, 1},
					[]int{1, 0, 1, 0},
				},
				Colors: map[int]string{
					0: "[0,0,0]",
					1: "[1,1,1]",
				},
			},
			w: 4, h: 4, patternSize: 3,
			out: Result{
				PatternSize: 3,
				Index:       []string{"[0,0,0]", "[1,1,1]"},
				Patterns: map[string]int{
					"[0,1,0,1,0,1,0,1,0]": 2,
					"[1,0,1,0,1,0,1,0,1]": 2,
				},
			},
		},
	}

	for title, tt := range tests {
		t.Run(title, func(t *testing.T) {
			HelperExtract(t, tt.in, tt.w, tt.h, tt.patternSize, tt.out)
		})
	}
}

func HelperExtract(t *testing.T, img FakeImg, w, h, patternSize int, exp Result) {
	img.t = t
	result := extract(img, w, h, patternSize)
	if !reflect.DeepEqual(result, exp) {
		t.Fatalf("exp \n%#v\ngot \n%#v", exp, result)
	}
}
