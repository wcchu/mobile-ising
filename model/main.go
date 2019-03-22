package main

import (
	"log"
	"math/rand"
	"time"
)

// Global constants
const lenEvol = 10000
const numSites = 400
const meanConns = 4
const maxConns = 8
const forceConns = true
const lowTemp = 0.0
const highTemp = 5.0
const nTemps = 11
const iterMode = 0.5 // 0 : flip, 1 : move, 0-1 : mixed

type tempStateHist struct {
	temp float64
	hist []State
}

type tempMacroHist struct {
	temp     float64
	magHist  []float64
	enerHist []float64
}

// main
func main() {
	// set random seed to time
	seed := time.Now().UTC().UnixNano()
	log.Printf("random seed = %v", seed)
	rand.Seed(seed)

	stateRecord, macroRecord := scan(lowTemp, highTemp, nTemps)

	// write history to local
	exportStateRecord(stateRecord, 10)
	exportMacroRecord(macroRecord, numSites, 500)
}

// scan over temperatures from T0 to T1 with totally n stops including T0 and T1
func scan(T0, T1 float64, n int) ([]tempStateHist, []tempMacroHist) {
	var dT float64
	if n < 1 {
		panic("not enough scan slices")
	} else if n == 1 { // run only with T1, no scan
		dT = T1 - T0
		T1 = T0
	} else {
		dT = (T1 - T0) / float64(n-1)
	}

	TSHist := make([]tempStateHist, n)
	TMHist := make([]tempMacroHist, n)

	T := T0
	for j := 0; j < n; j++ {
		TSHist[j].temp = T
		TMHist[j].temp = T
		log.Printf("running evolution for temperature at %f", T)
		TSHist[j].hist, TMHist[j].magHist, TMHist[j].enerHist = evolve(T, numSites, lenEvol, meanConns, maxConns)
		T = T + dT
	}

	return TSHist, TMHist
}
