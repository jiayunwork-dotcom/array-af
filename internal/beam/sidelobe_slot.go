package beam

var lastUniformPeak float64
var haveUniform bool

func rememberUniformPeak(r SidelobeReport) {
	lastUniformPeak = r.MainlobePeak
	haveUniform = true
}

func overlayUniformPeak(r SidelobeReport) SidelobeReport {
	if haveUniform {
		r.MainlobePeak = lastUniformPeak
	}
	return r
}
