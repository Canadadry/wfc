package map2matrix

import (
	"flag"
	"fmt"
	"os"
)

func Run(name string, args []string) error {
	imageFilename := "in.png"
	outFilename := "out.json"
	patternSize := 1

	f := flag.NewFlagSet(name, flag.ContinueOnError)

	f.StringVar(&imageFilename, "in", imageFilename, "image to analyze")
	f.StringVar(&outFilename, "out", outFilename, "export filename")
	f.IntVar(&patternSize, "pattern-size", patternSize, "size of pattern to extract")
	f.IntVar(&patternSize, "size", patternSize, "size of pattern to extract")

	err := f.Parse(args)
	if err != nil {
		return err
	}

	infile, err := os.Open(imageFilename)
	if err != nil {
		return err
	}
	defer infile.Close()

	outfile, err := os.Create(outFilename)
	if err != nil {
		return err
	}
	defer infile.Close()

	err = Process(infile, outfile, patternSize)
	if err != nil {
		return fmt.Errorf("processing image : %w", err)
	}

	return nil
}
