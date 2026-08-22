package geometry

import "math"

// 常用常量。
const (
	// Pi 是圆周率。
	Pi = math.Pi
	// TwoPi 是 2π。
	TwoPi = 2 * math.Pi
)

// DegToRad 把角度从度转换到弧度。
func DegToRad(deg float64) float64 { return deg * Pi / 180 }

// RadToDeg 把角度从弧度转换到度。
func RadToDeg(rad float64) float64 { return rad * 180 / Pi }

// Angle 提供角度规范化与比较工具。
type Angle struct{}

// Normalize 把角度折叠到 [0, 2π)。
func (Angle) Normalize(rad float64) float64 {
	rad = math.Mod(rad, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad
}

// NormalizeSigned 把角度折叠到 (−π, π]。
func (Angle) NormalizeSigned(rad float64) float64 {
	rad = math.Mod(rad+Pi, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad - Pi
}

// NearEqual 判断两个浮点数在 tol 容差内是否相等。
func NearEqual(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

// Clamp 把 v 限制在 [lo, hi]。
func Clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// InClosedInterval 判断 v 是否在 [lo, hi]（含端点）。
func InClosedInterval(v, lo, hi float64) bool {
	return v >= lo && v <= hi
}
