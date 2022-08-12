package generateFromMatrix

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math/rand"
	"os"
	"strings"
)

type Pixels map[string]color.RGBA

func (ps Pixels) load(indexes []string) error {
	for i := range indexes {
		c, err := colorFromString(indexes[i].Color)
		if err != nil {
			return fmt.Errorf("cannot decode color %d : %w", i, err)
		}
		ps[i].color = c
	}
	return nil
}

func (ps Pixels) toPng(out io.Writer, indexes []int, w, h int) error {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, ps[img[x+w*y]].color)
		}
	}
	return png.Encode(out, img)
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
