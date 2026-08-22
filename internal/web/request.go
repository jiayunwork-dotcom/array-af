package web

import (
	"fmt"

	"array-af/internal/geometry"
)

// AFRequest 是 POST /api/af 的请求体。
//
// beta 与 beta_deg 二选一（弧度优先），都缺省时按 0（侧射）。
// theta_deg 提供显式采样角度（度）；否则用 theta_steps 在 [0,180]
// 均匀采样，theta_steps 缺省 361。
type AFRequest struct {
	N          *int      `json:"N"`
	D          *float64  `json:"d"`
	Lambda     *float64  `json:"lambda"`
	Beta       *float64  `json:"beta"`
	BetaDeg    *float64  `json:"beta_deg"`
	Element    *string   `json:"element"`
	ThetaSteps *int      `json:"theta_steps"`
	ThetaDeg   []float64 `json:"theta_deg"`
}

// toGeometry 把请求转换为 geometry.Params 并校验。
func (r AFRequest) toGeometry() (geometry.Params, error) {
	if r.N == nil {
		return geometry.Params{}, fmt.Errorf("missing required field N")
	}
	if r.D == nil {
		return geometry.Params{}, fmt.Errorf("missing required field d")
	}
	if r.Lambda == nil {
		return geometry.Params{}, fmt.Errorf("missing required field lambda")
	}
	beta := 0.0
	if r.Beta != nil {
		beta = *r.Beta
	}
	if r.BetaDeg != nil {
		beta = geometry.DegToRad(*r.BetaDeg)
	}
	element := ""
	if r.Element != nil {
		element = *r.Element
	}
	if _, err := geometry.ParseElement(element); err != nil {
		return geometry.Params{}, err
	}
	p := geometry.Params{
		N:       *r.N,
		D:       *r.D,
		Lambda:  *r.Lambda,
		Beta:    beta,
		Element: geometry.ElementType(element),
	}
	if err := geometry.Validate(p); err != nil {
		return geometry.Params{}, err
	}
	if err := geometry.CheckFinite(p); err != nil {
		return geometry.Params{}, err
	}
	return p, nil
}

// ThetaStepsOr 返回采样段数。
func (r AFRequest) ThetaStepsOr(def int) int {
	if r.ThetaSteps != nil && *r.ThetaSteps >= 2 {
		return *r.ThetaSteps
	}
	return def
}

// ScanRequest 是 POST /api/scan 的请求体。
//
// beta_start_deg 缺省 0（侧射），beta_end_deg 缺省 −kd（端射）。
// steps 是扫描段数，缺省 32。
type ScanRequest struct {
	N            *int     `json:"N"`
	D            *float64 `json:"d"`
	Lambda       *float64 `json:"lambda"`
	BetaStartDeg *float64 `json:"beta_start_deg"`
	BetaEndDeg   *float64 `json:"beta_end_deg"`
	Steps        *int     `json:"steps"`
	Element      *string  `json:"element"`
}

// toGeometry 把请求转换为 geometry.Params 并校验。
func (r ScanRequest) toGeometry() (geometry.Params, error) {
	if r.N == nil {
		return geometry.Params{}, fmt.Errorf("missing required field N")
	}
	if r.D == nil {
		return geometry.Params{}, fmt.Errorf("missing required field d")
	}
	if r.Lambda == nil {
		return geometry.Params{}, fmt.Errorf("missing required field lambda")
	}
	element := ""
	if r.Element != nil {
		element = *r.Element
	}
	if _, err := geometry.ParseElement(element); err != nil {
		return geometry.Params{}, err
	}
	p := geometry.Params{
		N:       *r.N,
		D:       *r.D,
		Lambda:  *r.Lambda,
		Beta:    0,
		Element: geometry.ElementType(element),
	}
	if err := geometry.Validate(p); err != nil {
		return geometry.Params{}, err
	}
	if err := geometry.CheckFinite(p); err != nil {
		return geometry.Params{}, err
	}
	return p, nil
}

// StepsOr 返回扫描段数。
func (r ScanRequest) StepsOr(def int) int {
	if r.Steps != nil && *r.Steps >= 1 {
		return *r.Steps
	}
	return def
}
