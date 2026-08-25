package scan

import (
	"math"
)

func BetaSweep(start, end float64, steps int) []float64 {
	if steps < 1 {
		steps = 32
	}
	out := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		out = append(out, start+(end-start)*t)
	}
	return out
}

func BetaSweepDeg(startDeg, endDeg float64, steps int) []float64 {
	rad := BetaSweep(math.Pi*startDeg/180, math.Pi*endDeg/180, steps)
	out := make([]float64, 0, len(rad))
	for _, r := range rad {
		out = append(out, r*180/math.Pi)
	}
	return out
}

func StepsFromDelta(start, end, delta float64) int {
	if delta <= 0 {
		return 32
	}
	n := int(math.Ceil(math.Abs(end-start) / delta))
	if n < 1 {
		n = 1
	}
	return n
}
