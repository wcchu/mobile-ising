package main_test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	model "github.com/wcchu/mobile-ising/model"
)

// sampState is used in multiple tests
var sampState model.State = model.State{
	Locations: []model.Location{
		model.Location{X: 0.11, Y: 0.1},
		model.Location{X: 0.1, Y: 0.91},
		model.Location{X: 0.9, Y: 0.11},
		model.Location{X: 0.9, Y: 0.91},
	},
	Connections: []int{2, 2, 2, 2},
	Spins:       []int{1, -1, 1, -1},
}

func TestIterate(t *testing.T) {
	tests := []struct {
		temp  float64
		state model.State
		mag   float64
		ener  float64
	}{
		{ // initial state
			temp:  1.0,
			state: sampState,
			mag:   0.0,
			ener:  0.0,
		},
	}

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	for _, tt := range tests {
		predState, predMag, predEner := model.Iterate(tt.state, 0, tt.temp)
		// connections should not change
		if !reflect.DeepEqual(predState.Connections, tt.state.Connections) {
			t.Error("connections expected not changed but changed")
		}
		// magnetization at most one-spin change
		if math.Abs(predMag-tt.mag) > 2.0/float64(len(tt.state.Spins)) {
			t.Error("magnetization changed more than allowed")
		}
		// energy shift at most 2
		if math.Abs(predEner-tt.ener) > 2.0 {
			t.Error("energy shifted more than allowed")
		}
		// states
		numCh := 0
		for iloc, loc := range tt.state.Locations {
			if predState.Locations[iloc].X != loc.X || predState.Locations[iloc].Y != loc.Y {
				numCh++
			}
		}
		for ispin, spin := range tt.state.Spins {
			if predState.Spins[ispin] != spin {
				numCh++
			}
		}
		if numCh > 1 {
			t.Errorf("old state = %+v, new state = %+v; more than one site/feature changed", tt.state, predState)
		}
	}
}

func TestExcProb(t *testing.T) {
	tests := []struct {
		ener float64
		temp float64
		prob float64
	}{
		{
			ener: 1.0,
			temp: 0.0,
			prob: 0.0,
		},
		{
			ener: 1.0,
			temp: 1.0,
			prob: 0.367879,
		},
	}

	for _, tt := range tests {
		pred := model.ExcProb(tt.ener, tt.temp)
		if math.Abs(pred-tt.prob) > 1e-5 {
			t.Errorf("expected %f, got %f", tt.prob, pred)
		}
	}
}

func TestGetNeighbors(t *testing.T) {
	tests := []struct {
		site      model.SiteInfo
		locs      []model.Location
		neighbors []int
	}{
		{
			site: model.SiteInfo{
				ID:    0,
				Loc:   sampState.Locations[0],
				Conns: 1,
				Spin:  sampState.Spins[0],
			},
			locs:      sampState.Locations,
			neighbors: []int{1},
		},
	}

	for _, tt := range tests {
		pred := model.GetNeighbors(tt.site, tt.locs)
		if len(pred) != len(tt.neighbors) {
			t.Errorf("expected %d neighbors, got %d", len(tt.neighbors), len(pred))
		} else {
			for i, id := range tt.neighbors {
				if pred[i] != id {
					t.Errorf("expected the %d th neighbor with id %d, got id %d", i, id, pred[i])
				}
			}
		}
	}
}

func TestGetEnergy(t *testing.T) {
	tests := []struct {
		thisSpin    int
		neighborIDs []int
		allSpins    []int
		energy      float64
	}{
		{
			thisSpin:    1,
			neighborIDs: []int{1, 3, 6, 7},
			allSpins:    []int{1, -1, 1, -1, 1, -1, 1, -1},
			energy:      2,
		},
	}

	for _, tt := range tests {
		pred := model.GetEnergy(tt.thisSpin, tt.neighborIDs, tt.allSpins)
		if pred != tt.energy {
			t.Errorf("expected %f, got %f", tt.energy, pred)
		}
	}
}

func TestGetMag(t *testing.T) {
	tests := []struct {
		spins []int
		mag   float64
	}{
		{
			spins: []int{1, -1, 1, -1, 1, -1, 1, -1, 1},
			mag:   1.0 / 9.0,
		},
	}

	for _, tt := range tests {
		pred := model.GetMag(tt.spins)
		if pred != tt.mag {
			t.Errorf("expected %f, got %f", tt.mag, pred)
		}
	}
}
