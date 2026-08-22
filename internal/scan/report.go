package scan

import (
	"math"

	"array-af/internal/beam"
	"array-af/internal/geometry"
)

// Result 是一次扫描的完整输出。
type Result struct {
	// Rows 是逐 β 的主瓣表。
	Rows []Row
	// StartBetaDeg / EndBetaDeg 是实际扫描区间（度）。
	StartBetaDeg float64
	EndBetaDeg   float64
	// Steps 是实际扫描段数。
	Steps int
	// Summary 是汇总判定。
	Summary Summary
	// FirstMainlobeAngleDeg / LastMainlobeAngleDeg 是首末可见主瓣角（度）。
	// 若无可见主瓣则为 NaN。
	FirstMainlobeAngleDeg float64
	LastMainlobeAngleDeg  float64
}

// Run 执行 β 扫描并汇总。
func Run(p ScanParams) (Result, error) {
	if err := p.Validate(); err != nil {
		return Result{}, err
	}
	arr, err := geometry.NewArray(p.Array)
	if err != nil {
		return Result{}, err
	}
	start, end := p.EffectiveRange(arr)
	steps := p.EffectiveSteps()
	betas := BetaSweep(start, end, steps)

	res := Result{
		StartBetaDeg: start * 180 / math.Pi,
		EndBetaDeg:   end * 180 / math.Pi,
		Steps:        steps,
		Rows:         make([]Row, 0, len(betas)),
	}
	for _, beta := range betas {
		r := beam.AnalyzeWithBeta(arr, beta, p.SearchSteps, p.HpbwSteps)
		row := Row{
			Beta:            beta,
			BetaDeg:         beta * 180 / math.Pi,
			MainlobeDeg:     r.MainlobeAngleDeg,
			MainlobeVisible: r.MainlobeVisible,
			HpbwDeg:         math.NaN(),
			HasGrating:      r.HasGrating,
		}
		if r.Hpbw.Measurable {
			row.HpbwDeg = r.Hpbw.WidthDeg
		}
		res.Rows = append(res.Rows, row)
	}
	res.Summary = summarize(res.Rows)
	first, last := firstLastMainlobe(res.Rows)
	res.FirstMainlobeAngleDeg = first
	res.LastMainlobeAngleDeg = last
	return res, nil
}

// summarize 从行表推导交叉规则汇总。
func summarize(rows []Row) Summary {
	s := Summary{}
	if len(rows) == 0 {
		return s
	}
	first := rows[0]
	last := rows[len(rows)-1]
	s.MainlobeVisibleAtStart = first.MainlobeVisible
	s.MainlobeVisibleAtEnd = last.MainlobeVisible
	s.MainlobeStartDeg = first.MainlobeValue()
	s.MainlobeEndDeg = last.MainlobeValue()

	// 主瓣向端射移动：角度序列非递增（90°→0°）。
	s.MainlobeMovesTowardEndfire = true
	prev := math.Inf(1)
	for _, r := range rows {
		v := r.MainlobeValue()
		if math.IsNaN(v) {
			continue
		}
		if v > prev+1e-9 {
			s.MainlobeMovesTowardEndfire = false
			break
		}
		prev = v
	}

	// HPBW 变宽：末端宽度 > 起点宽度。
	startW := first.HpbwValue()
	endW := last.HpbwValue()
	if !math.IsNaN(startW) && !math.IsNaN(endW) && endW > startW+1e-9 {
		s.HpbwWidened = true
	}

	// 栅瓣：起点无、终点有。
	if !first.HasGrating && last.HasGrating {
		s.GratingAppears = true
	}
	for _, r := range rows {
		if r.HasGrating {
			s.GratingPresent = true
			break
		}
	}
	return s
}

// firstLastMainlobe 返回首末可见主瓣角。
func firstLastMainlobe(rows []Row) (float64, float64) {
	first, last := math.NaN(), math.NaN()
	for _, r := range rows {
		if r.MainlobeVisible {
			if math.IsNaN(first) {
				first = r.MainlobeDeg
			}
			last = r.MainlobeDeg
		}
	}
	return first, last
}
