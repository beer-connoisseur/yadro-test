package main

import (
	"log"
)

func main() {
	const (
		config         = "sunny_5_skiers/config.json"
		inputFilename  = "sunny_5_skiers/events"
		outputFilename = "output.txt"
	)

	cfg, err := New(config)
	if err != nil {
		log.Fatal(err)
	}

	err = parseEvents(inputFilename, outputFilename, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
