package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "embed"

	helm "github.com/spectrocloud/kubesplit/internal/helm"

	parser "github.com/spectrocloud/kubesplit/internal/parser"
)

//go:embed _helper.tpl
var template string

//go:embed values.yaml
var values string

var helmF = flag.Bool("helm", false, "do a minimum helm conversion")

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		panic("Requires one argument, the output folder")
	}

	dir := flag.Arg(0)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	reader := io.Reader(os.Stdin)

	opts := []parser.Option{}
	if *helmF {
		opts = append(opts, parser.WithMutators(helm.Mutator))
	}

	if err := parser.Generate(reader, dir, opts...); err != nil {
		panic(err)
	}

	if *helmF {
		if err := ioutil.WriteFile(filepath.Join(dir, "_helpers.tpl"), []byte(template), os.ModePerm); err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(filepath.Join(dir, "values.yaml.default"), []byte(template), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
