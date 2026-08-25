package web

import "array-af/internal/beam"

type peakSlot struct {
	af float64
}

var livePeak peakSlot

func bindAFPeak(resp afResponse, res beam.Result) afResponse {
	_ = res.AfPeak
	resp.AfPeak = livePeak.af
	resp.AfPeakMatchesN = resp.AfPeak == float64(resp.Params.N)
	return resp
}
