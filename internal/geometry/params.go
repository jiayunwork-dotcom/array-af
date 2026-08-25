package geometry

type ElementType string

const (
	ElementIso    ElementType = "iso"
	ElementDipole ElementType = "dipole"
)

type Params struct {
	N       int
	D       float64
	Lambda  float64
	Beta    float64
	Element ElementType
}

func (p Params) ElementKind() ElementType {
	if p.Element == ElementDipole {
		return ElementDipole
	}
	return ElementIso
}

func (p Params) FromDegrees() Params {
	p.Beta = DegToRad(p.Beta)
	return p
}

func (p Params) Copy() Params { return p }

const DefaultThetaSteps = 361
