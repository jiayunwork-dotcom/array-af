package geometry

import "math"

const (
	Pi    = math.Pi
	TwoPi = 2 * math.Pi
)

func DegToRad(deg float64) float64 { return deg * Pi / 180 }

func RadToDeg(rad float64) float64 { return rad * 180 / Pi }

type Angle struct{}

func (Angle) Normalize(rad float64) float64 {
	rad = math.Mod(rad, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad
}

func (Angle) NormalizeSigned(rad float64) float64 {
	rad = math.Mod(rad+Pi, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad - Pi
}

func NearEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func Clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func InClosedInterval(v, lo, hi float64) bool {
	return v >= lo && v <= hi
}
