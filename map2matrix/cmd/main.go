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

type Pixel struct {
	Color        string
	Density      float64
	NeighbourgsT []Neighbourg
	NeighbourgsB []Neighbourg
	NeighbourgsL []Neighbourg
	NeighbourgsR []Neighbourg
}

type Neighbourg struct {
	Index   int
	Density float64
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

type matrix map[string]neighbourg

func export(m matrix) []Pixel {
	out := make([]Pixel, len(m), len(m))

	keys := []string{}
	indexes := map[string]int{}
	for str := range m {
		indexes[str] = len(keys)
		keys = append(keys, str)
	}

	current := 0.0
	for i, str := range keys {
		ngbr := m[str]
		current += float64(ngbr.C)
		out[i].Color = str
		out[i].Density = current
		out[i].NeighbourgsT = make([]Neighbourg, 0, len(ngbr.T))
		subcurrent := 0.0
		for subStr, count := range ngbr.T {
			_, ok := indexes[subStr]
			if !ok {
				continue
			}
			subcurrent += float64(count)
			out[i].NeighbourgsT = append(out[i].NeighbourgsT, Neighbourg{
				Index:   indexes[subStr],
				Density: subcurrent,
			})
		}
		for j := range out[i].NeighbourgsT {
			out[i].NeighbourgsT[j].Density /= subcurrent
		}

		out[i].NeighbourgsB = make([]Neighbourg, 0, len(ngbr.B))
		subcurrent = 0
		for subStr, count := range ngbr.B {
			_, ok := indexes[subStr]
			if !ok {
				continue
			}
			subcurrent += float64(count)
			out[i].NeighbourgsB = append(out[i].NeighbourgsB, Neighbourg{
				Index:   indexes[subStr],
				Density: subcurrent,
			})
		}
		for j := range out[i].NeighbourgsB {
			out[i].NeighbourgsB[j].Density /= subcurrent
		}

		out[i].NeighbourgsL = make([]Neighbourg, 0, len(ngbr.L))
		subcurrent = 0
		for subStr, count := range ngbr.L {
			_, ok := indexes[subStr]
			if !ok {
				continue
			}
			subcurrent += float64(count)
			out[i].NeighbourgsL = append(out[i].NeighbourgsL, Neighbourg{
				Index:   indexes[subStr],
				Density: subcurrent,
			})
		}
		for j := range out[i].NeighbourgsL {
			out[i].NeighbourgsL[j].Density /= subcurrent
		}

		out[i].NeighbourgsR = make([]Neighbourg, 0, len(ngbr.R))
		subcurrent = 0
		for subStr, count := range ngbr.R {
			_, ok := indexes[subStr]
			if !ok {
				continue
			}
			subcurrent += float64(count)
			out[i].NeighbourgsR = append(out[i].NeighbourgsR, Neighbourg{
				Index:   indexes[subStr],
				Density: subcurrent,
			})
		}
		for j := range out[i].NeighbourgsR {
			out[i].NeighbourgsR[j].Density /= subcurrent
		}
	}
	for i := range out {
		out[i].Density /= current
	}
	return out
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
	out := matrix{}
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

	asJson, err := json.MarshalIndent(export(out), "", "  ")
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
