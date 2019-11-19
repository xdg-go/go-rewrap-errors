package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/xdg-go/go-rewrap-errors/internal/rewriter"
)

func printHelp(exitCode int) {
	fmt.Fprint(os.Stderr, "usage: go-rewrap-errors [options] [input-filename]\n\n")
	fmt.Fprint(os.Stderr, "If no input filename is provided, it will read from stdin.\n\n")
	pflag.PrintDefaults()
	os.Exit(exitCode)
}

func main() {
	// Setup options
	optWrite := pflag.BoolP("write", "w", false, "overwrite source file instead of writing to stdout")
	optHelp := pflag.BoolP("help", "h", false, "show this help text")
	pflag.Parse()
	if *optHelp {
		printHelp(0)
	}

	// Read original source
	var oldSource []byte
	var filename string
	var err error
	switch len(pflag.Args()) {
	case 0:
		filename = "stdin"
		oldSource, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("couldn't read from stdin: %v", err)
			os.Exit(1)
		}
	case 1:
		filename := os.Args[1]
		oldSource, err = ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("couldn't read from %s: %v", filename, err)
			os.Exit(1)
		}
	default:
		log.Print("Error: too many command line arguments\n\n")
		printHelp(1)
	}

	// Rewrite the original source
	newSource, err := rewriter.Rewrite(filename, oldSource)
	if err != nil {
		log.Fatal(err)
	}

	// Overwrite or print the new source
	if *optWrite {
		fi, err := os.Stat(filename)
		if err != nil {
			log.Fatal(err)
		}
		ioutil.WriteFile(filename, newSource, fi.Mode())
	} else {
		fmt.Print(string(newSource))
	}
}
