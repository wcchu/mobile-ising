package main

import (
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
// T = temperature
// N = num of sites
// L = max num of iterations
func evolve(T float64, N, L int, cmean float64, cmax int) ([]State, []Situation) {
	// build connectivity distribution
	ps := GetConnDist(cmean, cmax, forceConns)

	// build initial state/situation
	st := initState(N, ps)
	situ := Situation{Action: "init", Mag: GetMag(st.Spins)}
	// initial history
	var stateHist []State
	var situHist []Situation
	stateHist = append(stateHist, st)
	situHist = append(situHist, situ)

	// evolve state
	bar := pb.StartNew(L)
	for i := 1; i <= L; i++ {
		bar.Increment()
		st, situ = Iterate(st, T)
		stateHist = append(stateHist, st)
		situHist = append(situHist, situ)
		// if ferromagnetic, stop and return
		if math.Abs(situ.Mag) == 1.0 {
			return stateHist, situHist
		}
	}
	bar.FinishPrint("evolution done")

	// max num of iterations reached, return
	return stateHist, situHist
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
