package geometry

type Array struct {
	Params  Params
	Wave    Wave
	Element Element
}

func NewArray(p Params) (*Array, error) {
	if err := Validate(p); err != nil {
		return nil, err
	}
	if err := CheckFinite(p); err != nil {
		return nil, err
	}
	return &Array{
		Params:  p,
		Wave:    NewWave(p),
		Element: NewElement(p.ElementKind()),
	}, nil
}

func (a *Array) AF(theta float64) float64 {
	return ArrayFactor(a.Wave.N, a.Wave.Psi(theta))
}

func (a *Array) AFDeg(thetaDeg float64) float64 {
	return a.AF(DegToRad(thetaDeg))
}

func (a *Array) ElementFactor(theta float64) float64 {
	return a.Element.Factor(theta)
}

func (a *Array) ElementFactorDeg(thetaDeg float64) float64 {
	return a.Element.FactorDeg(thetaDeg)
}

func (a *Array) Pattern(theta float64) float64 {
	return Pattern(a.AF(theta), a.ElementFactor(theta))
}

func (a *Array) PatternDeg(thetaDeg float64) float64 {
	return a.Pattern(DegToRad(thetaDeg))
}

func (a *Array) Peak() float64 {
	return a.Wave.ArrayPeak()
}

func (a *Array) Kd() float64 { return a.Wave.KD }

func (a *Array) Wavenumber() float64 { return a.Wave.K }

func (a *Array) SpacingRatio() float64 { return a.Wave.NormalizedSpacing() }

func (a *Array) EndfireBeta() float64 { return -a.Wave.KD }

func (a *Array) ThetaGrid(steps int) []float64 {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	out := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		out = append(out, Pi*float64(i)/float64(steps))
	}
	return out
}

func (a *Array) ThetaDegGrid(steps int) []float64 {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	out := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		out = append(out, 180*float64(i)/float64(steps))
	}
	return out
}

func (a *Array) MaxAFVisible(steps int) (max, theta float64) {
	grid := a.ThetaGrid(steps)
	max = -1
	for _, th := range grid {
		v := a.AF(th)
		if v > max {
			max = v
			theta = th
		}
	}
	return max, theta
}

func (a *Array) StepRad(steps int) float64 {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	return Pi / float64(steps)
}
