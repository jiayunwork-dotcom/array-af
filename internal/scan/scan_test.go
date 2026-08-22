package scan

import (
	"math"
	"testing"

	"array-af/internal/geometry"
)

func TestScanMainlobeTracksBeta(t *testing.T) {
	res, err := Run(ScanParams{
		Array: geometry.Params{N: 8, D: 0.5, Lambda: 1},
		Steps: 16,
	})
	if err != nil {
		t.Fatalf("Run: unexpected error %v", err)
	}
	if len(res.Rows) != 17 {
		t.Errorf("rows: got %d, want 17", len(res.Rows))
	}
	first := res.Rows[0]
	last := res.Rows[len(res.Rows)-1]
	if math.Abs(first.MainlobeDeg-90) > 1e-6 {
		t.Errorf("first row mainlobe: got %vdeg, want 90deg (broadside)", first.MainlobeDeg)
	}
	if math.Abs(last.MainlobeDeg) > 1e-6 {
		t.Errorf("last row mainlobe: got %vdeg, want 0deg (endfire)", last.MainlobeDeg)
	}
	if !res.Summary.MainlobeMovesTowardEndfire {
		t.Errorf("expected mainlobe moving 90->0 across beta sweep")
	}
}

func TestScanHpbwWidens(t *testing.T) {
	res, err := Run(ScanParams{
		Array: geometry.Params{N: 8, D: 0.5, Lambda: 1},
		Steps: 8,
	})
	if err != nil {
		t.Fatalf("Run: unexpected error %v", err)
	}
	if !res.Summary.HpbwWidened {
		t.Errorf("expected hpbw to widen from broadside to endfire")
	}
	first := res.Rows[0].HpbwValue()
	last := res.Rows[len(res.Rows)-1].HpbwValue()
	if math.IsNaN(first) || math.IsNaN(last) {
		t.Errorf("expected measurable hpbw in rows, got first=%v last=%v", first, last)
	}
}

func TestScanGratingAppearsForWideSpacing(t *testing.T) {
	res, err := Run(ScanParams{
		Array: geometry.Params{N: 8, D: 1.2, Lambda: 1},
		Steps: 8,
	})
	if err != nil {
		t.Fatalf("Run: unexpected error %v", err)
	}
	if !res.Summary.GratingPresent {
		t.Errorf("d=1.2*lambda: expected grating present in sweep")
	}
	if !res.Rows[0].HasGrating {
		t.Errorf("d=1.2*lambda beta=0: expected grating at broadside")
	}
}

func TestScanNoGratingAtBroadsideForHalfWave(t *testing.T) {
	res, err := Run(ScanParams{
		Array: geometry.Params{N: 8, D: 0.5, Lambda: 1},
		Steps: 8,
	})
	if err != nil {
		t.Fatalf("Run: unexpected error %v", err)
	}
	if res.Rows[0].HasGrating {
		t.Errorf("d=lambda/2 beta=0: expected no grating, got %v", res.Rows[0].HasGrating)
	}
}

func TestScanRejectsBadParams(t *testing.T) {
	if _, err := Run(ScanParams{Array: geometry.Params{N: 1, D: 0.5, Lambda: 1}}); err == nil {
		t.Errorf("N=1: expected error, got nil")
	}
	if _, err := Run(ScanParams{Array: geometry.Params{N: 8, D: 0, Lambda: 1}}); err == nil {
		t.Errorf("d=0: expected error, got nil")
	}
	if _, err := Run(ScanParams{Array: geometry.Params{N: 8, D: 0.5, Lambda: -1}}); err == nil {
		t.Errorf("lambda=-1: expected error, got nil")
	}
}

func TestScanRespectsCustomRange(t *testing.T) {
	start := 0.0
	end := -1.0
	res, err := Run(ScanParams{
		Array:     geometry.Params{N: 8, D: 0.5, Lambda: 1},
		BetaStart: &start,
		BetaEnd:   &end,
		Steps:     4,
	})
	if err != nil {
		t.Fatalf("Run: unexpected error %v", err)
	}
	if len(res.Rows) != 5 {
		t.Errorf("rows: got %d, want 5", len(res.Rows))
	}
	if math.Abs(res.StartBetaDeg) > 1e-9 {
		t.Errorf("start beta: got %vdeg, want 0", res.StartBetaDeg)
	}
	wantEnd := -180 / math.Pi
	if math.Abs(res.EndBetaDeg-wantEnd) > 1e-6 {
		t.Errorf("end beta: got %vdeg, want %vdeg (-1 rad)", res.EndBetaDeg, wantEnd)
	}
}

func TestThresholdSpacingRatio(t *testing.T) {
	if got := ThresholdSpacingRatio(0); math.Abs(got-1) > 1e-12 {
		t.Errorf("ThresholdSpacingRatio(0): got %v, want 1", got)
	}
	if got := GratingThresholdSpacing(1, 0); math.Abs(got-1) > 1e-12 {
		t.Errorf("GratingThresholdSpacing(1,0): got %v, want 1", got)
	}
	if !RatioBelowThreshold(0.5, 0) {
		t.Errorf("d/lambda=0.5 should be below broadside grating threshold")
	}
}
