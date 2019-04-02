package main

import (
	"fmt"
	"math"
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
	ther := false
	end := 0
	bar := pb.StartNew(L)
	for r := 0; r < L; r++ { // loop over rounds
		bar.Increment()

		// check thermalization if it's not declared yet and there are enough past energy shifts
		// TODO: make better thermalization criteria
		if therRounds > 0 {
			if !ther && r >= therRounds-1 {
				// all energy shifts of the last therRounds rounds have to be zero to thermalize
				var sum float64
				for _, e := range enerHist[r-therRounds+1 : r+1] {
					sum = sum + math.Abs(e)
				}
				if sum == 0.0 {
					ther = true
				}
			}
		}

		if math.Abs(magHist[r]) < 1.0 && !ther {
			// not ferromagnetic nor thermalized yet, run iterations for a round
			state := stateHist[r]
			var mag, enerIter, enerRound float64
			for i := 0; i < D*D; i++ { // i is the id of the site
				state, mag, enerIter = Iterate(D, state, i, T)
				enerRound = enerRound + enerIter
			}
			stateHist[r+1], magHist[r+1], enerHist[r+1] = state, mag, enerRound
			end = r + 1 // the time index we want to record up to
		} else {
			// ferromagnetic or thermalized
			magHist[r+1] = magHist[r]
			enerHist[r+1] = 0.0
		}
	}
	bar.FinishPrint(fmt.Sprintf("evolution ends at %d rounds of iterations for T = %f", end, T))

	// return only the evolving part of history
	return stateHist[0 : end+1], magHist[0 : end+1], enerHist[0 : end+1]
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
