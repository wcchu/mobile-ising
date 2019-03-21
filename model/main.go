package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Global constants
const lenEvol = 50000
const numSites = 400
const meanConns = 4
const maxConns = 8
const forceConns = true

// main
func main() {
	// set random seed to time
	seed := time.Now().UTC().UnixNano()
	fmt.Printf("random seed = %v", seed)
	rand.Seed(seed)

	temperature := 1.0 // TODO: loop over temperatures
	_, hist := evolve(temperature, numSites, lenEvol, meanConns, maxConns)
	exportMagHist(hist)
}
