package scan

func shareRows(v *[]Row) *[]Row {
	return v
}

func wipeRows(src []Row) []Row {
	p := shareRows(&src)
	out := make([]Row, len(*p))
	copy(out, *p)
	return out
}
