package main

import (
	"math"
	"math/rand"
)

// Location of a site
type Location struct {
	x float64
	y float64
}

// State stores the state of the system
type State struct {
	locations     []Location
	connections   []int
	spins         []int
	magnetization float64
}

// evolve runs the evolution of state in time from a random initial state
// T = temperature
// N = num of sites
// L = max num of iterations
func evolve(T float64, N int64, L int64, cmean float64, cmax int64) (State, []float64) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax)

	// build initial state
	s := initState(N, ps)

	// evolve state
	var magHistory []float64
	for i := int64(0); i < L; i++ {
		s = Iterate(s, T)
		magHistory = append(magHistory, s.magnetization)
		// if ferromagnetic, stop and return
		if math.Abs(s.magnetization) == 1.0 {
			return s, magHistory
		}
	}

	// max num of iterations reached
	return s, magHistory
}

// initState creates a random initial state
// N = num of sites
// ps = prob distribution of having k connections
func initState(N int64, ps []float64) State {
	locs := make([]Location, N)
	conns := make([]int, N)
	sps := make([]int, N)
	var mag int64

	for i := int64(0); i < N; i++ {
		locs[i] = Location{x: rand.Float64(), y: rand.Float64()}
		c := AssignConn(ps, rand.Float64())
		if c < 0 {
			return State{}
		}
		conns[i] = c
		sps[i] = rand.Intn(2)*2 - 1
		mag = mag + int64(sps[i])
	}

	return State{
		locations:     locs,
		connections:   conns,
		spins:         sps,
		magnetization: float64(mag) / float64(N),
	}
}
