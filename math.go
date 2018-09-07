package assert

import "math"

// Mathematical utils functions.

// Comparison with an absolute epsilon.
func compareAbs(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// Comparison with a relative epsilon.
func compareRel(a, b, epsilon float64) bool {
	diff := math.Abs(a - b)
	largest := math.Max(math.Abs(a), math.Abs(b))
	return diff <= largest*epsilon
}
