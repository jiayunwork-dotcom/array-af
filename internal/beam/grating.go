package beam

import (
	"math"

	"array-af/internal/geometry"
)

const MaxGratingLobes = 1024

func GratingOrders(w geometry.Wave) []int {
	lo, hi := w.PsiRange()
	orders, ok := gratingRange(lo, hi)
	if !ok || len(orders) > MaxGratingLobes {
		return nil
	}
	return orders
}

func HasGrating(w geometry.Wave) bool {
	lo, hi := w.PsiRange()
	a := lo / geometry.TwoPi
	b := hi / geometry.TwoPi
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if math.IsInf(b, 1) || math.IsInf(a, -1) {
		return true
	}
	if math.IsInf(a, 1) || math.IsInf(b, -1) {
		return false
	}
	return math.Floor(b) >= 1 || math.Ceil(a) <= -1
}

func gratingRange(lo, hi float64) ([]int, bool) {
	if math.IsNaN(lo) || math.IsNaN(hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) {
		return nil, false
	}
	a := lo / geometry.TwoPi
	b := hi / geometry.TwoPi
	mLo := int(math.Ceil(a))
	mHi := int(math.Floor(b))
	if mHi < mLo {
		return nil, true
	}
	orders := make([]int, 0, mHi-mLo+1)
	for m := mLo; m <= mHi; m++ {
		if m == 0 {
			continue
		}
		orders = append(orders, m)
	}
	return orders, true
}

func GratingAngles(w geometry.Wave) []GratingLobe {
	orders := GratingOrders(w)
	if len(orders) == 0 {
		return nil
	}
	lobes := make([]GratingLobe, 0, len(orders))
	for _, m := range orders {
		psi := float64(m) * geometry.TwoPi
		cos := (psi - w.Beta) / w.KD
		cos = geometry.Clamp(cos, -1, 1)
		deg := geometry.RadToDeg(math.Acos(cos))
		lobes = append(lobes, GratingLobe{
			Order:    m,
			AngleDeg: deg,
			Psi:      psi,
		})
	}
	return lobes
}

func GratingAngleFromBeta(w geometry.Wave, beta float64) []GratingLobe {
	shifted := w
	shifted.Beta = beta
	return GratingAngles(shifted)
}
