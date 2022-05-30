package main

import (
	"fmt"
	"io"
	"log"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/winebarrel/poslog"
)

var version string

func init() {
	log.SetFlags(0)
}

func main() {
	file := parseArgs()
	defer file.Close()

	err := poslog.Parse(file, func(block *poslog.Block) {
		line, err := jsoniter.ConfigFastest.MarshalToString(block)

		if err != nil {
			panic(err)
		}

		fmt.Println(line)
	})

	if err != nil {
		log.Fatal(err)
	}
}

func parseArgs() io.ReadCloser {
	if len(os.Args) > 2 {
		log.Fatalf("usage: %s [-version] POSTGRESQL_LOG", os.Args[0])
	}

	if len(os.Args) == 1 {
		return os.Stdin
	}

	if os.Args[1] == "-version" {
		fmt.Fprintln(os.Stderr, version)
		os.Exit(0)
	}

	file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0)

	if err != nil {
		log.Fatal(err)
	}

	return file
}
