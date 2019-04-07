package main_test

import (
	"math"
	"testing"

	model "github.com/wcchu/mobile-ising/model"
)

func TestAveRuns(t *testing.T) {
	tests := []struct {
		hist    [][]float64
		aveHist []float64
	}{
		{
			hist: [][]float64{
				[]float64{1, 2, 3, 4},
				[]float64{2, 3, 4, 5},
			},
			aveHist: []float64{1.5, 2.5, 3.5, 4.5},
		},
	}

	for _, tt := range tests {
		pred := model.AveRuns(tt.hist)
		for round, value := range tt.aveHist {
			if math.Abs(value-pred[round]) > 1e-5 {
				t.Errorf("the %d th round value got %f, expected %f", round, pred[round], value)
			}
		}
	}
}
