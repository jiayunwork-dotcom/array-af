// Package beam 对均匀直线阵做波束层面的分析：
// 主瓣指向、半功率波束宽度、栅瓣、方向性近似与零点位置。
package beam

import (
	"array-af/internal/geometry"
)

// GratingLobe 是可见区内的一个栅瓣。
type GratingLobe struct {
	// Order 是栅瓣的整数阶 m（ψ=2πm，m≠0）。
	Order int
	// AngleDeg 是栅瓣所在角度（度）。
	AngleDeg float64
	// Psi 是栅瓣处的空间相位差（弧度），等于 2π·m。
	Psi float64
}

// NullPair 是一对零点（主瓣两侧）。
type NullPair struct {
	// LeftDeg 与 RightDeg 是第一对零点角度（度）。
	LeftDeg  float64
	RightDeg float64
	// LeftValid / RightValid 标记该侧零点是否在可见区内。
	LeftValid  bool
	RightValid bool
}

// Hpbw 是半功率波束宽度的测量结果。
type Hpbw struct {
	// WidthDeg 是半功率宽度（度）。
	WidthDeg float64
	// LeftDeg / RightDeg 是两个半功率点角度（度）。
	LeftDeg  float64
	RightDeg float64
	// Measurable 为 false 表示无法在可见区内找到完整一对半功率点
	//（例如主瓣不可见，或端射束半宽伸入虚空间）。
	Measurable bool
	// LeftClipped / RightClipped 表示对应侧半功率点落在可见区
	// 边界之外，边界值取 0° 或 180°。
	LeftClipped  bool
	RightClipped bool
}

// Directivity 是方向性近似结果。
type Directivity struct {
	// Approx 是方向性数值；仅在判定有效时有意义。
	Approx float64
	// Valid 表示该近似是否适用于当前配置。
	Valid bool
	// Reason 说明适用或不适用原因。
	Reason string
}

// Result 是一次完整波束分析的输出。
type Result struct {
	// MainlobeAngleDeg 是主瓣角度（度）。
	MainlobeAngleDeg float64
	// MainlobeVisible 表示主瓣（|ψ|=0 的 θ）是否在可见区。
	MainlobeVisible bool
	// MainlobeCos 是主瓣条件的 cosθ 值。
	MainlobeCos float64

	// Hpbw 是半功率宽度测量。
	Hpbw Hpbw

	// HasGrating 表示可见区内是否有栅瓣。
	HasGrating bool
	// GratingLobes 是可见区内全部栅瓣。
	GratingLobes []GratingLobe

	// Nulls 是第一对零点。
	Nulls NullPair

	// Directivity 是方向性近似。
	Directivity Directivity

	// AfPeak 是阵因子可见区峰值（理论上等于 N）。
	AfPeak float64
	// AfPeakThetaDeg 是峰值所在角度（度）。
	AfPeakThetaDeg float64
	// AfPeakMatchesN 表示可见区峰值是否等于阵元数 N。
	AfPeakMatchesN bool
}

// MainlobeDeg 返回主瓣角度（度）；不可见时返回 NaN。
func (r Result) MainlobeDeg() float64 {
	if !r.MainlobeVisible {
		return nan()
	}
	return r.MainlobeAngleDeg
}

// Beam 是把求解器与分析选项绑定后的分析器。
type Beam struct {
	Array *geometry.Array
	// SearchSteps 是可见区峰值/栅瓣搜索的采样段数。
	SearchSteps int
	// HpbwSteps 是 HPBW 扫描的采样段数（步长由此钉死）。
	HpbwSteps int
}

// New 构造分析器。steps 为 0 时使用默认步长。
func New(arr *geometry.Array, searchSteps, hpbwSteps int) *Beam {
	return &Beam{
		Array:       arr,
		SearchSteps: searchSteps,
		HpbwSteps:   hpbwSteps,
	}
}

// SearchStep 返回搜索步长（度）。
func (b *Beam) SearchStep() float64 {
	steps := b.SearchSteps
	if steps < 1 {
		steps = 720
	}
	return 180 / float64(steps)
}

// HpbwStep 返回 HPBW 扫描步长（度）。
func (b *Beam) HpbwStep() float64 {
	steps := b.HpbwSteps
	if steps < 1 {
		steps = 3600
	}
	return 180 / float64(steps)
}

func nan() float64 {
	var z float64
	return z / z
}
