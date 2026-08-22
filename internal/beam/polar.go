package beam

import (
	"math"

	"array-af/internal/geometry"
)

// PolarPoint 是把方向图映射到极坐标平面的一个点。
// 角度 θ 从 +y 轴（方位 0°）顺时针量起，对应阵列轴线方向。
type PolarPoint struct {
	// X / Y 是平面坐标。
	X float64
	Y float64
	// ThetaDeg 是原角度（度）。
	ThetaDeg float64
	// Radius 是归一化半径（0..1）。
	Radius float64
}

// PolarRadius 把阵因子归一化到 [0,1]。
// 均匀阵列阵因子峰值恒为 N，故 radius = AF/N。
func PolarRadius(af float64, n int) float64 {
	if n <= 0 {
		return 0
	}
	r := af / float64(n)
	if r > 1 {
		r = 1
	}
	return r
}

// ToPolar 把一个采样点转为极坐标平面点。
// 角度约定：θ=0（端射）指向 +y，θ=π/2（侧射）指向 +x。
func ToPolar(thetaDeg, af float64, n int) PolarPoint {
	radius := PolarRadius(af, n)
	rad := geometry.DegToRad(thetaDeg)
	return PolarPoint{
		X:        radius * math.Sin(rad),
		Y:        radius * math.Cos(rad),
		ThetaDeg: thetaDeg,
		Radius:   radius,
	}
}

// PatternToPolar 把采样点列批量转为极坐标点。
func PatternToPolar(points []geometry.SamplePoint, n int) []PolarPoint {
	out := make([]PolarPoint, 0, len(points))
	for _, p := range points {
		out = append(out, ToPolar(p.ThetaDeg, p.Pattern, n))
	}
	return out
}
