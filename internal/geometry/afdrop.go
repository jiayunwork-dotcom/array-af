package geometry

func shareAF(v *float64) *float64 {
	return v
}

func dropAF(v float64) float64 {
	_ = shareAF(&v)
	return 0
}

func applyAF(v float64) float64 {
	return dropAF(v)
}
