package generateFromMatrix

import (
	"encoding/json"
	"fmt"
	"io"
)

type Constraint struct {
	PatternSize int
	Index       []string
	Patterns    map[string]int
}

func (c Constraint) GetPatterns() Patterns {
	p := Patterns{}
	p.load(c.Patterns)
	return p
}

func (c Constraint) GetPixels() (Pixels, error) {
	ps := Pixels{}
	err := ps.load(c.Index)
	return ps, err
}

func load(in io.Reader) (Constraint, error) {
	c := Constraint{}
	err := json.NewDecoder(in).Decode(&c)
	return c, err
}

func Process(in io.Reader, out io.Writer, patternSize, w, h int, rand func() float64) error {

	c, err := load(in)
	if err != nil {
		return fmt.Errorf("while loading constraint file :%w", err)
	}

	indexes, err := generate(c, patternSize, w, h, rand)
	if err != nil {
		fmt.Errorf("while generating image : %w", err)
	}

	ps, err := c.GetPixels()
	if err != nil {
		fmt.Errorf("while decoding color indexes : %w", err)
	}

	return ps.toPng(out, indexes, w, h)
}

func generate(c Constraint, patternSize, w, h int, rand func() float64) ([]int, error) {
	indexes := make([]int, w*h)
	written := make([]bool, w*h)

	for i := 0; i <= (w - patternSize); i++ {
		for j := 0; j <= (h - patternSize); j++ {
			ok := false
			patterns := c.GetPatterns()
			for !ok {
				p, err := patterns.Pick(rand)
				if err != nil {
					return nil, err
				}
				ok = apply(written, indexes, p, i, j, w, patternSize)
			}
		}
	}
	return indexes, nil
}

func apply(written []bool, indexes, patterns []int, x, y, w, patternSize int) bool {
	if !canBePlaced(written, indexes, patterns, x, y, w, patternSize) {
		return false
	}
	update(written, indexes, patterns, x, y, w, patternSize)
	return true
}

func canBePlaced(written []bool, indexes, patterns []int, x, y, w, patternSize int) bool {
	for i := 0; i < patternSize; i++ {
		for j := 0; j < patternSize; j++ {
			if written[x+i+(j+i)*w] {
				if indexes[x+i+(j+i)*w] != patterns[i+j*patternSize] {
					return false
				}
			}
		}
	}
	return true
}

func update(written []bool, indexes, patterns []int, x, y, w, patternSize int) {
	for i := 0; i < patternSize; i++ {
		for j := 0; j < patternSize; j++ {
			written[x+i+(j+i)*w] = true
			indexes[x+i+(j+i)*w] = patterns[i+j*patternSize]
		}
	}
}
