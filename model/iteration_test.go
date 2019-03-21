package main_test

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	model "github.com/wcchu/mobile-ising/model"
)

func TestIterate(t *testing.T) {
	sampState := model.State{
		Locations: []model.Location{
			model.Location{X: 0.1, Y: 0.1},
			model.Location{X: 0.1, Y: 0.9},
			model.Location{X: 0.9, Y: 0.1},
			model.Location{X: 0.9, Y: 0.9},
		},
		Connections: []int{2, 2, 2, 2},
		Spins:       []int{1, -1, 1, -1},
	}

	tests := []struct {
		temp  float64
		state model.State
		mag   float64
	}{
		{
			temp:  1.0,
			state: sampState,
			mag:   0.0,
		},
	}

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	for _, tt := range tests {
		predState, predMag := model.Iterate(tt.state, tt.temp)
		// connections should not change
		if !reflect.DeepEqual(predState.Connections, tt.state.Connections) {
			t.Error("connections expected not changed but changed")
		}
		// magnetization at most one-spin change
		if math.Abs(predMag-tt.mag) > 2.0/float64(len(tt.state.Spins)) {
			t.Error("magnetization changed more than allowed")
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
