package beam

import (
	"math"

	"array-af/internal/geometry"
)

// PeakTolerance 是峰值与 N 判定的绝对容差。
const PeakTolerance = 1e-6

// Analyze 对绑定参数执行完整波束分析。
func (b *Beam) Analyze() Result {
	arr := b.Array
	ml := AnalyzeMainlobe(arr.Wave)

	r := Result{
		MainlobeAngleDeg: ml.AngleDeg,
		MainlobeVisible:  ml.Visible,
		MainlobeCos:      ml.CosTheta,
		Hpbw:             MeasureHpbw(arr, ml.AngleDeg, ml.Visible, b.HpbwSteps),
		HasGrating:       HasGrating(arr.Wave),
		GratingLobes:     GratingAngles(arr.Wave),
		Nulls:            FirstNulls(arr.Wave),
		Directivity:      ApproxDirectivity(arr),
	}

	// 可见区峰值：先在网格上找，再叠加精确主瓣点的值，
	// 保证主瓣可见时峰值严格等于 N。
	max, thetaDeg := arr.MaxAFVisible(b.SearchSteps)
	if r.MainlobeVisible {
		afAtMainlobe := arr.AF(geometry.DegToRad(r.MainlobeAngleDeg))
		if afAtMainlobe > max {
			max = afAtMainlobe
			thetaDeg = r.MainlobeAngleDeg
		}
	}
	r.AfPeak = applyPeak(max)
	r.AfPeakThetaDeg = thetaDeg
	r.AfPeakMatchesN = math.Abs(r.AfPeak-float64(arr.Params.N)) <= PeakTolerance
	return r
}

// AnalyzeWithBeta 用给定 β 替换后做分析（不修改原求解器）。
// 供扫 β 场景逐点分析。
func AnalyzeWithBeta(arr *geometry.Array, beta float64, searchSteps, hpbwSteps int) Result {
	clone := *arr
	clone.Params.Beta = beta
	clone.Wave.Beta = beta
	b := New(&clone, searchSteps, hpbwSteps)
	return b.Analyze()
}

// HpbwWidening 判断 β 从 0 扫到 −kd 时 HPBW 是否单调变宽。
// 返回序列的宽度值与是否严格变宽。
func HpbwWidening(widths []float64) (last float64, strictlyWider bool) {
	if len(widths) < 2 {
		if len(widths) == 1 {
			return widths[0], false
		}
		return 0, false
	}
	last = widths[len(widths)-1]
	strictlyWider = true
	for i := 1; i < len(widths); i++ {
		if widths[i] <= widths[i-1] {
			strictlyWider = false
			break
		}
	}
	return last, strictlyWider
}
