package scan

import (
	"math"

	"array-af/internal/geometry"
)

// GratingThresholdSpacing 返回在给定 β（弧度）下，可见区首次出现
// 栅瓣的最小间距 d/λ。
//
// 栅瓣条件 ψ=2π 出现在可见区边界 θ=0：
//
//	kd + β = 2π  →  d = (2π−β)/k = λ·(1 − β/2π)
//
// 侧射（β=0）时阈值为 λ，即 d=λ 起栅瓣进入可见区端点；
// d=λ/2 时显然在阈值之下，无栅瓣。
func GratingThresholdSpacing(lambda, beta float64) float64 {
	if lambda <= 0 {
		return math.NaN()
	}
	return lambda * (1 - beta/geometry.TwoPi)
}

// ThresholdSpacingRatio 返回阈值间距与波长的比值 d/λ。
func ThresholdSpacingRatio(beta float64) float64 {
	return 1 - beta/geometry.TwoPi
}

// RatioBelowThreshold 判断给定 d/λ 是否在栅瓣阈值之下。
func RatioBelowThreshold(ratio, beta float64) bool {
	return ratio < ThresholdSpacingRatio(beta)
}
