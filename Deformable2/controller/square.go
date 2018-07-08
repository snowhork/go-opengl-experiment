package controller

import (
	"gonum.org/v1/gonum/mat"
	"../basic"
	"math"
)

type square struct {
	particles [][]*particle
	verticalSprings [][]*spring
	horizontalSprings [][]*spring
	tilt1Springs [][]*spring
	tilt2Springs [][]*spring
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

//const Level = 5
//const Length = 0.25
//const k = 700
//const b = 20000

const Level = 5
const Length = 0.25
const k = 700
const b = 25000


func newSquare(x, y float64) *square {
	cx := x
	cy := y
	idx := Length*math.Cos(math.Pi/6)/(Level-1)
	idy := Length*math.Sin(math.Pi/6)/(Level-1)

	jdx := Length*math.Cos(math.Pi/6 + math.Pi/2)/(Level-1)
	jdy := Length*math.Sin(math.Pi/6 + math.Pi/2)/(Level-1)

	particles := make([][]*particle, Level)
	for i := 0; i < Level; i++ {
		particles[i] = make([]*particle, Level)
		for j := 0; j < Level; j++ {
			tmpPos := mat.NewVecDense(2, []float64{cx+idx*float64(i)+jdx*float64(j),cy+idy*float64(i)+jdy*float64(j)})
			particles[i][j] = &particle{pos: tmpPos, prev: tmpPos, force: basic.ZeroVec()}
		}
	}

	verticalSprings := make([][]*spring, Level)
	for i := 0; i < Level; i++ {
		verticalSprings[i] = make([]*spring, Level-1)
		for j := 0; j < Level-1; j++ {
			verticalSprings[i][j] = &spring{p1: particles[i][j], p2: particles[i][j+1], l: Length/(Level-1)}
		}
	}

	horizontalSprings := make([][]*spring, Level-1)
	for i := 0; i < Level-1; i++ {
		horizontalSprings[i] = make([]*spring, Level)
		for j := 0; j < Level; j++ {
			horizontalSprings[i][j] = &spring{p1: particles[i][j], p2: particles[i+1][j], l: Length/(Level-1)}
		}
	}

	tilt1Springs := make([][]*spring, Level-1)
	for i := 0; i < Level-1; i++ {
		tilt1Springs[i] = make([]*spring, Level-1)
		for j := 0; j < Level-1; j++ {
			tilt1Springs[i][j] = &spring{p1: particles[i][j], p2: particles[i+1][j+1], l: Length/(Level-1)*math.Sqrt2}
		}
	}

	tilt2Springs := make([][]*spring, Level-1)
	for i := 0; i < Level-1; i++ {
		tilt2Springs[i] = make([]*spring, Level-1)
		for j := 0; j < Level-1; j++ {
			tilt2Springs[i][j] = &spring{p1: particles[i+1][j], p2: particles[i][j+1], l: Length/(Level-1)*math.Sqrt2}
		}
	}


	return &square{particles: particles, verticalSprings:verticalSprings, horizontalSprings:horizontalSprings,
	tilt1Springs:tilt1Springs, tilt2Springs:tilt2Springs}
}

func (s *spring) update(dt float64) {
	p1, p2 := s.p1, s.p2

	v1, v2, d := &mat.VecDense{}, &mat.VecDense{}, &mat.VecDense{}

	v1.SubVec(p1.pos, p1.prev)
	v1.ScaleVec(dt, v1)

	v2.SubVec(p2.pos, p2.prev)
	v2.ScaleVec(dt, v2)

	d.SubVec(p2.pos, p1.pos)
	mag := math.Sqrt(d.At(0,0)*d.At(0,0)+d.At(1,0)*d.At(1,0))

	d.ScaleVec(1/mag, d)

	springForce := &mat.VecDense{}
	springForce.ScaleVec(k*(mag-s.l), d)

	damperForce := &mat.VecDense{}
	factor := b*((v2.At(0,0)*d.At(0,0)+v2.At(1,0)*d.At(1,0)) - (v1.At(0,0)*d.At(0,0)+v1.At(1,0)*d.At(1,0)))
	damperForce.ScaleVec(factor, d)

	p1.force.AddVec(p1.force, springForce)
	p1.force.AddVec(p1.force, damperForce)

	p2.force.SubVec(p2.force, springForce)
	p2.force.SubVec(p2.force, damperForce)
}

func (s *square) update(dt float64) {
	oldPos := make([][]*mat.VecDense, Level)
	for i := 0; i < Level; i++ {
		oldPos[i] = make([]*mat.VecDense, Level)
		for j := 0; j < Level; j++ {
			oldPos[i][j] = s.particles[i][j].pos
		}
	}


	for _, particles := range s.particles {
		for _, p := range particles {
			p.force = basic.NewVec(0, -4)
		}
	}

	for _, springs := range s.verticalSprings {
		for _ , spring := range springs {
			spring.update(dt)
		}
	}
	for _, springs := range s.horizontalSprings {
		for _ , spring := range springs {
			spring.update(dt)
		}
	}

	for _, springs := range s.tilt1Springs {
		for _ , spring := range springs {
			spring.update(dt)
		}
	}

	for _, springs := range s.tilt2Springs {
		for _ , spring := range springs {
			spring.update(dt)
		}
	}


	for _, particles := range s.particles {
		for _, p := range particles {
			if p.pos.At(1, 0) < -1.0 {
				posY := p.pos.At(1, 0)
				prevY := p.prev.At(1, 0)

				p.pos.SetVec(1, -1.0)
				p.prev.SetVec(1, -1.0 - (prevY+(-1.0-posY) + 1.0))
				//p.prev.SetVec(1, -1.0)
			}
		}
	}

	for _, particles := range s.particles {
		for _, p := range particles {
			p.accelerate(dt)
		}
	}

	newPos := make([][]*mat.VecDense, Level)
	for i := 0; i < Level; i++ {
		newPos[i] = make([]*mat.VecDense, Level)
		for j := 0; j < Level; j++ {
			newPos[i][j] = s.particles[i][j].pos
		}
	}

	//s.shapeMatch(oldPos, newPos)

}

func (s *square) shapeMatch(oldPos, newPos [][]*mat.VecDense) {
	c := basic.ZeroVec()
	for _, ps := range oldPos {
		for _, p := range ps {
			c.AddVec(c, p)
		}
	}
	c.ScaleVec(1.0/float64(Level*Level), c)

	t := basic.ZeroVec()
	for _, ps := range newPos {
		for _, p := range ps {
			t.AddVec(t, p)
		}
	}
	t.ScaleVec(1.0/float64(Level*Level), t)

	qs := make([][]*mat.Dense, Level)
	for i := 0; i < Level; i++ {
		qs[i] = make([]*mat.Dense, Level)
		for j := 0; j < Level; j++ {
			qs[i][j] = &mat.Dense{}
			qs[i][j].Sub(oldPos[i][j], c)
		}
	}

	ps := make([][]*mat.Dense, Level)
	for i := 0; i < Level; i++ {
		ps[i] = make([]*mat.Dense, Level)
		for j := 0; j < Level; j++ {
			ps[i][j] = &mat.Dense{}
			ps[i][j].Sub(newPos[i][j], t)
		}
	}

	Apq := basic.ZeroMat()

	for i := 0; i < Level; i++ {
		for j := 0; j < Level; j++ {
			factor := &mat.Dense{}
			factor.Mul(ps[i][j], qs[i][j].T())
			Apq.Add(Apq, factor)
		}
	}

	Aqq := basic.ZeroMat()

	for i := 0; i < Level; i++ {
		for j := 0; j < Level; j++ {
			factor := &mat.Dense{}
			factor.Mul(qs[i][j], qs[i][j].T())
			Aqq.Add(Aqq, factor)
		}
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

	alpha := 0.75
	R.Scale(1-alpha, R)
	L.Scale(alpha, L)

	Shape := &mat.Dense{}
	Shape.Add(R, L)

	for i := 0; i < Level; i++ {
		for j := 0; j < Level; j++ {
			s.particles[i][j].pos.MulVec(Shape, qs[i][j].ColView(0))
			s.particles[i][j].pos.AddVec(s.particles[i][j].pos, t)
		}
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


