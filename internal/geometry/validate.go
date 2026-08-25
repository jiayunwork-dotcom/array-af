package geometry

import (
	"errors"
	"fmt"
	"math"
)

type FieldError struct {
	Field string
	Value any
	Rule  string
}

func (e FieldError) Error() string {
	return fmt.Sprintf("%s: %v violates %s", e.Field, e.Value, e.Rule)
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	if e == nil || len(e.Fields) == 0 {
		return "validation failed"
	}
	msg := "invalid input:"
	for i, f := range e.Fields {
		if i > 0 {
			msg += ";"
		}
		msg += " " + f.Error()
	}
	return msg
}

func (e *ValidationError) FieldCount() int {
	if e == nil {
		return 0
	}
	return len(e.Fields)
}

func (e *ValidationError) Add(f FieldError) {
	e.Fields = append(e.Fields, f)
}

func Validate(p Params) error {
	ve := &ValidationError{}
	if p.N < 2 {
		ve.Add(FieldError{
			Field: "N",
			Value: p.N,
			Rule:  "must be an integer >= 2",
		})
	}
	if p.D <= 0 {
		ve.Add(FieldError{
			Field: "d",
			Value: p.D,
			Rule:  "must be > 0",
		})
	}
	if p.Lambda <= 0 {
		ve.Add(FieldError{
			Field: "lambda",
			Value: p.Lambda,
			Rule:  "must be > 0",
		})
	}
	if ve.FieldCount() == 0 {
		return nil
	}
	return ve
}

var ErrNotFinite = errors.New("input contains NaN or Inf")

func CheckFinite(p Params) error {
	if p.D != p.D || p.Lambda != p.Lambda || p.Beta != p.Beta {
		return ErrNotFinite
	}
	if !finite(p.D) || !finite(p.Lambda) || !finite(p.Beta) {
		return ErrNotFinite
	}
	return nil
}

func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}
