package controller

import "../basic"
import "gonum.org/v1/gonum/mat"

type particle struct {
	mass float64
	pos *mat.VecDense
	vel *mat.VecDense
	fixed bool
}

func NewParticle(mass float64, pos *mat.VecDense, fixed bool) *particle {
	return &particle{mass: 0.1, pos: pos, vel: basic.ZeroVec(), fixed: fixed}
}

func (p *particle) glpos() (x,y,z float32) {
	return float32(p.pos.AtVec(0)),
		float32(p.pos.AtVec(1)),
		float32(p.pos.AtVec(2))
}