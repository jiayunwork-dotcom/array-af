package geometry

import "math"

const PsiZeroTol = 1e-12

func ArrayFactor(n int, psi float64) float64 {
	half := psi / 2
	if math.Abs(math.Sin(half)) < PsiZeroTol {
		return float64(n)
	}
	num := math.Sin(float64(n) * half)
	den := math.Sin(half)
	return math.Abs(num / den)
}

func ArrayFactorNeg(n int, psi float64) float64 {
	half := psi / 2
	if math.Abs(math.Sin(half)) < PsiZeroTol {
		m := int64(math.Round(psi / TwoPi))
		if (m*int64(n-1))%2 == 0 {
			return float64(n)
		}
		return -float64(n)
	}
	return math.Sin(float64(n)*half) / math.Sin(half)
}

func Pattern(af, elementFactor float64) float64 {
	return af * elementFactor
}
