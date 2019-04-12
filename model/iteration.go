package main

import (
	"math"
	"math/rand"
)

// Iterate moves the state forward by one step
// inputs: D = size of map, inputState = input state, T = temperature
// outputs: state after iteration
func Iterate(D int, inputState State, siteID int, T float64) State {
	// duplicate the state that can change in this iteration
	state := make(State)
	for k, v := range inputState {
		state[k] = v
	}
	site := state[siteID]

	// get neighbors in 5x5 neighborhood of current operational site
	kiez := GetKiez(siteID, state, 2)
	currE := GetEnergy(kiez, D)

	if rand.Float64() > iterMode { // try flipping spin
		site.Spin = -site.Spin // flip
		kiez[siteID] = site    // insert update site into kiez
		flipE := GetEnergy(kiez, D)

		// if flipping drops energy, flip it;
		// if flipping raises energy, use conditional probability
		dE := flipE - currE
		if dE < 0 || rand.Float64() < ExcProb(dE, T) {
			// execute flipping
			state[siteID] = site
		}
	} else { // try moving site but keeping spin; in this case magnetization doesn't change
		moveID := PickMoveID(siteID, kiez, D)
		kiez.Move(siteID, moveID)
		moveE := GetEnergy(kiez, D)

		// if moving drops energy, flip it;
		// if moving raises energy, use conditional probability
		dE := moveE - currE
		if dE < 0 || rand.Float64() < ExcProb(dE, T) {
			// execute moving
			state.Move(siteID, moveID)
		}
	}
	// if neither action is taken, state and magnetization stay the same
	return state
}

// GetKiez returns the H-th order neighborhood around the operational site
// Number of returned sites = (2 * H + 1) ** 2 (including the operational site)
// inputs: id0 = id of operational site, state = state of system, H = see above
// outputs: neighborhood in State format
func GetKiez(id0 int, state State, H int) State {
	// validate size of kiez
	D := len(state)
	if 2*H+1 > D {
		panic("size of neighborhood larger then size of map")
	}
	// collect sites to form kiez
	loc0 := state[id0].Loc
	kiez := make(State)
	for id, site := range state {
		// the x and y distances between loc and loc0 have to be within range of H, considering boundary
		if dx, dy := DistXY(site.Loc, loc0, D); dx <= H && dy <= H {
			kiez[id] = SiteInfo{Loc: site.Loc, Spin: site.Spin}
		}
	}
	return kiez
}

// ExcProb returns the probability of excitation with given dE and T
func ExcProb(dE, T float64) float64 {
	if T == 0.0 { // exp(-dE/T) = 0
		return 0.0
	}
	return math.Exp(-dE / T)
}

// GetEnergy calculates the energy for H_i = -K sum(S_i * S_j) where i, j are 1st order neighbors
// inputs: neighborhood in State format, D = size of full map
// output: energy defined by Hamiltonian within the neighborhood
func GetEnergy(kiez State, D int) float64 {
	// set K to -1
	var K float64 = -1

	// calculate the summation part
	sum := 0
	for _, iSite := range kiez {
		for _, jSite := range kiez {
			if dx, dy := DistXY(iSite.Loc, jSite.Loc, D); dx+dy == 1 {
				sum = sum + iSite.Spin*jSite.Spin
			}
		}
	}

	// the summation counts twice for the same pair (i and j are symmetric)
	return K * float64(sum) / 2.0
}

// GetMag calculates the magnetization for given spins of a state
// input: state = State
// output: magnetization (between 0 and 1)
func GetMag(state State) float64 {
	N := len(state)
	var sum int
	for _, site := range state {
		sum = sum + site.Spin
	}
	return math.Abs(float64(sum)) / float64(N)
}

// DistXY calculates the distances in x and in y separately between two locations considering periodic boundary of size D
func DistXY(loc1, loc2 Location, D int) (int, int) {
	dxIn := Abs(loc1.X - loc2.X) // x distance in original map
	dx := Min(dxIn, D-dxIn)      // x distance considering cross-boundary
	dyIn := Abs(loc1.Y - loc2.Y) // y distance in original map
	dy := Min(dyIn, D-dyIn)      // y distance considering cross-boundary
	return dx, dy
}

// PickMoveID picks a direct neighbor that the operational site wants to swap with
// inputs: id0 = operational site id, kiez = neighborhood, D = size of full map
// output: id of the picked neighbor
func PickMoveID(id0 int, kiez State, D int) int {
	site0 := kiez[id0]
	// collect candidate IDs
	var candidates []int
	for id, site := range kiez {
		if dx, dy := DistXY(site.Loc, site0.Loc, D); dx+dy == 1 {
			candidates = append(candidates, id)
		}
	}
	// randomly pick one of candidates
	r := rand.Intn(len(candidates))
	return candidates[r]
}

// Move exchanges the locations of two sites while keeping their ids and spins unchanged
func (state State) Move(id1, id2 int) {
	site1 := SiteInfo{
		Loc:  state[id2].Loc,
		Spin: state[id1].Spin,
	}
	site2 := SiteInfo{
		Loc:  state[id1].Loc,
		Spin: state[id2].Spin,
	}
	state[id1] = site1
	state[id2] = site2
	return
}
