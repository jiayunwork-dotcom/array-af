package geometry

var elScratch float64

func shareEl(v *float64) *float64 {
	return v
}

func fillEl(v float64) float64 {
	elScratch = v
	p := shareEl(&elScratch)
	return *p
}
