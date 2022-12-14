package map2matrix

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
)

type Pattern interface {
	String() string
	// Rotate() Pattern
	// Mirror() Pattern
}

type Result struct {
	PatternSize int
	Index       []string
	Patterns    map[string]int
}

func Process(in io.Reader, out io.Writer, patternSize int) error {

	src, format, err := image.Decode(in)
	if err != nil {
		return fmt.Errorf("while decoding source img %s : %w", format, err)
	}

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	if patternSize > w || patternSize > h {
		return fmt.Errorf("pattern size bigger than image (%d,%d) > (%d,%d)", patternSize, patternSize, w, h)
	}

	result := extract(src, w, h, patternSize)

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	err = enc.Encode(result)
	if err != nil {
		return fmt.Errorf("while json encoding result : %w", err)
	}
	return nil
}

type PixelGetter interface {
	At(x, y int) color.Color
}

func extract(src PixelGetter, w, h, patternSize int) Result {
	result := Result{
		PatternSize: patternSize,
		Patterns:    map[string]int{},
	}
	colorIndex := ColorIndex{}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			colorIndex.Add(src.At(x, y))
		}
	}
	for x := 0; x <= (w - patternSize); x++ {
		for y := 0; y <= (h - patternSize); y++ {
			p := extractPatternAt(src, x, y, patternSize, colorIndex)
			result.Patterns[p.String()] = result.Patterns[p.String()] + 1
		}
	}

	result.Index = colorIndex.ToSlice()
	return result
}
