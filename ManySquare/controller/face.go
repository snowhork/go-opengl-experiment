package controller

import (
	"../basic"
	"math"
)

type face struct {
	p0 *basic.Point
	p1 *basic.Point
	depth float64
}

func (f *face) detectCollision(p *basic.Point) (*basic.Point, float64) {
	t := p.Sub(f.p0).Dot(f.p1.Sub(f.p0))/(f.p1.Sub(f.p0).Length2())

	if t < 0 || t > 1{
		return nil, 0
	}

	x := f.p0.Mult(t).Add(f.p1.Mult(1-t))
	d := x.Sub(p).Dot(f.normal())

	if d < 0 || d > f.depth {
		return nil, 0
	}

	return x, d
}

func (f *face) detectSinkV2(p *basic.Point) (*basic.Point, float64) {
	t := p.Sub(f.p0).Dot(f.p1.Sub(f.p0))/(f.p1.Sub(f.p0).Length2())

	t = math.Min(t, 1)
	t = math.Max(t, 0)

	x := f.p0.Mult(1-t).Add(f.p1.Mult(t))
	d := x.Sub(p).Dot(f.normal())

	if d < 0 || d > f.depth {
		return nil, 0
	}

	return x, d
}


func (f *face) detectCollisionV2(q0, q1 *basic.Point) (area float64, hitCenter *basic.Point, d float64) {
	u0, d0 := f.detectSinkV2(q0)
	u1, d1 := f.detectSinkV2(q1)

	if u0 == nil && u1 == nil {
		return 0, nil, 0
	}

	if u0 != nil && u1 != nil {
		return (d0*d1)*u0.Sub(u1).Length()/2, u0.Add(u1).Mult(0.5), math.Max(d0, d1)
	}

	aX := f.p0.X - f.p1.X
	aY := f.p0.Y - f.p1.Y

	bX := -(q0.X - q1.X)
	bY := -(q0.Y - q1.Y)

	cX := q1.X - f.p1.X
	cY := q1.Y - f.p1.Y

	det := aX*bY-bX*aY
	if det == 0 {
		return 0, nil, 0
		//panic("det must not be 0")
	}

	s, t := (bY*cX - bX*cY)/det, (-aY*cX + aX*cY)/det

	if !(s >= 0 && s <= 1 && t >= 0 && t <=1) {
		//log.Println(s,t)
		//panic("s,t must be in [0, 1]")
		return 0, nil, 0
	}

	crossPoint := f.p0.Mult(s).Add(f.p1.Mult(1-s))

	if u0 == nil {
		u0 = u1
		d0 = d1
	}

	return u0.Sub(crossPoint).Length()*d0/2, u0.Add(crossPoint).Mult(0.5), d0
}

func (f *face) normal() *basic.Point {
	return f.p1.Sub(f.p0).Rotation2D(math.Pi/2).Normalized()
}