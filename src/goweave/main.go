package main

import (
	"log"

	"./weave"
)

const (
	version = "v0.1"
)

func Print() {

}

// main is the main point of entry for running goweave
func main() {
	log.Println("goweave " + version)

	w := weave.NewWeave()
	w.Run()

}
