package scan

import (
	"errors"
	"fmt"

	"array-af/internal/geometry"
)

type ScanParams struct {
	Array       geometry.Params
	BetaStart   *float64
	BetaEnd     *float64
	Steps       int
	SearchSteps int
	HpbwSteps   int
}

func (p ScanParams) Validate() error {
	if err := geometry.Validate(p.Array); err != nil {
		return err
	}
	if p.Steps < 1 {
		return errors.New("scan steps must be >= 1")
	}
	if p.Steps > 10000 {
		return fmt.Errorf("scan steps too large: %d (max 10000)", p.Steps)
	}
	return nil
}

func (p ScanParams) EffectiveRange(arr *geometry.Array) (float64, float64) {
	start := 0.0
	if p.BetaStart != nil {
		start = *p.BetaStart
	}
	end := arr.EndfireBeta()
	if p.BetaEnd != nil {
		end = *p.BetaEnd
	}
	return start, end
}

func (p ScanParams) EffectiveSteps() int {
	if p.Steps < 1 {
		return 32
	}
	return p.Steps
}
