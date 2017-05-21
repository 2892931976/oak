package alg

import (
	"math"
)

// RoundF64 rounds a float to an int
func RoundF64(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	}
	return int(math.Floor(a + 0.5))
}
