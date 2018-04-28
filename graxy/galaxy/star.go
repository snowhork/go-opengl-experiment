package galaxy

import (
	"../basic"
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

type star struct {
	Current *basic.Point
	prev *basic.Point
	mass float32
	force *basic.Point
}

const VertexCount = 10

func newStar(p *basic.Point, v *basic.Point, mass float32) *star {
	return &star{
		Current: p,
		prev: p.Add(v.Mult(-1)),
		mass: mass,
		force: basic.Zero(),
	}
}

func (s *star) accelerate(delta float32) {
	next := s.Current.Add(s.Current.Sub(s.prev)).Add(s.force.Mult(delta*delta/s.mass))
	next.Z = 0
	s.prev = s.Current
	s.Current = next

	s.force = basic.Zero()
}

func (s *star) color() (float32,float32,float32,float32) {
	c := colorful.Hsv(240.0/0.50001*float64(1.0/s.mass), 1.0, 1.0)
	return float32(c.R),float32(c.G),float32(c.B), 1.0
}

func (s *star) array() []float32 {
	array := make([]float32, (VertexCount*3)*7)

	for i := 0; i < VertexCount; i++ {
		r := 0.01
		theta := math.Pi*2.0*float64(i)/ VertexCount

		array[(i*3)*7+0], array[(i*3)*7+1], array[(i*3)*7+2] = s.Current.Elements()
		array[(i*3)*7+3], array[(i*3)*7+4], array[(i*3)*7+5], array[(i*3)*7+6] = s.color()

		array[(i*3+1)*7+0], array[(i*3+1)*7+1], array[(i*3+1)*7+2] = s.Current.Add(
			&basic.Point{
				X: float32(r*math.Cos(theta)),
				Y: float32(r*math.Sin(theta)),
			}).Elements()
		array[(i*3+1)*7+3], array[(i*3+1)*7+4], array[(i*3+1)*7+5], array[(i*3+1)*7+6] = s.color()

		theta2 := math.Pi*2.0*float64(i+1)/ VertexCount
		array[(i*3+2)*7+0], array[(i*3+2)*7+1], array[(i*3+2)*7+2] = s.Current.Add(
			&basic.Point{
				X: float32(r*math.Cos(theta2)),
				Y: float32(r*math.Sin(theta2)),
			}).Elements()
		array[(i*3+2)*7+3], array[(i*3+2)*7+4], array[(i*3+2)*7+5], array[(i*3+2)*7+6] = s.color()
	}

	return array
}