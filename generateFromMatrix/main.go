package generateFromMatrix

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func Run(name string, args []string) error {
	constrainFile := "in.json"
	outFilename := "out.png"
	patternSize := 1
	w := 32
	h := 24
	seed := 0

	f := flag.NewFlagSet(name, flag.ContinueOnError)

	f.StringVar(&constrainFile, "in", constrainFile, "constraint to use for image generation")
	f.StringVar(&outFilename, "out", outFilename, "out image")
	f.IntVar(&patternSize, "pattern-size", patternSize, "size of pattern to extract")
	f.IntVar(&patternSize, "size", patternSize, "size of pattern to extract")
	f.IntVar(&w, "w", w, "width of generated image")
	f.IntVar(&h, "h", h, "heght of generated image")
	f.IntVar(&seed, "seed", seed, "seed for random generator, if 0 use number of second since 1 jan 1970")

	err := f.Parse(args)
	if err != nil {
		return err
	}

	infile, err := os.Open(constrainFile)
	if err != nil {
		return err
	}
	defer infile.Close()

	outfile, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	defer infile.Close()

	if seed == 0 {
		seed = int(time.Now().Unix())
	}

	err = Process(infile, outfile, patternSize, w, h, seed)
	if err != nil {
		return fmt.Errorf("processing image : %w", err)
	}

	return nil
}
