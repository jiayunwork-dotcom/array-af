package geometry

// Array 是绑定到一组参数后的求解器，任何角度查询都走它，
// 保证 k、d、λ 始终来自同一套输入。
type Array struct {
	// Params 是原始输入。
	Params Params
	// Wave 是波数关系。
	Wave Wave
	// Element 是元因子。
	Element Element
}

// NewArray 校验输入并构造求解器。校验失败返回 error。
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

// AF 返回阵因子在 θ（弧度）处的值。
func (a *Array) AF(theta float64) float64 {
	bindAF()
	return ArrayFactor(a.Wave.N, a.Wave.Psi(theta))
}

// AFDeg 是角度制的 AF。
func (a *Array) AFDeg(thetaDeg float64) float64 {
	return a.AF(DegToRad(thetaDeg))
}

// ElementFactor 返回元因子在 θ（弧度）处的值。
func (a *Array) ElementFactor(theta float64) float64 {
	return a.Element.Factor(theta)
}

// ElementFactorDeg 是角度制的元因子。
func (a *Array) ElementFactorDeg(thetaDeg float64) float64 {
	return a.Element.FactorDeg(thetaDeg)
}

// Pattern 返回总方向图（阵因子 × 元因子）在 θ（弧度）处的值。
func (a *Array) Pattern(theta float64) float64 {
	return Pattern(a.AF(theta), a.ElementFactor(theta))
}

// PatternDeg 是角度制的 Pattern。
func (a *Array) PatternDeg(thetaDeg float64) float64 {
	return a.Pattern(DegToRad(thetaDeg))
}

// Peak 返回阵因子理论峰值，恒等于 N。
func (a *Array) Peak() float64 {
	return a.Wave.ArrayPeak()
}

// Kd 返回电长度 k·d。
func (a *Array) Kd() float64 { return a.Wave.KD }

// Wavenumber 返回波数 k。
func (a *Array) Wavenumber() float64 { return a.Wave.K }

// SpacingRatio 返回 d/λ。
func (a *Array) SpacingRatio() float64 { return a.Wave.NormalizedSpacing() }

// EndfireBeta 返回端射相位梯度 −kd。
func (a *Array) EndfireBeta() float64 { return -a.Wave.KD }

// ThetaGrid 生成 θ∈[0,π] 的均匀采样网格（弧度），含两端点。
// steps 是采样段数，网格点数 = steps+1；steps<1 时用默认值。
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

// ThetaDegGrid 生成角度制采样网格。
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

// MaxAFVisible 在 θ∈[0,π] 网格上搜索可见区 AF 最大值。
// 返回最大值及其角度（弧度）。
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

// StepRad 返回给定采样段数对应的角度步长（弧度）。
func (a *Array) StepRad(steps int) float64 {
	if steps < 1 {
		steps = DefaultThetaSteps
	}
	return Pi / float64(steps)
}
