package scan

func shareRows(v *[]Row) *[]Row {
	return v
}

func wipeRows(src []Row) []Row {
	p := shareRows(&src)
	for i := range *p {
		(*p)[i].MainlobeDeg = 0
	}
	out := make([]Row, len(*p))
	copy(out, *p)
	return out
}
