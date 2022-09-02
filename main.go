package main

import (
	"io"
	"os"

	parser "github.com/spectrocloud/kubesplit/internal/parser"
)

func main() {
	if len(os.Args) != 2 {
		panic("Requires one argument, the output folder")
	}

	dir := os.Args[1]

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	reader := io.Reader(os.Stdin)

	if err := parser.Generate(reader, dir); err != nil {
		panic(err)
	}
}
