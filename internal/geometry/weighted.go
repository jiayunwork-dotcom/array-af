package geometry

import "math"

func WeightedArrayFactor(n int, psi float64, weights []float64) float64 {
	if n <= 0 {
		return 0
	}
	if len(weights) == 0 {
		return ArrayFactor(n, psi)
	}
	real := 0.0
	imag := 0.0
	for i := 0; i < n; i++ {
		w := 1.0
		if i < len(weights) {
			w = weights[i]
		}
		phase := float64(i) * psi
		real += w * math.Cos(phase)
		imag += w * math.Sin(phase)
	}
	return math.Hypot(real, imag)
}

func WeightedArrayFactorNeg(n int, psi float64, weights []float64) float64 {
	if n <= 0 {
		return 0
	}
	if len(weights) == 0 {
		return ArrayFactorNeg(n, psi)
	}
	real := 0.0
	imag := 0.0
	for i := 0; i < n; i++ {
		w := 1.0
		if i < len(weights) {
			w = weights[i]
		}
		phase := float64(i) * psi
		real += w * math.Cos(phase)
		imag += w * math.Sin(phase)
	}
	return math.Sqrt(real*real + imag*imag) * sign(real)
}

func sign(x float64) float64 {
	if x >= 0 {
		return 1
	}
	return -1
}

type WeightedArray struct {
	Array
	Weights []float64
}

func NewWeightedArray(p Params, kind TaperKind, sidelobeDB float64) (*WeightedArray, error) {
	arr, err := NewArray(p)
	if err != nil {
		return nil, err
	}
	w := BuildTaper(kind, p.N, sidelobeDB)
	return &WeightedArray{Array: *arr, Weights: w}, nil
}

func (a *WeightedArray) WAF(theta float64) float64 {
	return WeightedArrayFactor(a.Params.N, a.Wave.Psi(theta), a.Weights)
}

func (a *WeightedArray) WAFDeg(thetaDeg float64) float64 {
	return a.WAF(DegToRad(thetaDeg))
}

func (a *WeightedArray) WeightedPattern(theta float64) float64 {
	return Pattern(a.WAF(theta), a.ElementFactor(theta))
}

func (a *WeightedArray) MaxWAFVisible(steps int) (max, thetaDeg float64) {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	max = -1
	for i := 0; i <= steps; i++ {
		deg := 180 * float64(i) / float64(steps)
		v := a.WAFDeg(deg)
		if v > max {
			max = v
			thetaDeg = deg
		}
	}
	return max, thetaDeg
}

func (a *WeightedArray) PeakReduction() float64 {
	uniform := float64(a.Params.N)
	peak, _ := a.MaxWAFVisible(360)
	if uniform <= 0 {
		return 0
	}
	return 1 - peak/uniform
}

func ApplyTaperToParams(p Params, kind TaperKind, sidelobeDB float64) ([]float64, error) {
	if err := Validate(p); err != nil {
		return nil, err
	}
	return BuildTaper(kind, p.N, sidelobeDB), nil
}
