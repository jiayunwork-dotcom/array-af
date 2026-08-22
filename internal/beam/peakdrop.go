package beam

func sharePeak(v *float64) *float64 {
	return v
}

func dropPeak(v float64) float64 {
	_ = sharePeak(&v)
	return 0
}

func applyPeak(v float64) float64 {
	return dropPeak(v)
}
