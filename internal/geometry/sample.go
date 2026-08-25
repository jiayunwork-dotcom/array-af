package geometry

type SamplePoint struct {
	ThetaDeg float64
	Psi      float64
	AF       float64
	Element  float64
	Pattern  float64
}

func (a *Array) SamplePoints(steps int) []SamplePoint {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	grid := a.ThetaGrid(steps)
	out := make([]SamplePoint, 0, len(grid))
	for _, th := range grid {
		psi := a.Wave.Psi(th)
		af := a.AF(th)
		el := a.ElementFactor(th)
		out = append(out, SamplePoint{
			ThetaDeg: RadToDeg(th),
			Psi:      psi,
			AF:       af,
			Element:  el,
			Pattern:  Pattern(af, el),
		})
	}
	return out
}

func (a *Array) ToRadians(deg []float64) []float64 {
	out := make([]float64, len(deg))
	for i, d := range deg {
		out[i] = DegToRad(d)
	}
	return out
}

func (a *Array) AFAtDeg(deg float64) float64 {
	return a.AF(DegToRad(deg))
}

func (a *Array) PatternAtDeg(deg float64) float64 {
	return a.Pattern(DegToRad(deg))
}
