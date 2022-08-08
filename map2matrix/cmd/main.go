package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

func main() {
	if err := run(os.Args[1]); err != nil {
		fmt.Println("failed", err)
	}
}

type neighbourg struct {
	C int
	T map[string]int
	B map[string]int
	L map[string]int
	R map[string]int
}

func newNeighbourg() neighbourg {
	return neighbourg{
		C: 0,
		T: map[string]int{},
		B: map[string]int{},
		L: map[string]int{},
		R: map[string]int{},
	}
}

func run(filename string) error {
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer infile.Close()

	src, format, err := image.Decode(infile)
	if err != nil {
		return fmt.Errorf("while decoding img %s : %w", format, err)
	}

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	out := map[string]neighbourg{}
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := colorToString(src.At(x, y))
			at, ok := out[c]
			if !ok {
				at = newNeighbourg()
			}
			at.C++
			if x == 0 {
				at.L["left"] = at.L["left"] + 1
			} else {
				c := colorToString(src.At(x-1, y))
				at.L[c] = at.L[c] + 1
			}
			if x == w-1 {
				at.R["right"] = at.R["right"] + 1
			} else {
				c := colorToString(src.At(x+1, y))
				at.R[c] = at.R[c] + 1
			}

			if y == 0 {
				at.T["top"] = at.T["top"] + 1
			} else {
				c := colorToString(src.At(x, y-1))
				at.T[c] = at.T[c] + 1
			}
			if y == h-1 {
				at.B["bottom"] = at.B["bottom"] + 1
			} else {
				c := colorToString(src.At(x, y+1))
				at.B[c] = at.B[c] + 1
			}
			out[c] = at
		}
	}

	asJson, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("while json encoding result : %w", err)
	}
	fmt.Println(string(asJson))
	return nil
}

func colorToString(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("[%d,%d,%d]", r/257, g/257, b/257)
}
