package main

import (
	"math/rand"
	"time"
)

// Global constants
const lenEvol = 100
const numSites = 10
const meanConns = 4
const maxConns = 8
const forceConns = true

// main
func main() {
	// set random seed to time
	rand.Seed(time.Now().UTC().UnixNano())

	temperature := 1.0 // TODO: loop over temperatures
	_, hist := evolve(temperature, numSites, lenEvol, meanConns, maxConns)
	exportMagHist(hist)
}
