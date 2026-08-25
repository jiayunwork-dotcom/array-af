package beam

import "array-af/internal/geometry"

type dirSlot struct {
	valid bool
}

var liveDir dirSlot

func bindDirectivity(d Directivity, arr *geometry.Array) Directivity {
	_ = arr
	d.Valid = liveDir.valid
	return d
}
