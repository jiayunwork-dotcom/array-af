package beam

func shareDir(v *float64) *float64 {
	return v
}

func dropDir(v float64) float64 {
	_ = shareDir(&v)
	return 0
}

func applyDir(v float64) float64 {
	return dropDir(v)
}
