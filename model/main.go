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

	tempStateHist, tempMagHist := scan(0.0, 1.0, 2)

	// write history to local
	//exportStateHist(stateHist, 10)
	//exportMagHist(magHist, numSites, 100)
}

// scan over temperatures from T0 to T1 with totally n stops including T0 and T1
func scan(T0, T1 float64, n int) ([][]State, [][]float64) {
	var dT float64
	if n < 1 {
		panic("not enough scan slices")
	} else if n == 1 { // run only with T1, no scan
		dT = T1 - T0
		T1 = T0
	} else {
		dT = (T1 - T0) / float64(n-1)
	}

	TSHist := make([][]State, n)
	TMHist := make([][]float64, n)

	T := T0
	for j := 0; j < n; j++ {
		TSHist[j], TMHist[j] = evolve(T, numSites, lenEvol, meanConns, maxConns)
		T = T + dT
	}

	return TSHist, TMHist
}
