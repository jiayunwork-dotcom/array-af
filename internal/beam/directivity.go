package beam

import (
	"math"

	"array-af/internal/geometry"
)

// HalfWaveSpacingTol 是判断 d≈λ/2 的相对容差。
const HalfWaveSpacingTol = 0.01

// ApproxDirectivity 计算方向性近似 D≈N。
//
// 该近似仅在半波间距（d≈λ/2）且侧射（β=0）时成立：
// 此时均匀直线阵方向性约为 N，不做与 N 无关的常数修正。
// 其他配置返回 Valid=false 并说明原因。
func ApproxDirectivity(arr *geometry.Array) Directivity {
	n := float64(arr.Params.N)
	if !isHalfWave(arr) {
		return Directivity{
			Approx: 0,
			Valid:  false,
			Reason: "directionality N requires half-wave spacing",
		}
	}
	if !arr.Wave.IsBroadside() {
		return Directivity{
			Approx: 0,
			Valid:  false,
			Reason: "directionality N requires broadside beam",
		}
	}
	return Directivity{
		Approx: n,
		Valid:  true,
		Reason: "half-wave broadside uniform array",
	}
}

// isHalfWave 判断 d/λ 是否在半波间距容差内。
func isHalfWave(arr *geometry.Array) bool {
	ratio := arr.SpacingRatio()
	return math.Abs(ratio-0.5) <= HalfWaveSpacingTol*0.5
}

// DirectivityFromN 返回给定阵元数对应的近似方向性 N。
// 仅用于已知近似成立的场景。
func DirectivityFromN(n int) float64 {
	return float64(n)
}
