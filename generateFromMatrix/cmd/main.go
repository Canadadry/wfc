package main

import (
	"encoding/json"
	"fmt"
	// "image"
	"image/color"
	_ "image/png"
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

	fmt.Println(indexes)

	return nil
}

func pick(p Pickable) int {
	rng := rand.Float64()
	return p.Pick(rng)
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
