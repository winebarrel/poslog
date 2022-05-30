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
	file, fingerprint, fillParams := parseArgs()
	defer file.Close()

	proc := func(logBlk *poslog.LogBlock) {
		line, err := jsoniter.ConfigFastest.MarshalToString(logBlk)

		if err != nil {
			panic(err)
		}

		fmt.Println(line)
	}

	p := &poslog.Parser{
		Callback:    proc,
		Fingerprint: fingerprint,
		FillParams:  fillParams,
	}

	err := p.Parse(file)

	if err != nil {
		log.Fatal(err)
	}
}
