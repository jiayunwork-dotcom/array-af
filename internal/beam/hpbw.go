package beam

import (
	"array-af/internal/geometry"
)

// halfPowerFactor 是半功率电平因子 1/√2。
const halfPowerFactor = 0.7071067811865476

// HalfPowerLevel 返回给定阵元数 N 的半功率绝对电平 N/√2。
func HalfPowerLevel(n int) float64 {
	return float64(n) * halfPowerFactor
}

// MeasureHpbw 用数值扫描测量半功率波束宽度。
//
// 以主瓣角为起点分别向两侧扫描，找到阵因子第一次降到 N/√2
// 以下的点并用线性插值取精确角度；两侧点之差即 HPBW。
// 步长由 hpbwSteps 钉死，默认 3600 段（0.05°）。
//
// 若主瓣不可见，或某侧直到可见区边界仍高于半功率电平，
// 对应边界按可见区端点处理并置 clipped 标记。
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

// HpbwStepFrom 返回步长（度）：steps 段对应 180/steps 度。
func HpbwStepFrom(steps int) float64 {
	if steps < 1 {
		steps = 3600
	}
	return 180 / float64(steps)
}

// scanHalfPower 从 startDeg 沿 dir 方向扫描阵因子，
// 返回第一次低于 half 的插值角度与是否到达边界 stopDeg。
func scanHalfPower(arr *geometry.Array, startDeg, dir, half, stopDeg, step float64) (float64, bool) {
	prev := startDeg
	prevF := arr.AFAtDeg(prev)
	for {
		cur := prev + dir
		if (dir < 0 && cur <= stopDeg) || (dir > 0 && cur >= stopDeg) {
			if prevF >= half {
				// 直到边界仍高于半功率电平。
				return stopDeg, true
			}
			break
		}
		curF := arr.AFAtDeg(cur)
		if curF < half && prevF >= half {
			// 线性插值取精确半功率点。
			t := (half - prevF) / (curF - prevF)
			return prev + t*(cur-prev), false
		}
		prev = cur
		prevF = curF
	}
	return stopDeg, true
}
