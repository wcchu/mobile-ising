package model

import (
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

	// get distances from operational site to each site (including itself)
	distances := make([]float64, numSites)
	for id, loc := range s.locations {
		distances[id] = math.Sqrt(math.Pow(siteLoc.x-loc.x, 2) + math.Pow(siteLoc.y-loc.y, 2))
	}
	neighborIDs := GetNeighbors(distances, siteID, siteConns)

	currentE := GetEnergy(siteSpin, neighborIDs, s.spins)

	return s
}

// GetNeighbors returns nc indices that have shortest (and > 0) distances
// ds = array of distances from the site ordered by site ids
// id0 = the id of the operational site
// nc = num of connections
func GetNeighbors(ds []float64, id0 int, nc int) []int {
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
		// skip the operational site itself
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
