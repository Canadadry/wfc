package generateFromMatrix

import (
	"encoding/json"
	"fmt"
)

type Patterns struct {
	Patterns []string
	Count    []int
	Total    int
}

func (p *Patterns) load(patterns map[string]int) {
	if p.Patterns == nil {
		p.Patterns = make([]string, len(patterns))
		p.Count = make([]int, len(patterns))
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
	return nil
}

func (p *Patterns) Pick(rand func() float64) ([]int, error) {
	rng := rand()
	i, err := p.pick(rng)
	if err != nil {
		return nil, err
	}
	ret, err := explode(p.Patterns[i])
	if err != nil {
		return ret, err
	}
	return ret, p.remove(i)
}

func explode(p string) ([]int, error) {
	out := []int{}
	err := json.Unmarshal([]byte(p), &out)
	return out, err
}

func (p Patterns) pick(rng float64) (int, error) {
	// fmt.Println("pick")
	rng = rng * float64(p.Total)
	if len(p.Count) == 0 {
		return 0, fmt.Errorf("nothing left to pick")
	}
	current := 0.0
	for i, c := range p.Count {
		// fmt.Println(current, p.Total)
		current = current + float64(c)
		if current >= rng {
			return i, nil
		}
	}
	return 0, fmt.Errorf("pick with %v out of bound [0:%d]", rng, p.Total)
}
