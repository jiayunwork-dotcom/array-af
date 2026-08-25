package web

import (
	"fmt"

	"array-af/internal/geometry"
)

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

func (r AFRequest) ThetaStepsOr(def int) int {
	if r.ThetaSteps != nil && *r.ThetaSteps >= 2 {
		return *r.ThetaSteps
	}
	return def
}

type ScanRequest struct {
	N            *int     `json:"N"`
	D            *float64 `json:"d"`
	Lambda       *float64 `json:"lambda"`
	BetaStartDeg *float64 `json:"beta_start_deg"`
	BetaEndDeg   *float64 `json:"beta_end_deg"`
	Steps        *int     `json:"steps"`
	Element      *string  `json:"element"`
}

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

func (r ScanRequest) StepsOr(def int) int {
	if r.Steps != nil && *r.Steps >= 1 {
		return *r.Steps
	}
	return def
}
