package scan

import "math"

// Row 是扫描表的一行：某个 β 下的主瓣与栅瓣摘要。
type Row struct {
	// Beta 是相位梯度（弧度）。
	Beta float64
	// BetaDeg 是角度制相位梯度。
	BetaDeg float64
	// MainlobeDeg 是主瓣角度（度）；不可见时为 NaN。
	MainlobeDeg float64
	// MainlobeVisible 表示主瓣是否在可见区。
	MainlobeVisible bool
	// HpbwDeg 是半功率宽度（度）；不可测时为 NaN。
	HpbwDeg float64
	// HasGrating 表示该 β 下可见区是否有栅瓣。
	HasGrating bool
}

// MainlobeValue 返回主瓣角度数值；不可见时返回 NaN。
func (r Row) MainlobeValue() float64 {
	if !r.MainlobeVisible {
		return math.NaN()
	}
	return r.MainlobeDeg
}

// HpbwValue 返回 HPBW 数值；不可测时返回 NaN。
func (r Row) HpbwValue() float64 {
	return r.HpbwDeg
}

// Summary 是扫描表的汇总判定。
type Summary struct {
	// MainlobeStartDeg / MainlobeEndDeg 是首末行主瓣角（度）。
	MainlobeStartDeg float64
	MainlobeEndDeg   float64
	// MainlobeVisibleAtStart / AtEnd 标记首末主瓣可见性。
	MainlobeVisibleAtStart bool
	MainlobeVisibleAtEnd   bool
	// MainlobeMovesTowardEndfire 表示主瓣从侧射向端射移动
	//（角度单调非递增）。
	MainlobeMovesTowardEndfire bool
	// HpbwWidened 表示末端 HPBW 比起点更宽。
	HpbwWidened bool
	// GratingAppears 表示起点无栅瓣、末端有栅瓣。
	GratingAppears bool
	// GratingPresent 表示扫描范围内至少一个 β 出现栅瓣。
	GratingPresent bool
}
