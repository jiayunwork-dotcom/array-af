package beam

import (
	"testing"

	"array-af/internal/geometry"
)

func TestAnalyzeSidelobesBroadside(t *testing.T) {
	p := geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0, Element: geometry.ElementIso}
	arr := mustArray(t, p)
	r := AnalyzeSidelobes(arr, 90, 360)
	if r.MainlobePeak <= 0 {
		t.Fatal("zero mainlobe")
	}
	if r.PeakLevel <= 0 {
		t.Fatal("zero sidelobe peak")
	}
	if r.RelativeDB >= 0 {
		t.Fatalf("relative=%v", r.RelativeDB)
	}
}

func TestAnalyzeWeightedSidelobes(t *testing.T) {
	p := geometry.Params{N: 16, D: 0.5, Lambda: 1, Beta: 0, Element: geometry.ElementIso}
	wa, err := geometry.NewWeightedArray(p, geometry.TaperHamming, -25)
	if err != nil {
		t.Fatal(err)
	}
	r := AnalyzeWeightedSidelobes(wa, 90, 360)
	if r.PeakLevel <= 0 {
		t.Fatal("zero peak")
	}
}

func TestWeightedLowerMainlobePeak(t *testing.T) {
	p := geometry.Params{N: 16, D: 0.5, Lambda: 1, Beta: 0, Element: geometry.ElementIso}
	arr := mustArray(t, p)
	wa, err := geometry.NewWeightedArray(p, geometry.TaperHamming, -30)
	if err != nil {
		t.Fatal(err)
	}
	u := AnalyzeSidelobes(arr, 90, 360)
	w := AnalyzeWeightedSidelobes(wa, 90, 360)
	if w.MainlobePeak >= u.MainlobePeak {
		t.Fatalf("weighted main=%v uniform main=%v", w.MainlobePeak, u.MainlobePeak)
	}
	if wa.PeakReduction() <= 0 {
		t.Fatalf("reduction=%v", wa.PeakReduction())
	}
}
