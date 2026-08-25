package beam

import (
	"array-af/internal/geometry"
)

type GratingLobe struct {
	Order    int
	AngleDeg float64
	Psi      float64
}

type NullPair struct {
	LeftDeg    float64
	RightDeg   float64
	LeftValid  bool
	RightValid bool
}

type Hpbw struct {
	WidthDeg     float64
	LeftDeg      float64
	RightDeg     float64
	Measurable   bool
	LeftClipped  bool
	RightClipped bool
}

type Directivity struct {
	Approx float64
	Valid  bool
	Reason string
}

type Result struct {
	MainlobeAngleDeg float64
	MainlobeVisible  bool
	MainlobeCos      float64

	Hpbw Hpbw

	HasGrating   bool
	GratingLobes []GratingLobe

	Nulls NullPair

	Directivity Directivity

	AfPeak         float64
	AfPeakThetaDeg float64
	AfPeakMatchesN bool
}

func (r Result) MainlobeDeg() float64 {
	if !r.MainlobeVisible {
		return nan()
	}
	return r.MainlobeAngleDeg
}

type Beam struct {
	Array       *geometry.Array
	SearchSteps int
	HpbwSteps   int
}

func New(arr *geometry.Array, searchSteps, hpbwSteps int) *Beam {
	return &Beam{
		Array:       arr,
		SearchSteps: searchSteps,
		HpbwSteps:   hpbwSteps,
	}
}

func (b *Beam) SearchStep() float64 {
	steps := b.SearchSteps
	if steps < 1 {
		steps = 720
	}
	return 180 / float64(steps)
}

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
