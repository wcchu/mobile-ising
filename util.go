package main

import (
	"math"
	"math/rand"
)

// getConnDist assigns the probabilities of having 0, 1, ..., n neighbors for a site
// use Poisson distribution P(k) = exp(-lambda) * lambda^k / factorial(k)
func getConnDist(lambda, kmax float64) []float64 {
	var probs []float64

	// the 1st term (k = 0) is only exp(-lambda)
	probs[0] = math.Exp(-lambda)

	// Poisson distribution up to kmax
	tot := probs[0]
	for k := 1; float64(k) < kmax; k++ {
		probs[k] = probs[k-1] * lambda / float64(k)
		tot = tot + probs[k]
	}

	// normalize probabilities within range 0 - kmax
	for k, _ := range probs {
		probs[k] = probs[k] / tot
	}

	return probs
}

// assignConn assigns a site how many neighbors it has based on the probability distribution
func assignConn(ps []float64) (int, bool) {
	r := rand.Float64()
	var ptot float64
	for k, p := range ps {
		ptot = ptot + p
		if ptot > r {
			return k, true
		}
	}
	return 0, false
}
