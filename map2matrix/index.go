package map2matrix

import (
	"fmt"
	"image/color"
)

type ColorIndex map[string]int

func (ci ColorIndex) Add(c color.Color) {
	if _, ok := ci[colorToString(c)]; ok {
		return
	}
	ci[colorToString(c)] = len(ci)
}

func (ci ColorIndex) Get(c color.Color) (int, bool) {
	i, ok := ci[colorToString(c)]
	return i, ok
}

func (ci ColorIndex) ToSlice() []string {
	out := make([]string, len(ci))
	for k, i := range ci {
		out[i] = k
	}
	return out
}

func colorToString(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("[%d,%d,%d]", r/257, g/257, b/257)
}
