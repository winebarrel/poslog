package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var version string

func parseArgs() (io.ReadCloser, bool) {
	fingerprint := flag.Bool("fingerprint", false, "Add SQL fingerprint")
	showVersion := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	args := flag.Args()

	if len(args) >= 2 {
		printUsageAndExit()
	}

	if *showVersion {
		printVersionAndEixt()
	}

	file := os.Stdin

	if len(args) == 1 {
		var err error

		file, err = os.OpenFile(args[0], os.O_RDONLY, 0)

		if err != nil {
			log.Fatal(err)
		}
	}

	return file, *fingerprint
}

func printUsageAndExit() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func printVersionAndEixt() {
	fmt.Fprintln(os.Stderr, version)
	os.Exit(0)
}
