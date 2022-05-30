package main

import (
	"fmt"
	"log"

	jsoniter "github.com/json-iterator/go"
	"github.com/winebarrel/poslog"
)

func init() {
	log.SetFlags(0)
}

func main() {
	file, fingerprint := parseArgs()
	defer file.Close()

	err := poslog.Parse(file, fingerprint, func(block *poslog.Block) {
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
