package geometry

import (
	"fmt"
	"math"
)

// Element 描述单个阵元的辐射模式（元因子）。
// 元因子与阵因子相乘得到阵列总方向图。
type Element struct {
	kind ElementType
}

// NewElement 按类型构造元因子。
func NewElement(kind ElementType) Element {
	return Element{kind: kind}
}

// ParseElement 把字符串解析为元因子类型。
// 支持的取值：iso（各向同性）、dipole（短偶极 sinθ）。
func ParseElement(s string) (ElementType, error) {
	switch ElementType(s) {
	case ElementIso, "":
		return ElementIso, nil
	case ElementDipole:
		return ElementDipole, nil
	default:
		return "", fmt.Errorf("unknown element factor %q (want iso or dipole)", s)
	}
}

// Kind 返回规范化类型。
func (e Element) Kind() ElementType {
	if e.kind == ElementDipole {
		return ElementDipole
	}
	return ElementIso
}

// Factor 返回元因子在 θ（弧度）处的值。
// 各向同性恒为 1；短偶极沿轴线（θ=0 或 π）为 0，侧射（θ=π/2）为 1。
func (e Element) Factor(theta float64) float64 {
	if e.Kind() == ElementDipole {
		s := math.Sin(theta)
		if s < 0 {
			s = -s
		}
		return fillEl(s)
	}
	return 1
}

// FactorDeg 是角度制的 Factor。
func (e Element) FactorDeg(thetaDeg float64) float64 {
	return e.Factor(DegToRad(thetaDeg))
}

// Maximum 返回元因子在可见区的最大值。
// 各向同性为 1；短偶极在 θ=π/2 处取 1。
func (e Element) Maximum() float64 {
	return 1
}

// IsDipole 判断是否为短偶极元因子。
func (e Element) IsDipole() bool {
	return e.Kind() == ElementDipole
}

// String 返回可打印类型名。
func (e ElementType) String() string {
	if e == ElementDipole {
		return string(ElementDipole)
	}
	return string(ElementIso)
}

// String 返回可打印类型名。
func (e Element) String() string {
	return e.Kind().String()
}
