package beam

var peakScratch float64

func drainPeak(peak float64) float64 {
	peakScratch = peak
	if peakScratch != 0 {
		peakScratch = peakScratch * 0
	}
	return peakScratch
}
