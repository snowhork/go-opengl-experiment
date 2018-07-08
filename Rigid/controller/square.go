package controller

import (
	"github.com/lucasb-eyer/go-colorful"
	"../basic"
	"math"
)

type square struct {
	position *basic.Point
	velocity *basic.Point

	angle *basic.Point
	angularVelocity *basic.Point

	force *basic.Point
	torque *basic.Point

	radius float64
	color *colorful.Color

	mass          float64
	momentInteria float64
}

func newSquare(pos *basic.Point) *square {
	return &square{
		position: pos,
		velocity: basic.Zero(),

		angle: &basic.Point{Z: 1.0},
		angularVelocity: &basic.Point{Z: 1.0},

		force: basic.Zero(),
		torque: basic.Zero(),

		radius: 0.2,
		color: &colorful.Color{R: 0.8},

		mass:          1.0,
		momentInteria: 0.1,
	}
}

func (s *square) points() (p0, p1, p2, p3 *basic.Point) {

	r := &basic.Point{
		X: s.radius*math.Cos(s.angle.Z + math.Pi/4.0),
		Y: s.radius*math.Sin(s.angle.Z + math.Pi/4.0),
	}

	p0 = s.position.Add(r)
	p1 = s.position.Add(r.Rotation2D(math.Pi/2.0))
	p2 = s.position.Add(r.Rotation2D(math.Pi))
	p3 = s.position.Add(r.Rotation2D(math.Pi*3.0/2.0))

	return p0, p1, p2, p3
}

func (s *square) faces() (f0, f1, f2, f3 *face) {
	p0, p1, p2, p3 := s.points()
	depth := 0.02

	f0 = &face{p0: p0, p1: p1, depth: depth}
	f1 = &face{p0: p1, p1: p2, depth: depth}
	f2 = &face{p0: p2, p1: p3, depth: depth}
	f3 = &face{p0: p3, p1: p0, depth: depth}

	return f0, f1, f2, f3
}

func (s *square) addFource(f *basic.Point) {
	s.force = s.force.Add(f)
}

func (s *square) addTorqueToPoint(f *basic.Point, p *basic.Point) {
	s.torque = s.torque.Add(p.Sub(s.position).Cross(f))
}

func (s *square) addTorque(t *basic.Point) {
	s.torque = s.torque.Add(t)
}

func (s *square) accelerate(dt float64) {
	next := s.position.Add(s.velocity).Add(s.force.Mult(dt*dt/s.mass))
	s.velocity = next.Sub(s.position)
	s.position = next

	s.force = basic.Zero()

	s.angularVelocity = s.angularVelocity.Add(s.torque.Mult(dt/s.momentInteria))
	s.angle.Z += s.angularVelocity.Z*dt

	s.torque = basic.Zero()
}

func (s *square) colorElements() (float32, float32, float32, float32) {
	return float32(s.color.R), float32(s.color.G), float32(s.color.B), 1.0
}