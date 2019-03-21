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

// Situation is the summary of one iteration
type Situation struct {
	Action string
	Mag    float64
}

// Iterate moves the state forward by one step
// s = current state
// T = temperature
func Iterate(st State, T float64) (State, Situation) {
	// choose the operational site
	N := len(st.Locations)
	id := rand.Intn(N)
	site := siteInfo{
		id:    id,
		loc:   st.Locations[id],
		conns: st.Connections[id],
		spin:  st.Spins[id],
	}
	situ := Situation{
		Action: "none",
		Mag:    GetMag(st.Spins),
	}

	// get neighbors from current operational site
	currNeighbors := GetNeighbors(site, st.Locations)
	currE := GetEnergy(site.spin, currNeighbors, st.Spins)

	if rand.Float64() < 0.5 { // try flipping spin
		// in simplest case, flippedE = -currentE, but we calculate it using GetEnergy for completeness
		flippedE := GetEnergy(-site.spin, currNeighbors, st.Spins)

		// if flipping drops energy, flip it;
		// if flipping raises energy, use conditional probability
		dE := flippedE - currE
		if dE < 0 || rand.Float64() < math.Exp(-dE/T) {
			st.Spins[id] = -site.spin
			situ.Action = "flip"
			situ.Mag = GetMag(st.Spins)
		}
	} else { // try moving site but keeping spin; in this case magnetization doesn't change
		candSite := siteInfo{
			id:    id,
			loc:   Location{X: rand.Float64(), Y: rand.Float64()}, // random candidate location
			conns: st.Connections[id],
			spin:  st.Spins[id],
		}
		candNeighbors := GetNeighbors(candSite, st.Locations)
		candE := GetEnergy(site.spin, candNeighbors, st.Spins)

		dE := candE - currE
		if dE < 0 || rand.Float64() < math.Exp(-dE/T) {
			st.Locations[id] = candSite.loc
			situ.Action = "move"
		}
	}
	// if neither action is taken, state and situation stay the same
	return st, situ
}

// GetNeighbors returns nc indices that have shortest (and > 0) distances
// s = siteInfo of the operational site
// locs = locations of all sites
func GetNeighbors(s siteInfo, locs []Location) []int {
	// calculate distances from operational site to all sites
	ds := make([]float64, len(locs))
	for id, loc := range locs {
		ds[id] = math.Sqrt(math.Pow(s.loc.X-loc.X, 2) + math.Pow(s.loc.Y-loc.Y, 2))
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

// GetMag calculates the magnetization for given spins of a state
func GetMag(sps []int) float64 {
	N := len(sps)
	sum := 0
	for _, s := range sps {
		sum = sum + s
	}
	return float64(sum) / float64(N)
}
