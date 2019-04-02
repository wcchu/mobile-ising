package main_test

import (
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
