package main

import (
	"math/rand"
	"time"
)

// Global constants
const maxEvoTime = 1
const numSites = 5
const meanConns = 1
const maxConns = 2
const forceConns = true

// main
func main() {
	// set random seed to time
	rand.Seed(time.Now().UTC().UnixNano())

	temperature := 1.0 // TODO: loop over temperatures
	evolve(temperature, numSites, maxEvoTime, meanConns, maxConns)
}
