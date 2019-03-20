package main

import (
	"math"
	"math/rand"
	"sort"
)

type siteInfo struct {
	id    int
	loc   Location
	conns int
	spin  int
}

// Iterate moves the state forward by one step
// s = current state
// T = temperature
func Iterate(st State, T float64) State {
	// choose the operational site
	N := len(st.locations)
	id := rand.Intn(N)
	site := siteInfo{
		id:    id,
		loc:   st.locations[id],
		conns: st.connections[id],
		spin:  st.spins[id],
	}

	// get neighbors from current operational site
	currNeighbors := GetNeighbors(site, st.locations)
	currE := GetEnergy(site.spin, currNeighbors, st.spins)

	if rand.Float64() < 0.5 { // try flipping spin
		// in simplest case, flippedE = -currentE, but we calculate it using GetEnergy for completeness
		flippedE := GetEnergy(-site.spin, currNeighbors, st.spins)

		// if flipping drops energy, flip it;
		// if flipping raises energy, use conditional probability
		dE := flippedE - currE
		if dE < 0 || rand.Float64() < math.Exp(-dE/T) {
			st.spins[id] = -site.spin
			return st
		}
	} else { // try moving site but keeping spin
		candSite := siteInfo{
			id:    id,
			loc:   Location{x: rand.Float64(), y: rand.Float64()}, // random candidate location
			conns: st.connections[id],
			spin:  st.spins[id],
		}
		candNeighbors := GetNeighbors(candSite, st.locations)
		candE := GetEnergy(site.spin, candNeighbors, st.spins)

		dE := candE - currE
		if dE < 0 || rand.Float64() < math.Exp(-dE/T) {
			st.locations[id] = candSite.loc
			return st
		}
	}
	// if neither action is taken, return original state
	return st
}

// GetNeighbors returns nc indices that have shortest (and > 0) distances
// s = siteInfo of the operational site
// locs = locations of all sites
func GetNeighbors(s siteInfo, locs []Location) []int {
	// calculate distances from operational site to all sites
	ds := make([]float64, len(locs))
	for id, loc := range locs {
		ds[id] = math.Sqrt(math.Pow(s.loc.x-loc.x, 2) + math.Pow(s.loc.y-loc.y, 2))
	}

	// convert indices of ds to an array
	ids := make([]int, len(ds))
	vals := make([]float64, len(ds))
	for i, d := range ds {
		ids[i] = i
		vals[i] = d
	}

	// sort by value instead of index using Interface
	sort.Sort(TwoArrs{IDs: ids, Vals: vals})

	// collect neighbors
	var neighbors []int
	for _, id := range ids {
		// skip the operational site as its own neighbor
		if id != s.id {
			neighbors = append(neighbors, id)
		}
		// return when there are enough neighbors
		if len(neighbors) == s.conns {
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
