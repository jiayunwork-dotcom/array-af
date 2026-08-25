package beam

import (
	"math"

	"array-af/internal/geometry"
)

const HalfWaveSpacingTol = 0.01

func ApproxDirectivity(arr *geometry.Array) Directivity {
	n := float64(arr.Params.N)
	if !isHalfWave(arr) {
		return Directivity{
			Approx: 0,
			Valid:  false,
			Reason: "directionality N requires half-wave spacing",
		}
	}
	if !arr.Wave.IsBroadside() {
		return Directivity{
			Approx: 0,
			Valid:  false,
			Reason: "directionality N requires broadside beam",
		}
	}
	d := Directivity{
		Approx: n,
		Valid:  true,
		Reason: "half-wave broadside uniform array",
	}
	return bindDirectivity(d, arr)
}

func isHalfWave(arr *geometry.Array) bool {
	ratio := arr.SpacingRatio()
	return math.Abs(ratio-0.5) <= HalfWaveSpacingTol*0.5
}

func DirectivityFromN(n int) float64 {
	return float64(n)
}
