package main

import (
	"math/rand"
	"time"
)

// Global constants
const maxEvoTime = 10000
const numSites = 100
const meanConns = 4
const maxConns = 10

// main
func main() {
	// set random seed to time
	rand.Seed(time.Now().UTC().UnixNano())

	temperature := 1.0 // TODO: loop over temperatures
	evolve(temperature, numSites, maxEvoTime, meanConns, maxConns)
}
