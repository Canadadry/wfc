package map2matrix

import (
	"encoding/json"
	"image/color"
)

type PixelGetter interface {
	At(x, y int) color.Color
}

func extractPatternAt(src PixelGetter, x, y, patternSize int, ci ColorIndex) Pattern {
	pixels := make(pattern, patternSize*patternSize)
	for i := 0; i < patternSize; i++ {
		for j := 0; j < patternSize; j++ {
			pixels[i+j*patternSize], _ = ci.Get(src.At(x+i, y+j))
		}
	}
	return pixels
}

type pattern []int

func (p pattern) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

//
// func (p pattern) Rotate() Pattern {
// 	return p
// }
//
// func (p pattern) Mirror() Pattern {
// 	return p
// }
