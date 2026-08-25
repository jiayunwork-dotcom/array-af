package geometry

import "math"

type TaperKind int

const (
	TaperUniform TaperKind = iota
	TaperHamming
	TaperHanning
	TaperTaylor
)

func Hamming(n, i int) float64 {
	if n <= 1 {
		return 1
	}
	return 0.54 - 0.46*math.Cos(TwoPi*float64(i)/float64(n-1))
}

func Hanning(n, i int) float64 {
	if n <= 1 {
		return 1
	}
	return 0.5 * (1 - math.Cos(TwoPi*float64(i)/float64(n-1)))
}

func Taylor(n, i int, sidelobeLevelDB float64) float64 {
	if n <= 1 {
		return 1
	}
	A := math.Cosh(math.Log(math.Pow(10, sidelobeLevelDB/20)) / math.Pi)
	s2 := A * A
	m := int(math.Ceil(2*A*A-0.5)) + 1
	if m > n/2 {
		m = n / 2
	}
	if m < 1 {
		m = 1
	}
	w := make([]float64, n)
	for j := 0; j < n; j++ {
		w[j] = 1
	}
	for l := 1; l <= m; l++ {
		num := 1.0
		for p := 1; p <= m; p++ {
			if p != l {
				num *= 1 - float64(l*l)/s2/(float64(p)*float64(p)-0.25)
			}
		}
		den := float64(l) * float64(l)
		for p := 1; p <= m; p++ {
			if p != l {
				den *= 1 - float64(l*l)/float64(p*p)
			}
		}
		Fm := num / den
		for j := 0; j < n; j++ {
			w[j] += 2 * Fm * math.Cos(TwoPi*float64(l)*(float64(j)-float64(n-1)/2)/float64(n))
		}
	}
	max := 0.0
	for _, v := range w {
		if v > max {
			max = v
		}
	}
	if max <= 0 {
		return 1
	}
	return w[i] / max
}

func BuildTaper(kind TaperKind, n int, sidelobeDB float64) []float64 {
	if n < 1 {
		return nil
	}
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		switch kind {
		case TaperHamming:
			out[i] = Hamming(n, i)
		case TaperHanning:
			out[i] = Hanning(n, i)
		case TaperTaylor:
			out[i] = Taylor(n, i, sidelobeDB)
		default:
			out[i] = 1
		}
	}
	return out
}

func NormalizeWeights(w []float64) []float64 {
	if len(w) == 0 {
		return w
	}
	sum := 0.0
	for _, v := range w {
		sum += v
	}
	if sum <= 0 {
		return w
	}
	out := make([]float64, len(w))
	for i, v := range w {
		out[i] = v / sum * float64(len(w))
	}
	return out
}

func TaperEfficiency(w []float64) float64 {
	if len(w) == 0 {
		return 0
	}
	sum := 0.0
	sumSq := 0.0
	for _, v := range w {
		sum += v
		sumSq += v * v
	}
	if sumSq <= 0 {
		return 0
	}
	return sum * sum / (float64(len(w)) * sumSq)
}

func TaperPeakRatio(w []float64) float64 {
	if len(w) == 0 {
		return 0
	}
	max := 0.0
	sum := 0.0
	for _, v := range w {
		if v > max {
			max = v
		}
		sum += v
	}
	if sum <= 0 {
		return 0
	}
	mean := sum / float64(len(w))
	if mean <= 0 {
		return 0
	}
	return max / mean
}
