package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/brahama/cf-doc/doc"
	"github.com/brahama/cf-doc/print"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/tj/docopt"
)

var version = "v0.2.0"

const usage = `
  Usage:
    cf-docs [--no-required] [json | md | markdown] <path>...
    cf-docs -h | --help

  Examples:

    # View inputs and outputs
    $ cf-docs ./my-module

    # View inputs and outputs for variables.tf and outputs.tf only
    $ cf-docs variables.tf outputs.tf

    # Generate a JSON of inputs and outputs
    $ cf-docs json ./my-module

    # Generate markdown tables of inputs and outputs
    $ cf-docs md ./my-module

    # Generate markdown tables of inputs and outputs, but don't print "Required" column
    $ cf-docs --no-required md ./my-module

    # Generate markdown tables of inputs and outputs for the given module and ../config.tf
    $ cf-docs md ./my-module ../config.tf

  Options:
    -h, --help     show help information

`

func main() {
	args, err := docopt.Parse(usage, nil, true, version, true)
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	paths := args["<path>"].([]string)
	for _, p := range paths {
		pi, err := os.Stat(p)
		if err != nil {
			log.Fatal(err)
		}

		if !pi.IsDir() {
			names = append(names, p)
			continue
		}

		files, err := filepath.Glob(fmt.Sprintf("%s/*.tf", p))
		if err != nil {
			log.Fatal(err)
		}

		names = append(names, files...)
	}

	files := make(map[string]*ast.File, len(names))

	for _, name := range names {
		buf, err := ioutil.ReadFile(name)
		if err != nil {
			log.Fatal(err)
		}

		f, err := hcl.ParseBytes(buf)
		if err != nil {
			log.Fatal(err)
		}

		files[name] = f
	}

	doc := doc.Create(files)
	printRequired := !args["--no-required"].(bool)

	var out string

	switch {
	case args["markdown"].(bool):
		out, err = print.Markdown(doc, printRequired)
	case args["md"].(bool):
		out, err = print.Markdown(doc, printRequired)
	case args["json"].(bool):
		out, err = print.JSON(doc)
	default:
		out, err = print.Pretty(doc)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
}