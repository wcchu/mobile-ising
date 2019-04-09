package main

import (
	"math/rand"
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

// Evolve runs the evolution of state in time from a random initial state Q times,
// and
// inputs: T = temperature, N = num of sites, L = max time steps
// outputs: state history, magnetization history, energy shift history, time step reaching convergence
func Evolve(T float64, D, L, Q int) ([][]State, []float64, []float64) {
	// initialization
	stateHist := make([][]State, Q)
	magHist := make([][]float64, Q)
	enerHist := make([][]float64, Q)
	for run := 0; run < Q; run++ {
		stateHist[run] = make([]State, L+1)
		magHist[run] = make([]float64, L+1)
		enerHist[run] = make([]float64, L+1)
	}

	for run := 0; run < Q; run++ {
		stateHist[run][0] = initState(D)
		magHist[run][0] = GetMag(stateHist[run][0])
		enerHist[run][0] = 0.0

		// evolve state
		for r := 0; r < L; r++ { // loop over rounds
			state := stateHist[run][r]
			var enerIter, enerRound float64
			for i := 0; i < D*D; i++ { // i is the id of the site
				state, enerIter = Iterate(D, state, i, T)
				enerRound = enerRound + enerIter
			}
			stateHist[run][r+1], enerHist[run][r+1] = state, enerRound
			magHist[run][r+1] = GetMag(stateHist[run][r+1])
		}
	}

	return stateHist, AveRuns(magHist), AveRuns(enerHist)
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

// AveRuns averages the values across runs
func AveRuns(h [][]float64) []float64 {
	sumHist := make([]float64, len(h[0]))
	// summation
	for _, runHist := range h {
		for round, value := range runHist {
			sumHist[round] = sumHist[round] + value
		}
	}
	// divide by runs
	aveHist := make([]float64, len(h[0]))
	for round, value := range sumHist {
		aveHist[round] = value / float64(len(h))
	}
	return aveHist
}
