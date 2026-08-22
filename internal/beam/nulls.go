package beam

import (
	"math"

	"array-af/internal/geometry"
)

// FirstNulls 计算主瓣两侧第一对零点（阵因子零点）。
//
// 阵因子零点满足 sin(N·ψ/2)=0，即 ψ=2πm/N。
// 主瓣两侧最近的零点取 m=±1：
//
//	cosθ = (2πm/N − β)/(kd)
//
// 当 |cosθ|≤1 时该零点在可见区内，否则置 Valid=false。
func FirstNulls(w geometry.Wave) NullPair {
	np := NullPair{}
	// 较小的角度（左侧，靠近 θ=0）对应 m=+1。
	np.LeftDeg, np.LeftValid = nullAngle(w, 1)
	np.RightDeg, np.RightValid = nullAngle(w, -1)
	return np
}

// nullAngle 计算 m 阶零点角度（度）。
func nullAngle(w geometry.Wave, m int) (float64, bool) {
	psi := geometry.TwoPi * float64(m) / float64(w.N)
	cos := (psi - w.Beta) / w.KD
	if cos < -1 || cos > 1 {
		return nan(), false
	}
	cos = geometry.Clamp(cos, -1, 1)
	return geometry.RadToDeg(math.Acos(cos)), true
}

// NullAngles 列出可见区内全部零点（m=±1,±2,… 直至不可见）。
// 用于页面展示零点位置。
func NullAngles(w geometry.Wave) []float64 {
	// 从主瓣向两侧扩展，最多各看 N 个零点（超过 N 个的零点
	// 已远在可见区外或与栅瓣混叠，不再有意义）。
	out := make([]float64, 0, 2*min(w.N, 8))
	for m := -w.N; m <= w.N; m++ {
		if m == 0 {
			continue
		}
		if deg, ok := nullAngle(w, m); ok {
			out = append(out, deg)
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
