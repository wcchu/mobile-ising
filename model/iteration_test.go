package main_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	model "github.com/wcchu/mobile-ising/model"
)

// sampState is a map with 5x5 sites used in multiple tests
var sampState model.State = model.State{
	0: model.SiteInfo{Loc: model.Location{X: 0, Y: 0}, Spin: 1},
	1: model.SiteInfo{Loc: model.Location{X: 0, Y: 1}, Spin: -1},
	2: model.SiteInfo{Loc: model.Location{X: 0, Y: 2}, Spin: 1},
	3: model.SiteInfo{Loc: model.Location{X: 0, Y: 3}, Spin: -1},
	4: model.SiteInfo{Loc: model.Location{X: 0, Y: 4}, Spin: 1},

	5: model.SiteInfo{Loc: model.Location{X: 1, Y: 0}, Spin: -1},
	6: model.SiteInfo{Loc: model.Location{X: 1, Y: 1}, Spin: 1},
	7: model.SiteInfo{Loc: model.Location{X: 1, Y: 2}, Spin: -1},
	8: model.SiteInfo{Loc: model.Location{X: 1, Y: 3}, Spin: 1},
	9: model.SiteInfo{Loc: model.Location{X: 1, Y: 4}, Spin: -1},

	10: model.SiteInfo{Loc: model.Location{X: 2, Y: 0}, Spin: 1},
	11: model.SiteInfo{Loc: model.Location{X: 2, Y: 1}, Spin: -1},
	12: model.SiteInfo{Loc: model.Location{X: 2, Y: 2}, Spin: 1},
	13: model.SiteInfo{Loc: model.Location{X: 2, Y: 3}, Spin: -1},
	14: model.SiteInfo{Loc: model.Location{X: 2, Y: 4}, Spin: 1},

	15: model.SiteInfo{Loc: model.Location{X: 3, Y: 0}, Spin: -1},
	16: model.SiteInfo{Loc: model.Location{X: 3, Y: 1}, Spin: 1},
	17: model.SiteInfo{Loc: model.Location{X: 3, Y: 2}, Spin: -1},
	18: model.SiteInfo{Loc: model.Location{X: 3, Y: 3}, Spin: 1},
	19: model.SiteInfo{Loc: model.Location{X: 3, Y: 4}, Spin: -1},

	20: model.SiteInfo{Loc: model.Location{X: 4, Y: 0}, Spin: 1},
	21: model.SiteInfo{Loc: model.Location{X: 4, Y: 1}, Spin: -1},
	22: model.SiteInfo{Loc: model.Location{X: 4, Y: 2}, Spin: 1},
	23: model.SiteInfo{Loc: model.Location{X: 4, Y: 3}, Spin: -1},
	24: model.SiteInfo{Loc: model.Location{X: 4, Y: 4}, Spin: 1},
}

func TestIterate(t *testing.T) {
	tests := []struct {
		temp  float64
		size  int
		state model.State
	}{
		{ // initial state
			temp:  1.0,
			size:  5,
			state: sampState,
		},
	}

	seed := time.Now().UTC().UnixNano()
	rand.Seed(seed)

	for _, tt := range tests {
		finalState, _ := model.Iterate(tt.size, tt.state, 12, tt.temp)
		change := make(model.State)
		for id, finalSite := range finalState {
			initSite := tt.state[id]
			if finalSite.Loc.X != initSite.Loc.X || finalSite.Loc.Y != initSite.Loc.Y || finalSite.Spin != initSite.Spin {
				change[id] = finalSite
			}
		}
		if len(change) > 2 {
			t.Errorf("state was changed more than allowed; changed sites = %+v", change)
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

func TestGetKiez(t *testing.T) {
	centralKiez := model.State{
		6: model.SiteInfo{Loc: model.Location{X: 1, Y: 1}, Spin: 1},
		7: model.SiteInfo{Loc: model.Location{X: 1, Y: 2}, Spin: -1},
		8: model.SiteInfo{Loc: model.Location{X: 1, Y: 3}, Spin: 1},

		11: model.SiteInfo{Loc: model.Location{X: 2, Y: 1}, Spin: -1},
		12: model.SiteInfo{Loc: model.Location{X: 2, Y: 2}, Spin: 1},
		13: model.SiteInfo{Loc: model.Location{X: 2, Y: 3}, Spin: -1},

		16: model.SiteInfo{Loc: model.Location{X: 3, Y: 1}, Spin: 1},
		17: model.SiteInfo{Loc: model.Location{X: 3, Y: 2}, Spin: -1},
		18: model.SiteInfo{Loc: model.Location{X: 3, Y: 3}, Spin: 1},
	}
	tests := []struct {
		id    int
		state model.State
		order int
		kiez  model.State
	}{
		{
			id:    12,
			state: sampState,
			order: 1,
			kiez:  centralKiez,
		},
	}

	for _, tt := range tests {
		pred := model.GetKiez(tt.id, tt.state, tt.order)
		if len(pred) != len(tt.kiez) {
			t.Errorf("expected %d neighbors, got %d", len(tt.kiez), len(pred))
		} else {
			for id, site := range tt.kiez {
				predSite, ok := pred[id]
				if !ok {
					t.Errorf("site id %d not found in prediction", id)
				}
				if predSite.Spin != site.Spin {
					t.Errorf("expected site %d with spin %d, got spin %d", id, site.Spin, predSite.Spin)
				}
				if predSite.Loc.X != site.Loc.X || predSite.Loc.Y != site.Loc.Y {
					t.Errorf("expected site %d with location %d %d, got location %d %d",
						id, site.Loc.X, site.Loc.Y, predSite.Loc.X, predSite.Loc.Y)
				}
			}
		}
	}
}

func TestGetEnergy(t *testing.T) {
	tests := []struct {
		state    model.State
		boundary int
		energy   float64
	}{
		{
			state:    sampState,
			boundary: 5,
			energy:   30.0,
		},
	}

	for _, tt := range tests {
		pred := model.GetEnergy(tt.state, tt.boundary)
		if pred != tt.energy {
			t.Errorf("expected %f, got %f", tt.energy, pred)
		}
	}
}

func TestGetMag(t *testing.T) {
	tests := []struct {
		state model.State
		mag   float64
	}{
		{
			state: sampState,
			mag:   1.0 / 25.0,
		},
	}

	for _, tt := range tests {
		pred := model.GetMag(tt.state)
		if pred != tt.mag {
			t.Errorf("expected %f, got %f", tt.mag, pred)
		}
	}
}
