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

	x := f.p0.Mult(1-t).Add(f.p1.Mult(t))
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


func (f *face) detectCollisionV2(q0, q1 *basic.Point) (u0 *basic.Point, u1 *basic.Point, d float64) {
	u0, d = f.detectSinkV2(q0)
	if u0 == nil {
		u0, d = f.detectSinkV2(q1)
		if u0 == nil {
			return nil, nil, 0
		}
	}

	aX := f.p0.X - f.p1.X
	aY := f.p0.Y - f.p1.Y

	bX := -(q0.X - q1.X)
	bY := -(q0.Y - q1.Y)

	cX := q1.X - f.p1.X
	cY := q1.Y - f.p1.Y

	det := aX*bY-bX*aY
	if det == 0 {
		return nil, nil, 0
	}

	s, t := (bY*cX - bX*cY)/det, (-aY*cX + aX*cY)/det

	if !(s >= 0 && s <= 1 && t >= 0 && t <=1) {
		return nil, nil, 0
	}

	u1 = f.p0.Mult(s).Add(f.p1.Mult(1-s))

	return
}

func (f *face) normal() *basic.Point {
	return f.p1.Sub(f.p0).Rotation2D(math.Pi/2).Normalized()
}