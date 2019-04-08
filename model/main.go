package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Global constants
const evolLen = 500 // max num of iteration rounds in one evolution
const mapDim = 20   // map size; num of sites = mapDim^2
const lowTemp = 0.0
const highTemp = 4.0
const nTemps = 40
const nRuns = 10
const iterMode = 0.0 // 0 : flip, 1 : move, 0-1 : mixed
const maxCPUs = 4

type tempStateHist struct {
	temp float64
	hist [][]State // runs and iterations
}

type tempMacroHist struct {
	temp     float64
	magHist  []float64
	enerHist []float64
}

type empty struct{}

// main
func main() {
	// set random seed to time
	seed := time.Now().UTC().UnixNano()
	log.Printf("random seed = %v", seed)
	rand.Seed(seed)

	stateRecord, macroRecord := scan(lowTemp, highTemp, nTemps)

	// write history to local
	exportStateRecord(stateRecord, Min(10, evolLen))
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

	cpus := Min(runtime.NumCPU(), maxCPUs)
	runtime.GOMAXPROCS(cpus)
	log.Printf("number of cpus = %d", cpus)
	sem := make(chan struct{}, cpus) // limit semaphore to number of available cpus
	var wg sync.WaitGroup
	wg.Wait()
	for j := 0; j <= n; j++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(j int) {
			defer wg.Done()
			T := T0 + float64(j)*dT
			log.Printf("running T = %f", T)
			TSHist[j].temp, TMHist[j].temp = T, T
			TSHist[j].hist, TMHist[j].magHist, TMHist[j].enerHist = Evolve(T, mapDim, evolLen, nRuns)
			<-sem
		}(j)
	}
	wg.Wait()

	return TSHist, TMHist
}
