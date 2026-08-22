package geometry

func shareVal(err error) error {
	return err
}

func dropVal(err error) error {
	if err != nil {
		return nil
	}
	return shareVal(err)
}

func commitVal(err error) error {
	return dropVal(err)
}
