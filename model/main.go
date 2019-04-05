package main

import (
	"log"
	"math/rand"
	"time"
)

// Global constants
const evolLen = 100 // max num of iteration rounds in one evolution
const mapDim = 20   // total num of sites = mapDim^2
const forceConns = true
const lowTemp = 0.0
const highTemp = 4.0
const nTemps = 20
const iterMode = 0.0 // 0 : flip, 1 : move, 0-1 : mixed
const therRounds = 0 // define thermalization with the last numSite * therRounds iterations

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
	exportMacroRecord(macroRecord, Min(1000, evolLen))
}

// scan over temperatures from T0 to T1 with totally n+1 stops including T0 and T1
func scan(T0, T1 float64, n int) ([]tempStateHist, []tempMacroHist) {
	var dT float64
	if n == 0 { // run only with T1, no scan
		dT = T1 - T0
		T1 = T0
	} else {
		dT = (T1 - T0) / float64(n)
	}

	TSHist := make([]tempStateHist, n+1)
	TMHist := make([]tempMacroHist, n+1)

	T := T0
	for j := 0; j <= n; j++ {
		TSHist[j].temp = T
		TMHist[j].temp = T
		log.Printf("running evolution for temperature at %f", T)
		TSHist[j].hist, TMHist[j].magHist, TMHist[j].enerHist = evolve(T, mapDim, evolLen)
		T = T + dT
	}

	return TSHist, TMHist
}
