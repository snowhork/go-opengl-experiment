package line_segment

type Poly interface {
	calc(float64) float64
}

type Poly5 struct {
	a0, a1, a2, a3, a4, a5 float64
}

func (p *Poly5) calc(t float64) float64 {
	return p.a0*t*t*t*t*t + p.a1*t*t*t*t + p.a2*t*t*t + p.a3*t*t + p.a4*t + p.a5
}

type Poly4 struct {
	a0, a1, a2, a3, a4 float64
}

func (p *Poly4) calc(t float64) float64 {
	return p.a0*t*t*t*t + p.a1*t*t*t + p.a2*t*t + p.a3*t + p.a4
}

type Poly3 struct {
	a0, a1, a2, a3 float64
}

func (p *Poly3) calc(t float64) float64 {
	return p.a0*t*t*t + p.a1*t*t + p.a2*t + p.a3
}

type Poly2 struct {
	a0, a1, a2 float64
}

func (p *Poly2) calc(t float64) float64 {
	return p.a0*t*t + p.a1*t + p.a2
}

type Poly1 struct {
	a0, a1 float64
}

func (p *Poly1) calc(t float64) float64 {
	return p.a0*t + p.a1
}

type Poly0 struct {
	a0  float64
}

func (p *Poly0) calc(t float64) float64 {
	return p.a0
}

