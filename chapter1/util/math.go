package util

import "math"

func IsEqual(f1, f2 float64) bool {
	return math.Dim(f1, f2) < 0.01
}
