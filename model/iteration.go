package model

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

// Iterate moves the state forward by one step
// s = current state
// t = temperature
func Iterate(s State, t float64) State {
	// choose the operational site
	numSites := len(s.locations)
	siteID := rand.Intn(numSites)
	siteLoc := s.locations[siteID]
	siteConns := s.connections[siteID]
	siteSpin := s.spins[siteID]

	// get neighbors from current operational site
	currNeighbors := GetNeighbors(siteID, siteLoc, s.locations, siteConns)
	currE := GetEnergy(siteSpin, currNeighbors, s.spins)
	fmt.Printf("current energy = %f", currE)

	if rand.Float64() < 0.5 { // try flipping spin
		// in simplest case, flippedE = -currentE, but we calculate it using GetEnergy for completeness
		flippedE := GetEnergy(-siteSpin, currNeighbors, s.spins)
		fmt.Printf("flipped energy = %f", flippedE)

		excitation := flippedE - currE
		// if flipping drops energy, flip it
		// if flipping raises energy, use condition
		if excitation < 0 || rand.Float64() < math.Exp(-excitation/t) {
			s.spins[siteID] = -siteSpin
			return s
		}

	} else {
		// if moved
		candLoc := Location{x: rand.Float64(), y: rand.Float64()}
		fmt.Printf("proposed new location = %+v", candLoc)
		candNeighbors := GetNeighbors(siteID, candLoc, s.locations, siteConns)
		candE := GetEnergy(siteSpin, candNeighbors, s.spins)
		fmt.Printf("energy if moved to proposed location = %f", candE)

	}

	return s
}

// GetNeighbors returns nc indices that have shortest (and > 0) distances
// ds = array of distances from the site ordered by site ids
// id0 = the id of the operational site
// nc = num of connections
func GetNeighbors(id0 int, loc0 Location, locs []Location, nc int) []int {
	// calculate distances from all operational site to all sites
	distances := make([]float64, len(locs))
	for id, loc := range locs {
		distances[id] = math.Sqrt(math.Pow(loc0.x-loc.x, 2) + math.Pow(loc0.y-loc.y, 2))
	}

	// convert indices of ds to an array
	ids := make([]int, len(distances))
	vals := make([]float64, len(distances))
	for i, d := range distances {
		ids[i] = i
		vals[i] = d
	}

	// sort by value instead of index using Interface
	sort.Sort(TwoArrs{IDs: ids, Vals: vals})

	// collect neighbors
	var neighbors []int
	for _, id := range ids {
		// skip the operational site as its own neighbor
		if id != id0 {
			neighbors = append(neighbors, id)
		}
		// return when there are enough neighbors
		if len(neighbors) == nc {
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
