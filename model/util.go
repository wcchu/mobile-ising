package main

// TwoArrs is for sorting indexed values by values
type TwoArrs struct {
	IDs  []int
	Vals []float64
}

// Sort by value when sorting TwoArrs structure
func (x TwoArrs) Len() int           { return len(x.IDs) }
func (x TwoArrs) Less(i, j int) bool { return x.Vals[i] < x.Vals[j] }
func (x TwoArrs) Swap(i, j int) {
	x.IDs[i], x.IDs[j] = x.IDs[j], x.IDs[i]
	x.Vals[i], x.Vals[j] = x.Vals[j], x.Vals[i]
}

// Abs returns the absolute value of integer x
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns the min between x and y integers
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
