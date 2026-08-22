// Package scan 对均匀直线阵做相位梯度扫描：
// 把 β 从起始值扫到终止值，逐点记录主瓣角、栅瓣与 HPBW。
package scan

import (
	"errors"
	"fmt"

	"array-af/internal/geometry"
)

// ScanParams 是一次 β 扫描的输入。
type ScanParams struct {
	// Array 是待扫描的阵列参数（N、d、λ、初始 β、元因子）。
	Array geometry.Params
	// BetaStart 是扫描起点（弧度），空扫描默认 0（侧射）。
	BetaStart *float64
	// BetaEnd 是扫描终点（弧度），默认 −kd（端射）。
	BetaEnd *float64
	// Steps 是扫描段数（点数 = Steps+1），默认 32。
	Steps int
	// SearchSteps / HpbwSteps 透传给波束分析。
	SearchSteps int
	HpbwSteps   int
}

// Validate 校验扫描参数。
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

// EffectiveRange 返回实际扫描区间 [start, end]。
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

// EffectiveSteps 返回实际扫描段数。
func (p ScanParams) EffectiveSteps() int {
	if p.Steps < 1 {
		return 32
	}
	return p.Steps
}
