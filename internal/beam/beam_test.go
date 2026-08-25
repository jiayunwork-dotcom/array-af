package beam

import (
	"math"
	"testing"

	"array-af/internal/geometry"
)

func mustArray(t *testing.T, p geometry.Params) *geometry.Array {
	t.Helper()
	arr, err := geometry.NewArray(p)
	if err != nil {
		t.Fatalf("NewArray(%+v) unexpected error: %v", p, err)
	}
	return arr
}

func analyze(t *testing.T, p geometry.Params) Result {
	t.Helper()
	arr := mustArray(t, p)
	return New(arr, 720, 3600).Analyze()
}

func TestMainlobeBroadside(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	if !r.MainlobeVisible {
		t.Errorf("broadside mainlobe: expected visible, got invisible")
	}
	if math.Abs(r.MainlobeAngleDeg-90) > 1e-9 {
		t.Errorf("broadside mainlobe: got %vdeg, want 90deg", r.MainlobeAngleDeg)
	}
}

func TestMainlobeEndfire(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: -math.Pi})
	if !r.MainlobeVisible {
		t.Errorf("endfire mainlobe: expected visible, got invisible")
	}
	if math.Abs(r.MainlobeAngleDeg) > 1e-9 {
		t.Errorf("endfire mainlobe: got %vdeg, want 0deg", r.MainlobeAngleDeg)
	}
}

func TestMainlobeInvisibleBeyondKd(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: -3 * math.Pi})
	if r.MainlobeVisible {
		t.Errorf("expected invisible mainlobe for |beta|>kd, got visible at %vdeg", r.MainlobeAngleDeg)
	}
}

func TestNoGratingAtHalfWave(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	if r.HasGrating {
		t.Errorf("d=lambda/2: expected no grating lobes, got %d", len(r.GratingLobes))
	}
}

func TestGratingAboveLambda(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 1.2, Lambda: 1, Beta: 0})
	if !r.HasGrating {
		t.Errorf("d=1.2*lambda: expected grating lobes, got none")
	}
	if len(r.GratingLobes) < 2 {
		t.Errorf("d=1.2*lambda: expected at least 2 grating lobes (m=+1,-1), got %d", len(r.GratingLobes))
	}
	for _, gl := range r.GratingLobes {
		if gl.AngleDeg < 0 || gl.AngleDeg > 180 {
			t.Errorf("grating lobe order %d outside visible range: %vdeg", gl.Order, gl.AngleDeg)
		}
	}
}

func TestGratingAtHalfWaveEndfireTail(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: -math.Pi})
	if !r.HasGrating {
		t.Fatalf("d=lambda/2 endfire: expected grating at theta=180deg (opposite endfire)")
	}
	found := false
	for _, gl := range r.GratingLobes {
		if gl.Order == -1 && math.Abs(gl.AngleDeg-180) < 1e-6 {
			found = true
		}
	}
	if !found {
		t.Errorf("d=lambda/2 endfire: grating lobes = %+v, want order -1 at 180deg", r.GratingLobes)
	}
}

func TestFirstNullsBroadside8(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	wantLeft := math.Acos(0.25) * 180 / math.Pi
	wantRight := math.Acos(-0.25) * 180 / math.Pi
	if !r.Nulls.LeftValid || math.Abs(r.Nulls.LeftDeg-wantLeft) > 1e-6 {
		t.Errorf("first null left: got %v (valid=%v), want %v", r.Nulls.LeftDeg, r.Nulls.LeftValid, wantLeft)
	}
	if !r.Nulls.RightValid || math.Abs(r.Nulls.RightDeg-wantRight) > 1e-6 {
		t.Errorf("first null right: got %v (valid=%v), want %v", r.Nulls.RightDeg, r.Nulls.RightValid, wantRight)
	}
}

func TestHpbwNarrowsWithMoreElements(t *testing.T) {
	r4 := analyze(t, geometry.Params{N: 4, D: 0.5, Lambda: 1, Beta: 0})
	r8 := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	if !r4.Hpbw.Measurable || !r8.Hpbw.Measurable {
		t.Fatalf("expected measurable hpbw for N=4/8, got N=4:%v N=8:%v", r4.Hpbw.Measurable, r8.Hpbw.Measurable)
	}
	if r8.Hpbw.WidthDeg >= r4.Hpbw.WidthDeg {
		t.Errorf("hpbw(N=8)=%vdeg should be narrower than hpbw(N=4)=%vdeg", r8.Hpbw.WidthDeg, r4.Hpbw.WidthDeg)
	}
}

func TestHpbwWidensTowardEndfire(t *testing.T) {
	rSide := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	rEnd := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: -math.Pi})
	if !rSide.Hpbw.Measurable || !rEnd.Hpbw.Measurable {
		t.Fatalf("expected measurable hpbw, got side:%v end:%v", rSide.Hpbw.Measurable, rEnd.Hpbw.Measurable)
	}
	if rEnd.Hpbw.WidthDeg <= rSide.Hpbw.WidthDeg {
		t.Errorf("endfire hpbw=%vdeg should be wider than broadside hpbw=%vdeg", rEnd.Hpbw.WidthDeg, rSide.Hpbw.WidthDeg)
	}
}

func TestHpbwBroadside8Approx(t *testing.T) {
	r := analyze(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	if r.Hpbw.WidthDeg < 12.0 || r.Hpbw.WidthDeg > 13.5 {
		t.Errorf("broadside-8 hpbw: got %vdeg, want in [12, 13.5]", r.Hpbw.WidthDeg)
	}
}

func TestPeakMatchesN(t *testing.T) {
	configs := []geometry.Params{
		{N: 8, D: 0.5, Lambda: 1, Beta: 0},
		{N: 16, D: 0.25, Lambda: 1, Beta: -1.2},
		{N: 5, D: 1.2, Lambda: 1, Beta: 0},
	}
	for _, p := range configs {
		r := analyze(t, p)
		want := float64(p.N)
		if !r.AfPeakMatchesN {
			t.Errorf("af_peak=%v does not match N=%v (config %+v)", r.AfPeak, p.N, p)
		}
		if math.Abs(r.AfPeak-want) > 1e-6 {
			t.Errorf("af_peak: got %v, want %v (config %+v)", r.AfPeak, want, p)
		}
	}
}

func TestDirectivityHalfWaveBroadside(t *testing.T) {
	arr := mustArray(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: 0})
	d := ApproxDirectivity(arr)
	if !d.Valid {
		t.Errorf("half-wave broadside: expected valid directivity, got reason %q", d.Reason)
	}
	if math.Abs(d.Approx-8) > 1e-9 {
		t.Errorf("directivity: got %v, want 8", d.Approx)
	}
}

func TestDirectivityNotValidOffBroadside(t *testing.T) {
	arr := mustArray(t, geometry.Params{N: 8, D: 0.5, Lambda: 1, Beta: -1})
	d := ApproxDirectivity(arr)
	if d.Valid {
		t.Errorf("non-broadside: expected invalid directivity, got %v", d.Approx)
	}
}

func TestGratingAnglesSatisfyPsi(t *testing.T) {
	arr := mustArray(t, geometry.Params{N: 8, D: 1.2, Lambda: 1, Beta: 0})
	lobes := GratingAngles(arr.Wave)
	for _, gl := range lobes {
		psi := arr.Wave.PsiDeg(gl.AngleDeg)
		if math.Abs(psi-gl.Psi) > 1e-6 {
			t.Errorf("grating order %d: psi at angle %v = %v, want %v", gl.Order, gl.AngleDeg, psi, gl.Psi)
		}
		if math.Abs(math.Abs(gl.Psi)-2*math.Pi*math.Abs(float64(gl.Order))) > 1e-9 {
			t.Errorf("grating order %d: psi=%v not equal 2pi*m", gl.Order, gl.Psi)
		}
	}
}

func TestPolarRadiusNormalized(t *testing.T) {
	if got := PolarRadius(8, 8); math.Abs(got-1) > 1e-12 {
		t.Errorf("PolarRadius(8,8): got %v, want 1", got)
	}
	if got := PolarRadius(0, 8); math.Abs(got) > 1e-12 {
		t.Errorf("PolarRadius(0,8): got %v, want 0", got)
	}
	if got := PolarRadius(12, 8); math.Abs(got-1) > 1e-12 {
		t.Errorf("PolarRadius(12,8): got %v, want capped 1", got)
	}
}
