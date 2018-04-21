package fire_flower

import (
	"../basic"
	"math"
)
const PHI = (1+math.Sqrt2)/2.0

type face struct {
	points []*basic.Point
}

func newFace(p1,p2,p3 *basic.Point) *face {
	return &face{points: []*basic.Point{p1, p2, p3}}
}

func (f *face) Balance() *basic.Point {
	return f.points[0].Add(f.points[1]).Add(f.points[2]).Mult(1.0/3.0)
}

func (f *face) Subdivide() []*face {
	p0 := f.points[0]
	p1 := f.points[1]
	p2 := f.points[2]

	SphereR := p0.Length()

	p01 := p0.Add(p1).Normalized().Mult(SphereR)
	p12 := p1.Add(p2).Normalized().Mult(SphereR)
	p20 := p2.Add(p0).Normalized().Mult(SphereR)

	faces := []*face{
		{points: []*basic.Point{p0, p01, p1}},
		{points: []*basic.Point{p1, p12, p2}},
		{points: []*basic.Point{p2, p20, p0}},
		{points: []*basic.Point{p01, p12, p20}},
	}
	return faces
}

