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
// outputs: state history, magnetization history, energy shift history, time step reaching convergence
func evolve(T float64, N, L int, cmean float64, cmax int) ([]State, []float64, []float64) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax, forceConns)

	// initialization
	stateHist := make([]State, L+1) // initial state + L iterations
	magHist := make([]float64, L+1)
	enerHist := make([]float64, L+1)
	stateHist[0] = initState(N, ps)
	magHist[0] = GetMag(stateHist[0].Spins)
	enerHist[0] = 0.0

	// evolve state
	ther := false
	end := 0
	bar := pb.StartNew(L)
	for r := 0; r < L; r++ { // loop over rounds
		bar.Increment()

		// check thermalization if it's not declared yet and there are enough past energy shifts
		//nTher := therRounds * N
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

		if math.Abs(magHist[r]) < 1.0 && !ther {
			// not ferromagnetic nor thermalized yet, run iterations for a round
			state := stateHist[r]
			var mag, enerIter, enerRound float64
			for i := 0; i < N; i++ {
				state, mag, enerIter = Iterate(state, T)
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
