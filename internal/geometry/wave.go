package geometry

import "math"

type Wave struct {
	N      int
	K      float64
	KD     float64
	D      float64
	Lambda float64
	Beta   float64
}

func NewWave(p Params) Wave {
	k := TwoPi / p.Lambda
	return Wave{
		N:      p.N,
		K:      k,
		KD:     k * p.D,
		D:      p.D,
		Lambda: p.Lambda,
		Beta:   p.Beta,
	}
}

func (w Wave) Psi(theta float64) float64 {
	return w.KD*math.Cos(theta) + w.Beta
}

func (w Wave) PsiDeg(thetaDeg float64) float64 {
	return w.Psi(DegToRad(thetaDeg))
}

func (w Wave) PsiRange() (min, max float64) {
	return w.Beta - w.KD, w.Beta + w.KD
}

func (w Wave) MainlobePsi() float64 {
	return -w.Beta / w.KD
}

func (w Wave) MainlobeVisible() bool {
	cos := w.MainlobePsi()
	return cos >= -1 && cos <= 1
}

func (w Wave) IsBroadside() bool {
	return math.Abs(w.Beta) < 1e-12
}

func (w Wave) IsEndfire() bool {
	return math.Abs(w.Beta+w.KD) < 1e-12
}

func (w Wave) NormalizedSpacing() float64 {
	return w.D / w.Lambda
}

func (w Wave) ArrayPeak() float64 {
	return float64(w.N)
}
