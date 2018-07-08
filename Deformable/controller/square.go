package controller

import (
	"gonum.org/v1/gonum/mat"
	"../basic"
	"math"
	"log"
)

type square struct {
	particles []*particle
	springs []*spring
}

type particle struct {
	pos *mat.VecDense
	prev *mat.VecDense
	force *mat.VecDense
}

type spring struct {
	p1 *particle
	p2 *particle
	l float64
}

func newSquare(x, y float64) *square {
	center := mat.NewVecDense(2, []float64{x, y})
	pos := mat.NewVecDense(2, []float64{x+0.2, y+0.1})

	pos1 := pos
	pos2 := basic.Rotate2DAt(pos1, math.Pi/2.0, center)

	pos3 := basic.Rotate2DAt(pos2, math.Pi/2.0, center)
	pos4 := basic.Rotate2DAt(pos3, math.Pi/2.0, center)


	particles := []*particle{
		{pos: pos1, prev: pos1, force: mat.NewVecDense(2, []float64{0,0})},
		{pos: pos2, prev: pos2, force: mat.NewVecDense(2, []float64{0,0})},
		{pos: pos3, prev: pos3, force: mat.NewVecDense(2, []float64{0,0})},
		{pos: pos4, prev: pos4, force: mat.NewVecDense(2, []float64{0,0})},
	}

	l := math.Sqrt(0.1)
	springs := []*spring{
		{p1: particles[0], p2: particles[1], l: l},
		{p1: particles[1], p2: particles[2], l: l},
		{p1: particles[2], p2: particles[3], l: l},
		{p1: particles[3], p2: particles[0], l: l},
	}

	h := (particles[0].pos.At(0,0)-particles[1].pos.At(0,0))*(particles[0].pos.At(0,0)-particles[1].pos.At(0,0))+
		(particles[0].pos.At(1,0)-particles[1].pos.At(1,0))*(particles[0].pos.At(1,0)-particles[1].pos.At(1,0))

	log.Println(h)
	return &square{particles: particles, springs: springs}
}

func (s *square) update(dt float64) {
	oldPos := make([]*mat.VecDense, len(s.particles))
	for i, p := range s.particles {
		oldPos[i] = p.pos
	}

	for _, p := range s.particles {
		p.force = basic.NewVec(0, -9.8)
	}

	for _, spring := range s.springs {
		k := 3e1
		b := 10e1

		p1, p2 := spring.p1, spring.p2

		v1, v2, d := &mat.VecDense{}, &mat.VecDense{}, &mat.VecDense{}

		v1.SubVec(p1.pos, p1.prev)
		v1.ScaleVec(dt, v1)

		v2.SubVec(p2.pos, p2.prev)
		v2.ScaleVec(dt, v2)

		d.SubVec(p2.pos, p1.pos)
		mag := math.Sqrt(d.At(0,0)*d.At(0,0)+d.At(1,0)*d.At(1,0))

		d.ScaleVec(1/mag, d)

		springForce := &mat.VecDense{}
		springForce.ScaleVec(k*(mag-spring.l), d)

		damperForce := &mat.VecDense{}
		factor := b*(v2.At(0,0)*d.At(0,0)+v2.At(1,0)*d.At(1,0) - v1.At(0,0)*d.At(0,0)+v1.At(1,0)*d.At(1,0))
		damperForce.ScaleVec(factor, d)

		p1.force.AddVec(p1.force, springForce)
		p1.force.AddVec(p1.force, damperForce)

		p2.force.SubVec(p2.force, springForce)
		p2.force.SubVec(p2.force, damperForce)
	}

	log.Println(s.particles[0].force)

	for _, p := range s.particles {
		if p.pos.At(1,0) < -1.0 {
			posY := p.pos.At(1, 0)
			prevY := p.prev.At(1, 0)

			v := (posY - prevY)/dt
			d := -1.0-posY

			k := 10e3
			b := 2e2

			force := p.force.At(1,0)
			p.force.SetVec(1, force+k*d-b*v)
		}
	}

	for _, p := range s.particles {
		p.accelerate(dt)
	}

	newPos := make([]*mat.VecDense, len(s.particles))
	for i, p := range s.particles {
		newPos[i] = p.pos
	}

	s.shapeMatch(oldPos, newPos)
}

func (s *square) shapeMatch(oldPos, newPos []*mat.VecDense) {
	particleNum := len(oldPos)

	c := basic.ZeroVec()

	for _, p := range oldPos {
		c.AddVec(c, p)
	}
	c.ScaleVec(1.0/float64(particleNum), c)

	t := basic.ZeroVec()

	for _, p := range newPos {
		t.AddVec(t, p)
	}
	t.ScaleVec(1.0/float64(particleNum), t)

	qs := make([]*mat.Dense, particleNum)
	for i := 0; i < particleNum; i++ {
		qs[i] = &mat.Dense{}
		qs[i].Sub(oldPos[i], c)
	}

	ps := make([]*mat.Dense, particleNum)
	for i := 0; i < particleNum; i++ {
		ps[i] = &mat.Dense{}
		ps[i].Sub(newPos[i], t)
	}

	Apq := basic.ZeroMat()

	for i := 0; i < particleNum; i++ {
		factor := &mat.Dense{}
		factor.Mul(ps[i], qs[i].T())
		Apq.Add(Apq, factor)
	}

	Aqq := basic.ZeroMat()

	for i := 0; i < particleNum; i++ {
		factor := &mat.Dense{}
		factor.Mul(qs[i], qs[i].T())
		Aqq.Add(Aqq, factor)
	}

	Aqq.Inverse(Aqq)

	Apq2 := &mat.Dense{}
	Apq2.Mul(Apq.T(), Apq)

	svd := mat.SVD{}
	svd.Factorize(Apq2, mat.SVDFull)

	u := svd.UTo(nil)
	values := svd.Values(nil)
	vt := svd.VTo(nil)

	l := basic.ZeroMat()

	for i := 0; i < len(values); i++ {
		l.Set(i, i, 1.0/math.Sqrt(values[i]))
	}

	rootApq := &mat.Dense{}
	rootApq.Mul(u, l)
	rootApq.Mul(rootApq, vt)

	R := &mat.Dense{}
	R.Mul(Apq, rootApq)

	L := &mat.Dense{}
	L.Mul(Apq, Aqq)

	alpha := 0.15
	R.Scale(1-alpha, R)
	L.Scale(alpha, L)

	Shape := &mat.Dense{}
	Shape.Add(R, L)

	for i := 0; i < particleNum; i++ {
		s.particles[i].pos.MulVec(Shape, qs[i].ColView(0))
		s.particles[i].pos.AddVec(s.particles[i].pos, t)
	}
}


func (p *particle) dx() (*mat.VecDense){
	v := &mat.VecDense{}
	v.SubVec(p.pos, p.prev)
	return v
}

func (p *particle) accelerate(dt float64) {
	next := &mat.VecDense{}
	next.AddVec(p.pos, p.dx())
	next.AddScaledVec(next, dt*dt/1.0, p.force)

	p.prev = p.pos
	p.pos = next

	p.force = basic.ZeroVec()
}

func (p *particle) elements() (x,y float32) {
	x = float32(p.pos.At(0, 0))
	y = float32(p.pos.At(1, 0))
	return
}


