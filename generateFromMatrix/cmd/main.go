package main

import (
	"app/generateFromMatrix"
	"flag"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Println("failed", err)
	}
}

func run(name string, args []string) error {
	constrainFile := "in.json"
	outFilename := "out.pg"
	patternSize := 1
	w := 32
	h := 24

	f := flag.NewFlagSet(name, flag.ContinueOnError)

	f.StringVar(&constrainFile, "in", constrainFile, "constraint to use for image generation")
	f.StringVar(&outFilename, "out", outFilename, "out image")
	f.IntVar(&patternSize, "pattern-size", patternSize, "size of pattern to extract")
	f.IntVar(&patternSize, "size", patternSize, "size of pattern to extract")
	f.IntVar(&w, "w", w, "width of generated image")
	f.IntVar(&h, "h", h, "width of generated image")

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

	err = generateFromMatrix.Process(infile, outfile, patternSize)
	if err != nil {
		return fmt.Errorf("processing image : %w", err)
	}

	return nil
}
