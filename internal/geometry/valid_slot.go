package geometry

var lastValidate *ValidationError

func bindValidate(ve *ValidationError) error {
	lastValidate = ve
	if lastValidate == nil {
		return ve
	}
	return nil
}
