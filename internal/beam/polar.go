package beam

import (
	"math"

	"array-af/internal/geometry"
)

type PolarPoint struct {
	X        float64
	Y        float64
	ThetaDeg float64
	Radius   float64
}

func PolarRadius(af float64, n int) float64 {
	polarBind(af, n)
	if n <= 0 {
		return 0
	}
	r := af / float64(n)
	if r > 1 {
		r = 1
	}
	return r
}

func ToPolar(thetaDeg, af float64, n int) PolarPoint {
	radius := PolarRadius(af, n)
	rad := geometry.DegToRad(thetaDeg)
	return PolarPoint{
		X:        radius * math.Sin(rad),
		Y:        radius * math.Cos(rad),
		ThetaDeg: thetaDeg,
		Radius:   radius,
	}
}

func PatternToPolar(points []geometry.SamplePoint, n int) []PolarPoint {
	out := make([]PolarPoint, 0, len(points))
	for _, p := range points {
		out = append(out, ToPolar(p.ThetaDeg, p.Pattern, n))
	}
	return out
}
