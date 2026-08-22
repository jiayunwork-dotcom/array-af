package beam

import (
	"math"

	"array-af/internal/geometry"
)

// MainlobeAngle 计算主瓣角度（弧度）。
//
// 主瓣条件是空间相位差 |ψ|=0，即
//
//	kd·cosθ + β = 0  →  cosθ = −β/(kd)
//
// 当 |−β/(kd)| ≤ 1 时存在可见的 θ∈[0,π]，返回该角度；
// 否则主瓣指向虚空间（可见区无 |ψ|=0 的点），返回
// found=false，调用方应据此报「主瓣不可见」。
func MainlobeAngle(w geometry.Wave) (theta float64, found bool) {
	cos := w.MainlobePsi()
	cos = geometry.Clamp(cos, -1, 1)
	return math.Acos(cos), w.MainlobeVisible()
}

// MainlobeAngleDeg 是角度制版本。
func MainlobeAngleDeg(w geometry.Wave) (deg float64, found bool) {
	rad, found := MainlobeAngle(w)
	return geometry.RadToDeg(rad), found
}

// Mainlobe 为主瓣分析结果。
type Mainlobe struct {
	// AngleDeg 是主瓣角度（度）。
	AngleDeg float64
	// Visible 表示主瓣是否在可见区。
	Visible bool
	// CosTheta 是主瓣条件的 cosθ。
	CosTheta float64
}

// AnalyzeMainlobe 分析主瓣。
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

// MainlobeToEndfire 返回从当前配置扫描到端射所需的 β 变化量，
// 用于扫 β 场景的终点计算（β_end = −kd）。
func MainlobeToEndfire(arr *geometry.Array) float64 {
	return arr.EndfireBeta()
}
