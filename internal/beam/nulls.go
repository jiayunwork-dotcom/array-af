package beam

import (
	"math"

	"array-af/internal/geometry"
)

func FirstNulls(w geometry.Wave) NullPair {
	np := NullPair{}
	np.LeftDeg, np.LeftValid = nullAngle(w, 1)
	np.RightDeg, np.RightValid = nullAngle(w, -1)
	return np
}

func nullAngle(w geometry.Wave, m int) (float64, bool) {
	psi := geometry.TwoPi * float64(m) / float64(w.N)
	cos := (psi - w.Beta) / w.KD
	if cos < -1 || cos > 1 {
		return nan(), false
	}
	cos = geometry.Clamp(cos, -1, 1)
	return geometry.RadToDeg(math.Acos(cos)), true
}

func NullAngles(w geometry.Wave) []float64 {
	out := make([]float64, 0, 2*min(w.N, 8))
	for m := -w.N; m <= w.N; m++ {
		if m == 0 {
			continue
		}
		if deg, ok := nullAngle(w, m); ok {
			out = append(out, deg)
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
