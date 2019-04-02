package main

import (
	"math"
	"math/rand"
	"sort"
)

type SiteInfo struct {
	ID    int
	Loc   Location
	Conns int
	Spin  int
}

// Iterate moves the state forward by one step
// inputs: u = input state, T = temperature
// outputs: end state, magnetization, change of energy
func Iterate(stInp State, id int, T float64) (State, float64, float64) {
	// prepare state and mag before iteration
	st := stInp
	st.Locations = make([]Location, len(stInp.Locations))
	st.Spins = make([]int, len(stInp.Spins))
	copy(st.Locations, stInp.Locations)
	copy(st.Spins, stInp.Spins)

	// default magnetization and change of energy
	mag := GetMag(st.Spins)
	shift := 0.0

	site := SiteInfo{
		ID:    id,
		Loc:   st.Locations[id],
		Conns: st.Connections[id],
		Spin:  st.Spins[id],
	}

	// get neighbors from current operational site
	currNeighbors := GetNeighbors(site, st.Locations)
	currE := GetEnergy(site.Spin, currNeighbors, st.Spins)

	if rand.Float64() > iterMode { // try flipping spin
		// in simplest case, flippedE = -currentE, but we calculate it using GetEnergy for completeness
		flipE := GetEnergy(-site.Spin, currNeighbors, st.Spins)

		// if flipping drops energy, flip it;
		// if flipping raises energy, use conditional probability
		dE := flipE - currE
		if dE < 0 || rand.Float64() < ExcProb(dE, T) {
			// execute flipping
			st.Spins[id] = -site.Spin
			mag = GetMag(st.Spins)
			shift = dE
		}
	} else { // try moving site but keeping spin; in this case magnetization doesn't change
		moveSite := SiteInfo{
			ID:    id,
			Loc:   Location{X: rand.Float64(), Y: rand.Float64()}, // random candidate location
			Conns: site.Conns,
			Spin:  site.Spin,
		}
		moveNeighbors := GetNeighbors(moveSite, st.Locations)
		moveE := GetEnergy(site.Spin, moveNeighbors, st.Spins)

		// if moving drops energy, flip it;
		// if moving raises energy, use conditional probability
		dE := moveE - currE
		if dE < 0 || rand.Float64() < ExcProb(dE, T) {
			// execute moving
			st.Locations[id] = moveSite.Loc
			// mag doesn't change with moving
			shift = dE
		}
	}
	// if neither action is taken, state and magnetization stay the same
	return st, mag, shift
}

// ExcProb returns the probability of excitation with given dE and T
func ExcProb(dE, T float64) float64 {
	if T == 0.0 { // exp(-dE/T) = 0
		return 0.0
	}
	return math.Exp(-dE / T)
}

// GetNeighbors returns nc indices that have shortest (and > 0) distances
// s = SiteInfo of the operational site
// locs = locations of all sites
func GetNeighbors(s SiteInfo, locs []Location) []int {
	// make mirror images as extension of original map to satisfy periodic consition
	N := len(locs)                  // num of sites
	bigMap := make([]Location, N*9) // site locations on the 9-block extended map
	var mapID int = 0
	for _, xShift := range []int{-1, 0, 1} {
		for _, yShift := range []int{-1, 0, 1} {
			for id, loc := range locs {
				bigMap[id+mapID*N] = Location{X: loc.X + float64(xShift), Y: loc.Y + float64(yShift)}
			}
			mapID = mapID + 1
		}
	}

	// calculate distances from operational site to all sites
	ds := make([]float64, N*9)
	for id, loc := range bigMap {
		ds[id] = math.Sqrt(math.Pow(s.Loc.X-loc.X, 2) + math.Pow(s.Loc.Y-loc.Y, 2))
	}

	// sort sites by the distance to operational site
	// convert indices of ds to an array
	ids := make([]int, N*9)
	vals := make([]float64, N*9)
	for i, d := range ds {
		ids[i] = i
		vals[i] = d
	}
	// sort by value instead of index using Interface
	sort.Sort(TwoArrs{IDs: ids, Vals: vals})

	// collect neighbors
	var neighbors []int
	for _, id := range ids {
		// convert extended site id on extended map to corresponding original id
		realID := id % N
		// skip the operational site as its own neighbor
		if realID != s.ID {
			neighbors = append(neighbors, realID)
		}
		// return when there are enough neighbors
		if len(neighbors) == s.Conns {
			return neighbors
		}
	}
	return []int{}
}

// GetEnergy calculates the energy for H_i = -K sum(S_i * S_j)
// s0 = the spin of operational site
// ns = the ids of the neighbor sites
// ss = the list of spins of the whole system
func GetEnergy(s0 int, ns, ss []int) float64 {
	// set K to -1 for simplicity
	var K float64 = -1

	// calculate the summation part
	sum := int(0)
	for _, id := range ns {
		sum = sum + s0*ss[id]
	}

	return K * float64(sum)
}

// GetMag calculates the magnetization for given spins of a state
func GetMag(sps []int) float64 {
	N := len(sps)
	sum := 0
	for _, s := range sps {
		sum = sum + s
	}
	return float64(sum) / float64(N)
}
