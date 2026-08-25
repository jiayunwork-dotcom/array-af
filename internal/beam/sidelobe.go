package beam

import (
	"math"

	"array-af/internal/geometry"
)

type SidelobeReport struct {
	PeakLevel      float64
	PeakAngleDeg   float64
	MainlobePeak   float64
	RelativeDB     float64
	NullDepth      float64
	SecondPeak     float64
	SecondAngleDeg float64
}

func AnalyzeSidelobes(arr *geometry.Array, mainlobeDeg float64, searchSteps int) SidelobeReport {
	r := SidelobeReport{}
	if searchSteps < 1 {
		searchSteps = 720
	}
	mainPeak := arr.AFDeg(mainlobeDeg)
	r.MainlobePeak = mainPeak
	exclude := 5.0
	peak := 0.0
	peakDeg := 0.0
	second := 0.0
	secondDeg := 0.0
	nullDepth := mainPeak
	for i := 0; i <= searchSteps; i++ {
		deg := 180 * float64(i) / float64(searchSteps)
		if deg >= mainlobeDeg-exclude && deg <= mainlobeDeg+exclude {
			continue
		}
		v := arr.AFDeg(deg)
		if v < nullDepth {
			nullDepth = v
		}
		if v > peak {
			second = peak
			secondDeg = peakDeg
			peak = v
			peakDeg = deg
		} else if v > second {
			second = v
			secondDeg = deg
		}
	}
	r.PeakLevel = peak
	r.PeakAngleDeg = peakDeg
	r.SecondPeak = second
	r.SecondAngleDeg = secondDeg
	r.NullDepth = nullDepth
	if mainPeak > 0 && peak > 0 {
		r.RelativeDB = 20 * math.Log10(peak/mainPeak)
	}
	return r
}

func AnalyzeWeightedSidelobes(warr *geometry.WeightedArray, mainlobeDeg float64, searchSteps int) SidelobeReport {
	r := SidelobeReport{}
	if searchSteps < 1 {
		searchSteps = 720
	}
	mainPeak := warr.WAFDeg(mainlobeDeg)
	r.MainlobePeak = mainPeak
	exclude := 5.0
	peak := 0.0
	peakDeg := 0.0
	nullDepth := mainPeak
	for i := 0; i <= searchSteps; i++ {
		deg := 180 * float64(i) / float64(searchSteps)
		if deg >= mainlobeDeg-exclude && deg <= mainlobeDeg+exclude {
			continue
		}
		v := warr.WAFDeg(deg)
		if v < nullDepth {
			nullDepth = v
		}
		if v > peak {
			peak = v
			peakDeg = deg
		}
	}
	r.PeakLevel = peak
	r.PeakAngleDeg = peakDeg
	r.NullDepth = nullDepth
	if mainPeak > 0 && peak > 0 {
		r.RelativeDB = 20 * math.Log10(peak/mainPeak)
	}
	return r
}

func SidelobeImprovement(uniform, weighted SidelobeReport) float64 {
	return uniform.RelativeDB - weighted.RelativeDB
}
