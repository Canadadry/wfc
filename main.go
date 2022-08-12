package main

import (
	"app/generateFromMatrix"
	"app/map2matrix"
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println("failed :", err)
	}
}

const errMsg = "usage %s command [args ...]\ncommand: \n - generate\n - build\n"

func run(args []string) error {
	if len(args) <= 1 {
		return fmt.Errorf(errMsg, args[0])
	}
	switch args[1] {
	case "generate":
		return generateFromMatrix.Run(args[1], args[2:])
	case "build":
		return map2matrix.Run(args[1], args[2:])
	}
	return fmt.Errorf(errMsg, args[0])
}
