package beam

func shareHP(v *float64) *float64 {
	return v
}

func dropHP(v float64) float64 {
	p := shareHP(&v)
	return *p
}

func applyHP(v float64) float64 {
	return dropHP(v)
}
