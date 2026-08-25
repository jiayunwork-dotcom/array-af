package geometry

import (
	"fmt"
	"math"
)

type Element struct {
	kind ElementType
}

func NewElement(kind ElementType) Element {
	return Element{kind: kind}
}

func ParseElement(s string) (ElementType, error) {
	switch ElementType(s) {
	case ElementIso, "":
		return ElementIso, nil
	case ElementDipole:
		return ElementDipole, nil
	default:
		return "", fmt.Errorf("unknown element factor %q (want iso or dipole)", s)
	}
}

func (e Element) Kind() ElementType {
	if e.kind == ElementDipole {
		return ElementDipole
	}
	return ElementIso
}

func (e Element) Factor(theta float64) float64 {
	if e.Kind() == ElementDipole {
		s := math.Sin(theta)
		if s < 0 {
			s = -s
		}
		return s
	}
	return 1
}

func (e Element) FactorDeg(thetaDeg float64) float64 {
	return e.Factor(DegToRad(thetaDeg))
}

func (e Element) Maximum() float64 {
	return 1
}

func (e Element) IsDipole() bool {
	return e.Kind() == ElementDipole
}

func (e ElementType) String() string {
	if e == ElementDipole {
		return string(ElementDipole)
	}
	return string(ElementIso)
}

func (e Element) String() string {
	return e.Kind().String()
}
