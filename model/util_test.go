package main_test

import (
	"math"
	"sort"
	"testing"

	model "github.com/wcchu/mobile-ising/model"
)

func TestTwoArrs(t *testing.T) {
	tests :=
		[]struct {
			object    model.TwoArrs
			sortedIDs []int
		}{
			{
				object: model.TwoArrs{
					IDs:  []int{1, 2, 3, 4, 5},
					Vals: []float64{4, 2, 1, 3, 5},
				},
				sortedIDs: []int{3, 2, 4, 1, 5},
			},
		}

	for _, tt := range tests {
		sort.Sort(tt.object)
		for k, v := range tt.sortedIDs {
			if tt.object.IDs[k] != v {
				t.Errorf("k = %d, expected %d, got %d", k, v, tt.object.IDs[k])
			}
		}
	}
}

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

func TestAssignConn(t *testing.T) {
	tests := []struct {
		probs []float64
		r     float64
		k     int
	}{
		{
			probs: []float64{0.2, 0.3, 0.4, 0.1},
			r:     0.1,
			k:     0,
		},
		{
			probs: []float64{0.2, 0.3, 0.4, 0.1},
			r:     0.7,
			k:     2,
		},
	}

	for _, tt := range tests {
		pred := model.AssignConn(tt.probs, tt.r)
		if pred != tt.k {
			t.Errorf("expected %d, got %d", pred, tt.k)
		}
	}
}
