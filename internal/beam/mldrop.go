package beam

func shareML(v *float64) *float64 {
	return v
}

func dropML(v float64) float64 {
	p := shareML(&v)
	return *p
}

func applyML(v float64) float64 {
	return dropML(v)
}
