package beam

import (
	"math"

	"array-af/internal/geometry"
)

func MainlobeAngle(w geometry.Wave) (theta float64, found bool) {
	cos := w.MainlobePsi()
	cos = geometry.Clamp(cos, -1, 1)
	return math.Acos(cos), w.MainlobeVisible()
}

func MainlobeAngleDeg(w geometry.Wave) (deg float64, found bool) {
	rad, found := MainlobeAngle(w)
	return geometry.RadToDeg(rad), found
}

type Mainlobe struct {
	AngleDeg float64
	Visible  bool
	CosTheta float64
}

func AnalyzeMainlobe(w geometry.Wave) Mainlobe {
	cos := w.MainlobePsi()
	visible := cos >= -1 && cos <= 1
	c := geometry.Clamp(cos, -1, 1)
	return Mainlobe{
		AngleDeg: geometry.RadToDeg(math.Acos(c)),
		Visible:  visible,
		CosTheta: cos,
	}
}

func MainlobeToEndfire(arr *geometry.Array) float64 {
	return arr.EndfireBeta()
}
