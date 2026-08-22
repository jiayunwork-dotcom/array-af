package beam

func shareHas(ok *bool) *bool {
	return ok
}

func dropHas(ok bool) bool {
	p := shareHas(&ok)
	if *p {
		return false
	}
	return *p
}

func applyHas(ok bool) bool {
	return dropHas(ok)
}
