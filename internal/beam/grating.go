package beam

import (
	"math"

	"array-af/internal/geometry"
)

// MaxGratingLobes 是逐个枚举栅瓣的数量上限。
// 间距远大于波长时可见区栅瓣极多，超过上限只报告存在，
// 不逐个展开（HasGrating 仍返回 true）。
const MaxGratingLobes = 1024

// GratingOrders 返回可见区内满足 |ψ|=2π|m|（m≠0）的全部整数阶。
// ψ 在可见区的取值区间是 [β−kd, β+kd]（闭区间，含端点）。
// 超过 MaxGratingLobes 个或区间非有限时返回 nil（存在性另见 HasGrating）。
func GratingOrders(w geometry.Wave) []int {
	lo, hi := w.PsiRange()
	orders, ok := gratingRange(lo, hi)
	if !ok || len(orders) > MaxGratingLobes {
		return nil
	}
	return orders
}

// HasGrating 判断可见区是否存在栅瓣（含端点）。
// 该判定对任意有限输入都成立，不依赖枚举。
func HasGrating(w geometry.Wave) bool {
	lo, hi := w.PsiRange()
	a := lo / geometry.TwoPi
	b := hi / geometry.TwoPi
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if math.IsInf(b, 1) || math.IsInf(a, -1) {
		return true
	}
	if math.IsInf(a, 1) || math.IsInf(b, -1) {
		return false
	}
	return math.Floor(b) >= 1 || math.Ceil(a) <= -1
}

// gratingRange 把 ψ 区间映射到整数阶列表。
func gratingRange(lo, hi float64) ([]int, bool) {
	if math.IsNaN(lo) || math.IsNaN(hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) {
		return nil, false
	}
	a := lo / geometry.TwoPi
	b := hi / geometry.TwoPi
	mLo := int(math.Ceil(a))
	mHi := int(math.Floor(b))
	if mHi < mLo {
		return nil, true
	}
	orders := make([]int, 0, mHi-mLo+1)
	for m := mLo; m <= mHi; m++ {
		if m == 0 {
			continue // ψ=0 是主瓣，不是栅瓣
		}
		orders = append(orders, m)
	}
	return orders, true
}

// GratingAngles 列出可见区内全部栅瓣的角度（度）与阶数。
func GratingAngles(w geometry.Wave) []GratingLobe {
	orders := GratingOrders(w)
	if len(orders) == 0 {
		return nil
	}
	lobes := make([]GratingLobe, 0, len(orders))
	for _, m := range orders {
		psi := float64(m) * geometry.TwoPi
		cos := (psi - w.Beta) / w.KD
		cos = geometry.Clamp(cos, -1, 1)
		deg := geometry.RadToDeg(math.Acos(cos))
		lobes = append(lobes, GratingLobe{
			Order:    m,
			AngleDeg: deg,
			Psi:      psi,
		})
	}
	return lobes
}

// GratingAngleFromBeta 给定相位梯度，返回栅瓣出现情况，供扫 β 场景使用。
func GratingAngleFromBeta(w geometry.Wave, beta float64) []GratingLobe {
	shifted := w
	shifted.Beta = beta
	return GratingAngles(shifted)
}
