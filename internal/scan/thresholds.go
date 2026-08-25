package scan

import (
	"math"

	"array-af/internal/geometry"
)

func GratingThresholdSpacing(lambda, beta float64) float64 {
	if lambda <= 0 {
		return math.NaN()
	}
	return lambda * (1 - beta/geometry.TwoPi)
}

func ThresholdSpacingRatio(beta float64) float64 {
	return 1 - beta/geometry.TwoPi
}

func RatioBelowThreshold(ratio, beta float64) bool {
	return ratio < ThresholdSpacingRatio(beta)
}
