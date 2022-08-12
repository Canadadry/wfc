package generateFromMatrix

import (
	"fmt"
	"io"
)

type Constraint struct {
	PatternSize int
	Index       []string
	Patterns    map[string]int
}

func (c Constraint) Patterns() Patterns {
	p := Patterns{}
	p.load(c.Patterns)
	return p
}

func (c Constraint) Pixels() (Pixels, error) {
	ps := Pixels{}
	err := ps.load(c.Index)
	return ps, err
}

func load(in io.Reader) (Constraint, error) {
	c := Constraint{}
	err = json.NewDecoder(in).Decode(&c)
	return c, err
}

func run(in io.Reader, out io.Writer, patternSize, w, h int) error {

	c, err := load(in)
	if err != nil {
		return fmt.Errorf("while loading constraint file :%w", err)
	}

	indexes := make([]int, w*h)
	written := make([]bool, w*h)

	for i := 0; i <= (w - patternSize); i++ {
		for j := 0; j <= (h - patternSize); j++ {
			ok := false
			patterns := c.Patterns()
			for !ok {
				p, err = patterns.Pick()
				if err != nil {
					return err
				}
				ok := apply(written, indexes, p, x, y, patternSize)
			}
		}
	}

	ps, err := c.Pixels()
	if err != nil {
		fmt.Errorf("while decoding color indexes : %w", err)
	}
	return ps.toPng(out, indexes, w, h)
}

func apply(written []bool, indexes, patterns []int, x, y, patternSize int) bool {

}
