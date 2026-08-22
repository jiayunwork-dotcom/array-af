package geometry

import "math"

// Wave 描述阵列的工作波数关系，是 ψ 与阵因子的唯一来源。
// k、d、λ 在这一个类型里统一，避免不同文件各自算一遍导致公式分叉。
type Wave struct {
	// N 是阵元数。
	N int
	// K 是自由空间波数 2π/λ。
	K float64
	// KD 是电长度 k·d。
	KD float64
	// D 与 Lambda 是原始输入，供检查与报告。
	D      float64
	Lambda float64
	// Beta 是相位梯度（弧度）。
	Beta float64
}

// NewWave 从已校验的 Params 构造波数关系。
// 调用方必须先 Validate，否则 kd 可能无意义。
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

// Psi 计算空间相位差 ψ(θ) = kd·cosθ + β（弧度制 θ）。
func (w Wave) Psi(theta float64) float64 {
	return w.KD*math.Cos(theta) + w.Beta
}

// PsiDeg 是角度制的 Psi。
func (w Wave) PsiDeg(thetaDeg float64) float64 {
	return w.Psi(DegToRad(thetaDeg))
}

// PsiRange 返回可见区 θ∈[0,π] 内 ψ 的取值范围 [min, max]。
// 由于 cosθ 在 [0,π] 上从 1 单调降到 −1，范围是
// [β−kd, β+kd]。
func (w Wave) PsiRange() (min, max float64) {
	return w.Beta - w.KD, w.Beta + w.KD
}

// MainlobePsi 返回主瓣条件 ψ=0 对应的 cosθ 值。
func (w Wave) MainlobePsi() float64 {
	return -w.Beta / w.KD
}

// MainlobeVisible 判断 |ψ|=0 是否在可见区内，
// 即是否存在 θ∈[0,π] 使 ψ(θ)=0。
func (w Wave) MainlobeVisible() bool {
	cos := w.MainlobePsi()
	return cos >= -1 && cos <= 1
}

// IsBroadside 判断当前相位梯度是否为侧射（β=0，容差内）。
func (w Wave) IsBroadside() bool {
	return math.Abs(w.Beta) < 1e-12
}

// IsEndfire 判断当前相位梯度是否为端射（β=−kd，容差内）。
func (w Wave) IsEndfire() bool {
	return math.Abs(w.Beta+w.KD) < 1e-12
}

// NormalizedSpacing 返回 d/λ（以波长为单位的间距）。
func (w Wave) NormalizedSpacing() float64 {
	return w.D / w.Lambda
}

// ArrayPeak 返回均匀激励阵因子的理论峰值，即阵元数 N。
// 任何配置下 |sin(Nψ/2)/sin(ψ/2)| 的最大值都是 N。
func (w Wave) ArrayPeak() float64 {
	return float64(w.N)
}
