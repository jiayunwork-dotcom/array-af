package scan

import "math"

type Row struct {
	Beta            float64
	BetaDeg         float64
	MainlobeDeg     float64
	MainlobeVisible bool
	HpbwDeg         float64
	HasGrating      bool
}

func (r Row) MainlobeValue() float64 {
	if !r.MainlobeVisible {
		return math.NaN()
	}
	return r.MainlobeDeg
}

func (r Row) HpbwValue() float64 {
	return r.HpbwDeg
}

type Summary struct {
	MainlobeStartDeg           float64
	MainlobeEndDeg             float64
	MainlobeVisibleAtStart     bool
	MainlobeVisibleAtEnd       bool
	MainlobeMovesTowardEndfire bool
	HpbwWidened                bool
	GratingAppears             bool
	GratingPresent             bool
}
