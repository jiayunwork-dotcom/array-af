package geometry

// SamplePoint 是角度网格上的一个采样值。
type SamplePoint struct {
	// ThetaDeg 是角度（度）。
	ThetaDeg float64
	// Psi 是空间相位差（弧度）。
	Psi float64
	// AF 是阵因子。
	AF float64
	// Element 是元因子。
	Element float64
	// Pattern 是总方向图（AF × Element）。
	Pattern float64
}

// SamplePoints 在 θ∈[0,π] 上按 steps 段均匀采样，返回点列。
// steps<1 时用默认段数。
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

// ToRadians 把采样点角度数组转换为弧度数组。
func (a *Array) ToRadians(deg []float64) []float64 {
	out := make([]float64, len(deg))
	for i, d := range deg {
		out[i] = DegToRad(d)
	}
	return out
}

// AFAtDeg 是角度制的快捷查询，供插值使用。
func (a *Array) AFAtDeg(deg float64) float64 {
	return a.AF(DegToRad(deg))
}

// PatternAtDeg 是角度制的快捷查询。
func (a *Array) PatternAtDeg(deg float64) float64 {
	return a.Pattern(DegToRad(deg))
}
