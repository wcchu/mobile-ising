package main

import (
	"log"
	"math/rand"
	"time"
)

// Global constants
const lenEvol = 100000
const numSites = 400
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
	stateHist, magHist, L := evolve(temperature, numSites, lenEvol, meanConns, maxConns)
	// only export history up to step L; the rest are either converged or out of range
	exportStateHist(stateHist[0:L+1], 100)
	exportMagHist(magHist[0:L+1], numSites, 100)
}
