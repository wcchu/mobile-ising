package main

import (
	"fmt"
	"math/rand"

	"gopkg.in/cheggaaa/pb.v1"
)

// Location is the location of the site on the lattice map
type Location struct {
	X int
	Y int
}

// SiteInfo describes a site
type SiteInfo struct {
	Loc  Location
	Spin int
}

// State maps site ID to SiteInfo
type State map[int]SiteInfo

// evolve runs the evolution of state in time from a random initial state
// inputs: T = temperature, N = num of sites, L = max time steps
// outputs: state history, magnetization history, energy shift history, time step reaching convergence
func evolve(T float64, D, L int) ([]State, []float64, []float64) {
	// initialization
	stateHist := make([]State, L+1) // initial state + L iterations
	magHist := make([]float64, L+1)
	enerHist := make([]float64, L+1)
	stateHist[0] = initState(D)
	magHist[0] = GetMag(stateHist[0])
	enerHist[0] = 0.0

	// evolve state
	bar := pb.StartNew(L)
	for r := 0; r < L; r++ { // loop over rounds
		bar.Increment()
		state := stateHist[r]
		var mag, enerIter, enerRound float64
		for i := 0; i < D*D; i++ { // i is the id of the site
			state, mag, enerIter = Iterate(D, state, i, T)
			enerRound = enerRound + enerIter
		}
		stateHist[r+1], magHist[r+1], enerHist[r+1] = state, mag, enerRound
	}
	bar.FinishPrint(fmt.Sprintf("evolution ends for T = %f", T))

	// return only the evolving part of history
	return stateHist, magHist, enerHist
}

// initState creates an initial state with random spins
// D = dimension of map
func initState(D int) State {
	state := make(State)
	var i int
	for ix := 0; ix < D; ix++ {
		for iy := 0; iy < D; iy++ {
			state[i] = SiteInfo{
				Loc:  Location{X: ix, Y: iy},
				Spin: rand.Intn(2)*2 - 1,
			}
			i = i + 1
		}
	}
	return state
}
