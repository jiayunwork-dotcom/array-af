package geometry

import (
	"math"
	"testing"
)

func TestHammingEndpoints(t *testing.T) {
	w0 := Hamming(8, 0)
	w7 := Hamming(8, 7)
	if w0 > 0.2 || w7 > 0.2 {
		t.Fatalf("endpoints too high w0=%v w7=%v", w0, w7)
	}
	center := Hamming(8, 4)
	if center <= w0 {
		t.Fatalf("center=%v w0=%v", center, w0)
	}
}

func TestHanningSymmetry(t *testing.T) {
	n := 9
	for i := 0; i < n/2; i++ {
		a := Hanning(n, i)
		b := Hanning(n, n-1-i)
		if math.Abs(a-b) > 1e-12 {
			t.Fatalf("i=%d a=%v b=%v", i, a, b)
		}
	}
}

func TestBuildTaperUniform(t *testing.T) {
	w := BuildTaper(TaperUniform, 5, 0)
	for i, v := range w {
		if v != 1 {
			t.Fatalf("i=%d v=%v", i, v)
		}
	}
}

func TestTaylorWeightsBounded(t *testing.T) {
	w := BuildTaper(TaperTaylor, 16, -30)
	max := 0.0
	for _, v := range w {
		a := v
		if a < 0 {
			a = -a
		}
		if a > max {
			max = a
		}
	}
	if max <= 0 {
		t.Fatalf("max=%v", max)
	}
}

func TestTaperEfficiency(t *testing.T) {
	uniform := BuildTaper(TaperUniform, 8, 0)
	ham := BuildTaper(TaperHamming, 8, 0)
	effU := TaperEfficiency(uniform)
	effH := TaperEfficiency(ham)
	if effH >= effU {
		t.Fatalf("effU=%v effH=%v", effU, effH)
	}
}

func TestNormalizeWeights(t *testing.T) {
	w := []float64{1, 2, 3, 4}
	nw := NormalizeWeights(w)
	sum := 0.0
	for _, v := range nw {
		sum += v
	}
	if math.Abs(sum-float64(len(w))) > 1e-9 {
		t.Fatalf("sum=%v", sum)
	}
}

func TestWeightedArrayFactorBroadside(t *testing.T) {
	p := Params{N: 8, D: 0.5, Lambda: 1, Beta: 0, Element: ElementIso}
	w := BuildTaper(TaperHamming, p.N, 0)
	psi := 0.0
	af := WeightedArrayFactor(p.N, psi, w)
	if af <= 0 {
		t.Fatalf("af=%v", af)
	}
}

func TestNewWeightedArray(t *testing.T) {
	p := Params{N: 8, D: 0.5, Lambda: 1, Beta: 0, Element: ElementIso}
	wa, err := NewWeightedArray(p, TaperHamming, -25)
	if err != nil {
		t.Fatal(err)
	}
	peak, _ := wa.MaxWAFVisible(180)
	if peak <= 0 {
		t.Fatal("zero peak")
	}
	red := wa.PeakReduction()
	if red <= 0 {
		t.Fatalf("reduction=%v", red)
	}
}
