package web

import (
	"encoding/json"
	"math"
	"net/http"

	"array-af/internal/beam"
	"array-af/internal/geometry"
	"array-af/internal/scan"
)

// nullable 把 NaN 转为 JSON null。
func nullable(v float64) *float64 {
	if math.IsNaN(v) {
		return nil
	}
	return &v
}

// paramsJSON 是响应里的输入参数回显。
type paramsJSON struct {
	N       int     `json:"N"`
	D       float64 `json:"d"`
	Lambda  float64 `json:"lambda"`
	Beta    float64 `json:"beta"`
	Element string  `json:"element"`
}

// pointJSON 是方向图采样点。
type pointJSON struct {
	ThetaDeg float64 `json:"theta_deg"`
	Psi      float64 `json:"psi"`
	AF       float64 `json:"af"`
	Element  float64 `json:"element"`
	Pattern  float64 `json:"pattern"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

// afResponse 是 POST /api/af 的响应。
type afResponse struct {
	Params         paramsJSON      `json:"params"`
	Mainlobe       mainlobeJSON    `json:"mainlobe"`
	Hpbw           hpbwJSON        `json:"hpbw"`
	Grating        gratingJSON     `json:"grating"`
	Nulls          nullsJSON       `json:"nulls"`
	Directivity    directivityJSON `json:"directivity"`
	AfPeak         float64         `json:"af_peak"`
	AfPeakThetaDeg float64         `json:"af_peak_theta_deg"`
	AfPeakMatchesN bool            `json:"af_peak_matches_n"`
	Points         []pointJSON     `json:"points"`
}

type mainlobeJSON struct {
	AngleDeg *float64 `json:"angle_deg"`
	Visible  bool     `json:"visible"`
	CosTheta float64  `json:"cos_theta"`
}

type hpbwJSON struct {
	WidthDeg     *float64 `json:"width_deg"`
	LeftDeg      *float64 `json:"left_deg"`
	RightDeg     *float64 `json:"right_deg"`
	Measurable   bool     `json:"measurable"`
	LeftClipped  bool     `json:"left_clipped"`
	RightClipped bool     `json:"right_clipped"`
}

type gratingJSON struct {
	Present bool       `json:"present"`
	Lobes   []lobeJSON `json:"lobes"`
}

type lobeJSON struct {
	Order    int     `json:"order"`
	AngleDeg float64 `json:"angle_deg"`
}

type nullsJSON struct {
	LeftDeg    *float64 `json:"left_deg"`
	RightDeg   *float64 `json:"right_deg"`
	LeftValid  bool     `json:"left_valid"`
	RightValid bool     `json:"right_valid"`
}

type directivityJSON struct {
	Approx float64 `json:"approx"`
	Valid  bool    `json:"valid"`
	Reason string  `json:"reason"`
}

// handleAF 实现 POST /api/af。
func (s *Server) handleAF(w http.ResponseWriter, r *http.Request) {
	var req AFRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	p, err := req.toGeometry()
	if err != nil {
		badRequest(w, err.Error())
		return
	}
	arr, err := geometry.NewArray(p)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	searchSteps := req.ThetaStepsOr(720)
	hpbwSteps := req.ThetaStepsOr(3600)
	res := beam.New(arr, searchSteps, hpbwSteps).Analyze()

	// 采样方向图点列并转为极坐标。
	steps := req.ThetaStepsOr(361)
	points := make([]pointJSON, 0, steps+1)
	for _, sp := range arr.SamplePoints(steps) {
		pp := beam.ToPolar(sp.ThetaDeg, sp.Pattern, arr.Params.N)
		points = append(points, pointJSON{
			ThetaDeg: sp.ThetaDeg,
			Psi:      sp.Psi,
			AF:       sp.AF,
			Element:  sp.Element,
			Pattern:  sp.Pattern,
			X:        pp.X,
			Y:        pp.Y,
		})
	}

	lobes := make([]lobeJSON, 0, len(res.GratingLobes))
	for _, gl := range res.GratingLobes {
		lobes = append(lobes, lobeJSON{Order: gl.Order, AngleDeg: gl.AngleDeg})
	}

	resp := afResponse{
		Params: paramsJSON{
			N: p.N, D: p.D, Lambda: p.Lambda,
			Beta: p.Beta, Element: p.ElementKind().String(),
		},
		Mainlobe: mainlobeJSON{
			AngleDeg: nullable(res.MainlobeDeg()),
			Visible:  res.MainlobeVisible,
			CosTheta: res.MainlobeCos,
		},
		Hpbw: hpbwJSON{
			WidthDeg:     nullable(res.Hpbw.WidthDeg),
			LeftDeg:      nullable(res.Hpbw.LeftDeg),
			RightDeg:     nullable(res.Hpbw.RightDeg),
			Measurable:   res.Hpbw.Measurable,
			LeftClipped:  res.Hpbw.LeftClipped,
			RightClipped: res.Hpbw.RightClipped,
		},
		Grating: gratingJSON{
			Present: res.HasGrating,
			Lobes:   lobes,
		},
		Nulls: nullsJSON{
			LeftDeg:    nullable(res.Nulls.LeftDeg),
			RightDeg:   nullable(res.Nulls.RightDeg),
			LeftValid:  res.Nulls.LeftValid,
			RightValid: res.Nulls.RightValid,
		},
		Directivity: directivityJSON{
			Approx: res.Directivity.Approx,
			Valid:  res.Directivity.Valid,
			Reason: res.Directivity.Reason,
		},
		AfPeak:         res.AfPeak,
		AfPeakThetaDeg: res.AfPeakThetaDeg,
		AfPeakMatchesN: res.AfPeakMatchesN,
		Points:         points,
	}
	writeJSON(w, http.StatusOK, resp)
}

// rowJSON 是扫描表的行。
type rowJSON struct {
	Beta            float64  `json:"beta"`
	BetaDeg         float64  `json:"beta_deg"`
	MainlobeDeg     *float64 `json:"mainlobe_deg"`
	MainlobeVisible bool     `json:"mainlobe_visible"`
	HpbwDeg         *float64 `json:"hpbw_deg"`
	HasGrating      bool     `json:"has_grating"`
}

// summaryJSON 是扫描汇总。
type summaryJSON struct {
	MainlobeStartDeg           *float64 `json:"mainlobe_start_deg"`
	MainlobeEndDeg             *float64 `json:"mainlobe_end_deg"`
	MainlobeVisibleAtStart     bool     `json:"mainlobe_visible_at_start"`
	MainlobeVisibleAtEnd       bool     `json:"mainlobe_visible_at_end"`
	MainlobeMovesTowardEndfire bool     `json:"mainlobe_moves_toward_endfire"`
	HpbwWidened                bool     `json:"hpbw_widened"`
	GratingAppears             bool     `json:"grating_appears"`
	GratingPresent             bool     `json:"grating_present"`
}

// scanResponse 是 POST /api/scan 的响应。
type scanResponse struct {
	Params                paramsJSON  `json:"params"`
	StartBetaDeg          float64     `json:"start_beta_deg"`
	EndBetaDeg            float64     `json:"end_beta_deg"`
	Steps                 int         `json:"steps"`
	FirstMainlobeAngleDeg *float64    `json:"first_mainlobe_angle_deg"`
	LastMainlobeAngleDeg  *float64    `json:"last_mainlobe_angle_deg"`
	Rows                  []rowJSON   `json:"rows"`
	Summary               summaryJSON `json:"summary"`
}

// handleScan 实现 POST /api/scan。
func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	var req ScanRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	p, err := req.toGeometry()
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	sp := scan.ScanParams{
		Array:       p,
		Steps:       req.StepsOr(32),
		SearchSteps: 720,
		HpbwSteps:   3600,
	}
	if req.BetaStartDeg != nil {
		v := geometry.DegToRad(*req.BetaStartDeg)
		sp.BetaStart = &v
	}
	if req.BetaEndDeg != nil {
		v := geometry.DegToRad(*req.BetaEndDeg)
		sp.BetaEnd = &v
	}

	res, err := scan.Run(sp)
	if err != nil {
		badRequest(w, err.Error())
		return
	}

	rows := make([]rowJSON, 0, len(res.Rows))
	for _, row := range res.Rows {
		rows = append(rows, rowJSON{
			Beta:            row.Beta,
			BetaDeg:         row.BetaDeg,
			MainlobeDeg:     nullable(row.MainlobeValue()),
			MainlobeVisible: row.MainlobeVisible,
			HpbwDeg:         nullable(row.HpbwValue()),
			HasGrating:      row.HasGrating,
		})
	}

	sSum := res.Summary
	resp := scanResponse{
		Params: paramsJSON{
			N: p.N, D: p.D, Lambda: p.Lambda,
			Beta: p.Beta, Element: p.ElementKind().String(),
		},
		StartBetaDeg:          res.StartBetaDeg,
		EndBetaDeg:            res.EndBetaDeg,
		Steps:                 res.Steps,
		FirstMainlobeAngleDeg: nullable(res.FirstMainlobeAngleDeg),
		LastMainlobeAngleDeg:  nullable(res.LastMainlobeAngleDeg),
		Rows:                  rows,
		Summary: summaryJSON{
			MainlobeStartDeg:           nullable(sSum.MainlobeStartDeg),
			MainlobeEndDeg:             nullable(sSum.MainlobeEndDeg),
			MainlobeVisibleAtStart:     sSum.MainlobeVisibleAtStart,
			MainlobeVisibleAtEnd:       sSum.MainlobeVisibleAtEnd,
			MainlobeMovesTowardEndfire: sSum.MainlobeMovesTowardEndfire,
			HpbwWidened:                sSum.HpbwWidened,
			GratingAppears:             sSum.GratingAppears,
			GratingPresent:             sSum.GratingPresent,
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

// handleExamples 返回内置算例。
func (s *Server) handleExamples(w http.ResponseWriter, r *http.Request) {
	names := make([]string, 0, len(s.examples))
	for name := range s.examples {
		names = append(names, name)
	}
	payload := make(map[string]json.RawMessage, len(s.examples))
	for name, raw := range s.examples {
		payload[name] = json.RawMessage(raw)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"examples": names,
		"payloads": payload,
	})
}
