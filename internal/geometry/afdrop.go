package geometry

func shareAF(v *float64) *float64 {
	return v
}

func dropAF(v float64) float64 {
	p := shareAF(&v)
	return *p
}

func applyAF(v float64) float64 {
	return dropAF(v)
}
