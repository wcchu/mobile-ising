package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Global constants
const evolLen = 500  // max num of iteration rounds in one evolution
const mapDim = 20    // map size; num of sites = mapDim^2
const lowTemp = 0.0  // lower limit of temperature scan
const highTemp = 4.0 // upper limit of temperature scan
const nTemps = 40    // num of temperatures
const nRuns = 10     // num of runs at each temperature
const iterMode = 0.0 // interaction mode, 0 : flip, 1 : move, 0-1 : mixed
const maxCPUs = 4    // max cpus to use

// state history
type tempStateHist struct {
	temp float64
	hist [][]State // runs and iterations
}

// macro-state history
type tempMacroHist struct {
	temp     float64
	magHist  []float64
	enerHist []float64
}

type empty struct{}

func main() {
	// set random seed to time
	seed := time.Now().UTC().UnixNano()
	log.Printf("random seed = %v", seed)
	rand.Seed(seed)

	// simulation
	stateRecord, macroRecord := scan()

	// write history to local
	exportStateRecord(stateRecord, Min(10, evolLen))
	exportMacroRecord(macroRecord, Min(1000, evolLen))
}

// scan over temperatures from T0 to T1 with totally n+1 stops including T0 and T1
func scan() ([]tempStateHist, []tempMacroHist) {
	// TODO: remove input args and use global consts directly
	var dT float64
	if nTemps == 0 { // run only with lowTemp, no scan
		dT = 0.0
	} else {
		dT = (highTemp - lowTemp) / float64(nTemps)
	}

	TSHist := make([]tempStateHist, nTemps+1)
	TMHist := make([]tempMacroHist, nTemps+1)

	cpus := Min(runtime.NumCPU(), maxCPUs)
	runtime.GOMAXPROCS(cpus)
	log.Printf("number of cpus = %d", cpus)
	sem := make(chan struct{}, cpus) // limit semaphore to number of cpus
	var wg sync.WaitGroup
	wg.Wait()
	for j := 0; j <= nTemps; j++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(j int) {
			defer wg.Done()
			T := lowTemp + float64(j)*dT
			log.Printf("running T = %f", T)
			TSHist[j].temp, TMHist[j].temp = T, T
			TSHist[j].hist, TMHist[j].magHist, TMHist[j].enerHist = Evolve(T, mapDim, evolLen, nRuns)
			<-sem
		}(j)
	}
	wg.Wait()

	return TSHist, TMHist
}
