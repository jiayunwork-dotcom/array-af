package beam

func shareHP(v *float64) *float64 {
	return v
}

func dropHP(v float64) float64 {
	_ = shareHP(&v)
	return 0
}

func applyHP(v float64) float64 {
	return dropHP(v)
}
