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

type Patterns struct {
	Patterns []string
	Count    []int
	Total    int
}

func (p *Patterns) load(patterns map[string]int) {
	if p.Patterns == nil {
		p.Patterns = make([]string, len(patterns))
		p.Density = make([]float64, len(patterns))
	}

	i := 0
	for k, v := range patterns {
		p.Total = p.Total + v
		p.Patterns[i] = k
		p.Count[i] = v
		i++
	}
}

func (p *Patterns) remove(i int) error {
	if i < 0 || i >= len(p.Patterns) {
		return fmt.Errorf("remove %d out of bound [0:%d]", i, len(p.Patterns))
	}
	p.Total = p.Total - p.Count[i]
	p.Patterns[i] = p.Patterns[len(p.Patterns)-1]
	p.Count[i] = p.Count[len(p.Count)-1]

	p.Patterns = p.Patterns[0 : len(p.Patterns)-1]
	p.Count = p.Count[0 : len(p.Count)-1]
}

func (p *Patterns) Pick() (string, error) {
	rng := rand.Float64()
	i, err := p.pick(rng)
	if err != nil {
		return "", err
	}
	return p.Patterns[i], p.remove(i)
}

func (p Patterns) pick(rng float64) (int, error) {
	if len(p.Count) == 0 {
		return 0, fmt.Errorf("nothing left to pick")
	}
	current := 0
	for i, c := range p.Count {
		current = current + c
		if current > rng*p.Total {
			return i, nil
		}
	}
	return 0, fmt.Errorf("pick with %d(%v) out of bound [0:%d]", p.Count*rng, rng, p.Count)
}
