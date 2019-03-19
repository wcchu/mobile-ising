package model

import (
	"fmt"
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
// t = temperature
// n = num of sites
// m = max num of iterations
func evolve(t float64, n int64, m int64, cmean float64, cmax int64) (State, []float64) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax)

	// build initial state
	s := initState(n, ps)

	// evolve state
	var magHistory []float64
	for i := int64(0); i < m; i++ {
		s = Iterate(s, t)
		magHistory = append(magHistory, s.magnetization)
		// if ferromagnetic, stop and return
		if math.Abs(s.magnetization) == 1.0 {
			fmt.Printf("reached ferromagnetism in %d steps", i+1)
			return s, magHistory
		}
	}

	// max num of iterations reached
	fmt.Printf("final magnetization %f after %d steps", s.magnetization, m)
	return s, magHistory
}

// initState creates a random initial state
func initState(n int64, ps []float64) State {
	var locs []Location
	var conns, sps []int
	var mag int64

	for i := int64(0); i < n; i++ {
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
		spins:         sps,
		magnetization: float64(mag) / float64(n),
	}
}
