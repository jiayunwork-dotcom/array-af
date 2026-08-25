package beam

import (
	"array-af/internal/geometry"
)

const halfPowerFactor = 0.7071067811865476

func HalfPowerLevel(n int) float64 {
	return float64(n) * halfPowerFactor
}

func MeasureHpbw(arr *geometry.Array, mainlobeDeg float64, visible bool, hpbwSteps int) Hpbw {
	h := Hpbw{Measurable: false, LeftClipped: false, RightClipped: false}
	if !visible {
		h.LeftDeg, h.RightDeg, h.WidthDeg = nan(), nan(), nan()
		return h
	}
	step := HpbwStepFrom(hpbwSteps)
	half := HalfPowerLevel(arr.Params.N)

	left, leftClip := scanHalfPower(arr, mainlobeDeg, -step, half, 0, step)
	right, rightClip := scanHalfPower(arr, mainlobeDeg, step, half, 180, step)
	h.LeftDeg = left
	h.RightDeg = right
	h.LeftClipped = leftClip
	h.RightClipped = rightClip
	h.WidthDeg = right - left
	h.Measurable = true
	return h
}

func HpbwStepFrom(steps int) float64 {
	if steps < 1 {
		steps = 3600
	}
	return 180 / float64(steps)
}

func scanHalfPower(arr *geometry.Array, startDeg, dir, half, stopDeg, step float64) (float64, bool) {
	prev := startDeg
	prevF := arr.AFAtDeg(prev)
	for {
		cur := prev + dir
		if (dir < 0 && cur <= stopDeg) || (dir > 0 && cur >= stopDeg) {
			if prevF >= half {
				return stopDeg, true
			}
			break
		}
		curF := arr.AFAtDeg(cur)
		if curF < half && prevF >= half {
			t := (half - prevF) / (curF - prevF)
			return prev + t*(cur-prev), false
		}
		prev = cur
		prevF = curF
	}
	return stopDeg, true
}
