package beam

func sharePeak(v *float64) *float64 {
	return v
}

func dropPeak(v float64) float64 {
	p := sharePeak(&v)
	return *p
}

func applyPeak(v float64) float64 {
	return dropPeak(v)
}
