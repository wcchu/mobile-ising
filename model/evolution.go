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
	locations   []Location
	connections []int
	spins       []int
}

// evolve runs the evolution of state in time from a random initial state
// T = temperature
// N = num of sites
// L = max num of iterations
func evolve(T float64, N int64, L int64, cmean float64, cmax int64) (State, []Situation) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax, forceConns)

	var situ Situation
	var hist []Situation

	// build initial state
	st := initState(N, ps)
	hist = append(hist, Situation{action: "init", mag: getMag(st.spins)})

	// evolve state
	for i := int64(1); i <= L; i++ {
		st, situ = Iterate(st, T)
		hist = append(hist, situ)
		// if ferromagnetic, stop and return
		if math.Abs(situ.mag) == 1.0 {
			return st, hist
		}
	}

	// max num of iterations reached, return
	return st, hist
}

// initState creates a random initial state
// N = num of sites
// ps = prob distribution of having k connections
func initState(N int64, ps []float64) State {
	locs := make([]Location, N)
	conns := make([]int, N)
	sps := make([]int, N)

	for i := int64(0); i < N; i++ {
		locs[i] = Location{x: rand.Float64(), y: rand.Float64()}
		c := AssignConn(ps, rand.Float64())
		if c < 0 {
			return State{}
		}
		conns[i] = c
		sps[i] = rand.Intn(2)*2 - 1
	}

	return State{
		locations:   locs,
		connections: conns,
		spins:       sps,
	}
}
