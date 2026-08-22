package geometry

import (
	"math"
	"testing"
)

func mustArray(t *testing.T, p Params) *Array {
	t.Helper()
	arr, err := NewArray(p)
	if err != nil {
		t.Fatalf("NewArray(%+v) unexpected error: %v", p, err)
	}
	return arr
}

func TestValidateRejectsInvalidParams(t *testing.T) {
	cases := []struct {
		name  string
		param Params
		want  string
	}{
		{"N=1", Params{N: 1, D: 0.5, Lambda: 1}, "N"},
		{"N=0", Params{N: 0, D: 0.5, Lambda: 1}, "N"},
		{"d=0", Params{N: 8, D: 0, Lambda: 1}, "d"},
		{"d<0", Params{N: 8, D: -0.5, Lambda: 1}, "d"},
		{"lambda=0", Params{N: 8, D: 0.5, Lambda: 0}, "lambda"},
		{"lambda<0", Params{N: 8, D: 0.5, Lambda: -2}, "lambda"},
	}
	for _, c := range cases {
		err := Validate(c.param)
		if err == nil {
			t.Errorf("Validate(%s): expected error, got nil", c.name)
			continue
		}
		ve, ok := err.(*ValidationError)
		if !ok {
			t.Errorf("Validate(%s): expected *ValidationError, got %T: %v", c.name, err, err)
			continue
		}
		if !containsField(ve, c.want) {
			t.Errorf("Validate(%s): expected field %q in %v", c.name, c.want, ve.Fields)
		}
	}
}

func containsField(ve *ValidationError, field string) bool {
	for _, f := range ve.Fields {
		if f.Field == field {
			return true
		}
	}
	return false
}

func TestNewArrayRejectsBadInput(t *testing.T) {
	bad := []Params{
		{N: 1, D: 0.5, Lambda: 1},
		{N: 8, D: 0, Lambda: 1},
		{N: 8, D: 0.5, Lambda: -1},
	}
	for _, p := range bad {
		if _, err := NewArray(p); err == nil {
			t.Errorf("NewArray(%+v): expected error, got nil", p)
		}
	}
}

func TestWavePsiDefinition(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1, Beta: math.Pi / 4})
	k := 2 * math.Pi / arr.Params.Lambda
	kd := k * arr.Params.D
	want := kd*math.Cos(math.Pi/3) + math.Pi/4
	got := arr.Wave.Psi(math.Pi / 3)
	if math.Abs(got-want) > 1e-12 {
		t.Errorf("Psi(60deg): got %v, want %v", got, want)
	}
}

func TestPsiRangeMatchesCosExtremes(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1, Beta: 0.3})
	lo, hi := arr.Wave.PsiRange()
	wantLo := arr.Wave.Psi(math.Pi)
	wantHi := arr.Wave.Psi(0)
	if math.Abs(lo-wantLo) > 1e-12 || math.Abs(hi-wantHi) > 1e-12 {
		t.Errorf("PsiRange: got [%v,%v], want [%v,%v]", lo, hi, wantLo, wantHi)
	}
}

func TestAFPeakEqualsN(t *testing.T) {
	configs := []Params{
		{N: 4, D: 0.5, Lambda: 1},
		{N: 8, D: 0.5, Lambda: 1},
		{N: 16, D: 1.2, Lambda: 1},
		{N: 5, D: 0.75, Lambda: 1, Beta: -1},
	}
	for _, p := range configs {
		arr := mustArray(t, p)
		peak, _ := arr.MaxAFVisible(7200)
		want := float64(p.N)
		if math.Abs(peak-want) > 1e-6 {
			t.Errorf("AF peak (N=%d,d=%v,lambda=%v,beta=%v): got %v, want %v", p.N, p.D, p.Lambda, p.Beta, peak, want)
		}
	}
}

func TestAFLimitAtPsiZero(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1})
	cos := Clamp(arr.Wave.MainlobePsi(), -1, 1)
	thetaMain := math.Acos(cos)
	if !arr.Wave.MainlobeVisible() {
		t.Fatalf("expected visible mainlobe for broadside-8")
	}
	got := arr.AF(thetaMain)
	want := 8.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("AF at psi=0: got %v, want %v", got, want)
	}
	got2 := ArrayFactor(8, 1e-13)
	if math.Abs(got2-want) > 1e-9 {
		t.Errorf("ArrayFactor(8, 1e-13): got %v, want %v", got2, want)
	}
}

func TestArrayFactorZeroWhenPsiIsNull(t *testing.T) {
	psi := 2 * math.Pi / 8
	got := ArrayFactor(8, psi)
	if got > 1e-9 {
		t.Errorf("ArrayFactor(8, psi=2pi/8): got %v, want near 0", got)
	}
}

func TestElementFactorIsoIsOne(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1, Element: ElementIso})
	for _, deg := range []float64{0, 45, 90, 180} {
		if got := arr.ElementFactorDeg(deg); got != 1 {
			t.Errorf("iso element at %vdeg: got %v, want 1", deg, got)
		}
	}
}

func TestElementFactorDipoleIsSine(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1, Element: ElementDipole})
	if got := arr.ElementFactorDeg(90); math.Abs(got-1) > 1e-12 {
		t.Errorf("dipole at 90deg: got %v, want 1", got)
	}
	if got := arr.ElementFactorDeg(0); math.Abs(got) > 1e-12 {
		t.Errorf("dipole at 0deg: got %v, want 0", got)
	}
	if got := arr.ElementFactorDeg(30); math.Abs(got-0.5) > 1e-9 {
		t.Errorf("dipole at 30deg: got %v, want 0.5", got)
	}
}

func TestParseElement(t *testing.T) {
	if kind, err := ParseElement("dipole"); err != nil || kind != ElementDipole {
		t.Errorf("ParseElement(dipole): got (%v,%v), want (dipole,nil)", kind, err)
	}
	if kind, err := ParseElement("iso"); err != nil || kind != ElementIso {
		t.Errorf("ParseElement(iso): got (%v,%v), want (iso,nil)", kind, err)
	}
	if _, err := ParseElement("magic"); err == nil {
		t.Errorf("ParseElement(magic): expected error, got nil")
	}
}

func TestEndfireBetaMatchesKd(t *testing.T) {
	arr := mustArray(t, Params{N: 8, D: 0.5, Lambda: 1})
	want := -2 * math.Pi * arr.Params.D / arr.Params.Lambda
	if got := arr.EndfireBeta(); math.Abs(got-want) > 1e-12 {
		t.Errorf("EndfireBeta: got %v, want %v", got, want)
	}
}
