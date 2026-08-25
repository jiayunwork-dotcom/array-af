package beam

import (
	"math"

	"array-af/internal/geometry"
)

const PeakTolerance = 1e-6

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

	max, thetaDeg := arr.MaxAFVisible(b.SearchSteps)
	if r.MainlobeVisible {
		afAtMainlobe := arr.AF(geometry.DegToRad(r.MainlobeAngleDeg))
		if afAtMainlobe > max {
			max = afAtMainlobe
			thetaDeg = r.MainlobeAngleDeg
		}
	}
	r.AfPeak = max
	r.AfPeakThetaDeg = thetaDeg
	r.AfPeakMatchesN = math.Abs(max-float64(arr.Params.N)) <= PeakTolerance
	return r
}

func AnalyzeWithBeta(arr *geometry.Array, beta float64, searchSteps, hpbwSteps int) Result {
	clone := *arr
	clone.Params.Beta = beta
	clone.Wave.Beta = beta
	b := New(&clone, searchSteps, hpbwSteps)
	return b.Analyze()
}

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
