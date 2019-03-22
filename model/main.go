package main

import (
	"log"
	"math/rand"
	"time"
)

// Global constants
const lenEvol = 500000
const numSites = 100
const meanConns = 4
const maxConns = 8
const forceConns = true

// main
func main() {
	// set random seed to time
	seed := time.Now().UTC().UnixNano()
	log.Printf("random seed = %v", seed)
	rand.Seed(seed)

	temperature := 1.0 // TODO: loop over temperatures
	stateHist, magHist := evolve(temperature, numSites, lenEvol, meanConns, maxConns)

	// write history to local
	exportStateHist(stateHist, 10)
	exportMagHist(magHist, numSites, 100)
}
