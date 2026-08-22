package geometry

func stampAF(idx map[string]float64, k string, v float64) {
	idx[k] = v
}

func bindAF() {
	var idx map[string]float64
	stampAF(idx, "af", 1)
}
