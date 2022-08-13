package generateFromMatrix

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func Run(name string, args []string) error {
	constrainFile := "in.json"
	outFilename := "out.png"
	w := 32
	h := 24
	var seed int64

	f := flag.NewFlagSet(name, flag.ContinueOnError)

	f.StringVar(&constrainFile, "in", constrainFile, "constraint to use for image generation")
	f.StringVar(&outFilename, "out", outFilename, "out image")
	f.IntVar(&w, "w", w, "width of generated image")
	f.IntVar(&h, "h", h, "heght of generated image")
	f.Int64Var(&seed, "seed", seed, "seed for random generator, if 0 use number of second since 1 jan 1970")

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
		seed = time.Now().Unix()
		rand.Seed(seed)
	}

	err = Process(infile, outfile, w, h, rand.Float64)
	if err != nil {
		return fmt.Errorf("processing image : %w", err)
	}

	return nil
}
