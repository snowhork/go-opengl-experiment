package basic

import "math"

type Point struct {
	X,Y,Z float64
}

func NewPoint(X,Y,Z float64) *Point {
	return &Point{
		X: X,
		Y: Y,
		Z: Z,
	}
}

func (p *Point) Length() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y + p.Z*p.Z)
}

func (p *Point) Length2() float64 {
	return p.X*p.X + p.Y*p.Y + p.Z*p.Z
}

func (p *Point) Add(p2 *Point) *Point {
	return &Point{
		X: p.X+p2.X,
		Y: p.Y+p2.Y,
		Z: p.Z+p2.Z,
	}
}

func (p *Point) Sub(p2 *Point) *Point {
	return &Point{
		X: p.X-p2.X,
		Y: p.Y-p2.Y,
		Z: p.Z-p2.Z,
	}
}

func (p *Point) Mult(k float64) *Point {
	return &Point{
		X: p.X*k,
		Y: p.Y*k,
		Z: p.Z*k,
	}
}

func (p *Point) Dot(p2 *Point) float64 {
	return p.X*p2.X + p.Y*p2.Y + p.Z*p2.Z
}

func (p *Point) Cross(p2 *Point) *Point {
	return &Point{
		X: p.Y*p2.Z-p.Z*p2.Y,
		Y: p.Z*p2.X-p.X*p2.Z,
		Z: p.X*p2.Y-p.Y*p2.X,
	}
}

func (p *Point) Normalized() *Point {
	return p.Mult(1.0/p.Length())
}

func (p *Point) Elements() (float32,float32,float32)  {
	return float32(p.X), float32(p.Y), float32(p.Z)
}

func (p *Point) Rotation2D(theta float64) *Point {
	return &Point{
		X: math.Cos(theta)*p.X - math.Sin(theta)*p.Y,
		Y: math.Sin(theta)*p.X + math.Cos(theta)*p.Y,
		Z: p.Z,
	}
}

func Zero() *Point {
	return &Point{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}
}
