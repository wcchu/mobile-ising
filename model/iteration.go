package model

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

type twoArrs struct {
	ids  []int
	vals []float64
}

// Iterate moves the state forward by one step
// s = current state
// t = temperature
func Iterate(s State, t float64) State {
	// choose the operational site
	numSites := len(s.locations)
	siteID := rand.Intn(numSites)
	siteLoc := s.locations[siteID]
	siteConns := s.connections[siteID]
	//siteSpin := s.spins[siteID]

	// get distances from operational site to each site (including itself)
	distances := make([]float64, numSites)
	for id, loc := range s.locations {
		distances[id] = math.Sqrt(math.Pow(siteLoc.x-loc.x, 2) + math.Pow(siteLoc.y-loc.y, 2))
	}
	neighborIDs := GetNeighbors(distances, siteConns)
	fmt.Printf("neighbors = %v", neighborIDs)
	return s
}

// GetNeighbors returns c indices that have shortest (and > 0) distances
func GetNeighbors(ds []float64, c int) []int {
	// convert indices of ds to an array
	ids := make([]int, len(ds))
	vals := make([]float64, len(ds))
	for i, d := range ds {
		ids[i] = i
		vals[i] = d
	}

	// sort by value instead of index using Interface
	sort.Sort(twoArrs{ids: ids, vals: vals})

	// exclude the first one which is the site itself (distance = 0)
	return ids[1 : c+1]
}

func (x twoArrs) Len() int           { return len(x.ids) }
func (x twoArrs) Less(i, j int) bool { return x.vals[i] < x.vals[j] }
func (x twoArrs) Swap(i, j int) {
	x.ids[i], x.ids[j] = x.ids[j], x.ids[i]
	x.vals[i], x.vals[j] = x.vals[j], x.vals[i]
}
