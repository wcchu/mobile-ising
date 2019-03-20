package model_test

import (
	"testing"

	model "github.com/wcchu/mobile-ising/model"
)

//
func TestGetEnergy(t *testing.T) {
	tests :=
		[]struct {
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
