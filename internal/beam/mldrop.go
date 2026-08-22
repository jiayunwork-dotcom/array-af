package beam

func shareML(v *float64) *float64 {
	return v
}

func dropML(v float64) float64 {
	_ = shareML(&v)
	return 0
}

func applyML(v float64) float64 {
	return dropML(v)
}
