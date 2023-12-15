package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

//go:embed arr.tmpl

var enumArrTemplate string

// flags
var (
	typeName   = flag.String("type", "", "name of type")
	exportVar  = flag.Bool("var", false, "export array variable")
	exportFunc = flag.Bool("func", true, "create exported function that returns array")
	output     = flag.String("output", "", "name of generated file, default - T_array.go")
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of enumarr:\n")
	fmt.Fprintf(os.Stderr, "\tenumarr [flags] -type T [files...]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("enumarr: ")
	flag.Usage = Usage
	flag.Parse()

	// validate required 'type' flag
	if len(*typeName) == 0 {
		log.Print("you should specify type")
		flag.Usage()
		os.Exit(1)
	}

	// files to parse
	files := flag.Args()
	if len(files) == 0 {
		// Default: process all files in current directory.
		dir, err := os.ReadDir(".")
		if err != nil {
			log.Print("failed to read files from current directory", err)
			flag.Usage()
			os.Exit(1)
		}
		for _, f := range dir {
			fname := f.Name()
			if strings.HasSuffix(fname, ".go") {
				files = append(files, f.Name())
			}
		}

		if len(files) == 0 {
			log.Print("no files to read!")
			flag.Usage()
			os.Exit(1)
		}
	}

	// output file
	if *output == "" {
		*output = strings.ToLower((*typeName)[:1]) + (*typeName)[1:] + "_array.go"
	}

	g := &enumarr{
		TypeName:   *typeName,
		ExportVar:  *exportVar,
		ExportFunc: *exportFunc,
		Files:      files,
		Output:     *output,
	}

	if err := g.Run(); err != nil {
		log.Fatal("error: ", err)
	}

	log.Println("Array generated!")
}
