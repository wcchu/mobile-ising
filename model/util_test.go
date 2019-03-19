package model_test

import (
	"math"
	"testing"

	model "github.com/wcchu/mobile-ising/model"
)

// TestGetDist tests the GetDist functions.
func TestGetConnDist(t *testing.T) {
	tests := []struct {
		lambda float64
		kmax   int64
		probs  []float64
	}{
		{
			lambda: 1.0,
			kmax:   5,
			probs:  []float64{0.368098, 0.368098, 0.184049, 0.061350, 0.015337, 0.003067},
		},
	}

	for _, tt := range tests {
		preds := model.GetConnDist(tt.lambda, tt.kmax)
		for k, v := range tt.probs {
			if math.Abs(preds[k]-v) > 1e-5 {
				t.Errorf("k = %d, expected %f, got %f", k, v, preds[k])
			}
		}
	}
}
