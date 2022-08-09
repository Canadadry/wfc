package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"strings"
)

func main() {
	if err := run(os.Args[1], 32, 24); err != nil {
		fmt.Println("failed", err)
	}
}

type Pickable interface {
	Pick(p float64) int
}

type Pixel struct {
	Color        string
	color        color.RGBA
	Density      float64
	NeighbourgsT Neighbourgs
	NeighbourgsB Neighbourgs
	NeighbourgsL Neighbourgs
	NeighbourgsR Neighbourgs
}

type Pixels []Pixel

func (ps Pixels) Pick(probability float64) int {
	for i, pix := range ps {
		if pix.Density > probability {
			return i
		}
	}
	return len(ps) - 1
}

type Neighbourg struct {
	Index   int
	Density float64
}

type Neighbourgs []Neighbourg

func (ps Neighbourgs) Pick(probability float64) int {
	for i, pix := range ps {
		if pix.Density > probability {
			return i
		}
	}
	return len(ps) - 1
}

func run(filename string, w, h int) error {
	infile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer infile.Close()

	pixels := Pixels{}
	err = json.NewDecoder(infile).Decode(&pixels)
	if err != nil {
		return err
	}

	indexes := make([]int, w*h)

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			if i > 0 {
				leftIndex := indexes[i-1+w*j]
				indexes[i+w*j] = pick(pixels[leftIndex].NeighbourgsR)
			} else {
				if j > 0 {
					leftIndex := indexes[w*(j-1)]
					indexes[i+w*j] = pick(pixels[leftIndex].NeighbourgsR)
				} else {
					indexes[i+w*j] = pick(pixels)
				}
			}
		}
	}

	for i := range pixels {
		c, err := colorFromString(pixels[i].Color)
		if err != nil {
			return fmt.Errorf("cannot decode color %d : %w", i, err)
		}
		pixels[i].color = c
	}

	return pixels.toPng(indexes, w, h)

}

func pick(p Pickable) int {
	rng := rand.Float64()
	return p.Pick(rng)
}

func (ps Pixels) toPng(indexes []int, w, h int) error {
	outfile, err := os.Create("out.png")
	if err != nil {
		return err
	}

	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			out.Set(x, y, ps[indexes[x+w*y]].color)
		}
	}
	defer outfile.Close()
	return png.Encode(outfile, out)
}

func colorFromString(str string) (color.RGBA, error) {
	c := [3]int{}
	err := json.NewDecoder(strings.NewReader(str)).Decode(&c)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{
		R: uint8(c[0]),
		G: uint8(c[1]),
		B: uint8(c[2]),
		A: 255,
	}, nil
}
