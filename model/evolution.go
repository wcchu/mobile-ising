package main

import (
	"fmt"
	"math"
	"math/rand"

	"gopkg.in/cheggaaa/pb.v1"
)

// Location of a site
type Location struct {
	X float64
	Y float64
}

// State stores the state of the system
type State struct {
	Locations   []Location
	Connections []int
	Spins       []int
}

// evolve runs the evolution of state in time from a random initial state
// inputs: T = temperature, N = num of sites, L = max time steps
// outputs: state history, magnetization history, time step reaching convergence
func evolve(T float64, N, L int, cmean float64, cmax int) ([]State, []float64) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax, forceConns)

	// initialization
	stateHist := make([]State, L+1) // initial state + L iterations
	magHist := make([]float64, L+1)
	stateHist[0] = initState(N, ps)
	magHist[0] = GetMag(stateHist[0].Spins)

	// evolve state
	end := 0
	bar := pb.StartNew(L)
	for i := 0; i < L; i++ {
		bar.Increment()
		if math.Abs(magHist[i]) < 1.0 { // not converged yet, iterate
			stateHist[i+1], magHist[i+1] = Iterate(stateHist[i], T)
			end = i + 1 // the time index we want to record up to
		} else { // converged
			magHist[i+1] = magHist[i]
		}
	}
	bar.FinishPrint(fmt.Sprintf("evolution ends at time step %d for T = %f", end, T))

	// return only the evolving part of history
	return stateHist[0 : end+1], magHist[0 : end+1]
}

// initState creates a random initial state
// N = num of sites
// ps = prob distribution of having k connections
func initState(N int, ps []float64) State {
	locs := make([]Location, N)
	conns := make([]int, N)
	sps := make([]int, N)

	for i := 0; i < N; i++ {
		locs[i] = Location{X: rand.Float64(), Y: rand.Float64()}
		c := AssignConn(ps, rand.Float64())
		if c < 0 {
			return State{}
		}
		conns[i] = c
		sps[i] = rand.Intn(2)*2 - 1
	}

	return State{
		Locations:   locs,
		Connections: conns,
		Spins:       sps,
	}
}
