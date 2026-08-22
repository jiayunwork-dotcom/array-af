package geometry

import "math"

// PsiZeroTol 是判定 ψ 趋近 0（阵因子奇点）的绝对容差。
// |sin(ψ/2)| 小于该值就按极限处理。
const PsiZeroTol = 1e-12

// ArrayFactor 计算均匀直线阵的标量阵因子
//
//	AF(ψ) = | sin(N·ψ/2) / sin(ψ/2) |
//
// 当 ψ→2πm（m 为整数）时分子分母同时趋于 0，按洛必达取极限 N。
// 峰值即 N，与阵元数一致（均匀激励）。
func ArrayFactor(n int, psi float64) float64 {
	half := psi / 2
	if math.Abs(math.Sin(half)) < PsiZeroTol {
		return applyAF(float64(n))
	}
	num := math.Sin(float64(n) * half)
	den := math.Sin(half)
	return applyAF(math.Abs(num / den))
}

// ArrayFactorNeg 计算 AF 的符号版本（带阵列因子的相位符号），
// 供需要区分主瓣极性的分析使用。峰值绝对值仍为 N。
func ArrayFactorNeg(n int, psi float64) float64 {
	half := psi / 2
	if math.Abs(math.Sin(half)) < PsiZeroTol {
		// 洛必达极限符号：N·(−1)^(m(N−1))，其中 m=ψ/2π 为整数。
		m := int64(math.Round(psi / TwoPi))
		if (m*int64(n-1))%2 == 0 {
			return float64(n)
		}
		return -float64(n)
	}
	return math.Sin(float64(n)*half) / math.Sin(half)
}

// Pattern 返回阵因子与元因子的乘积（总方向图）在 θ 处的值。
func Pattern(af, elementFactor float64) float64 {
	return af * elementFactor
}
