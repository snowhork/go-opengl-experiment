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

func (face *face) Balance() *basic.Point {
	return &basic.Point{
		X: (face.points[0].X + face.points[1].X + face.points[2].X)/3.0,
		Y: (face.points[0].Y + face.points[1].Y + face.points[2].Y)/3.0,
		Z: (face.points[0].Z + face.points[1].Z + face.points[2].Z)/3.0,
	}
}

func (input *face) Subdivide() []*face {
	SphereR := float32(math.Sqrt(PHI*(math.Sqrt(5))))

	p0 := input.points[0]
	p1 := input.points[1]
	p2 := input.points[2]

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

