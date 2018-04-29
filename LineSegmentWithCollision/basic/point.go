package basic

import "math"

type Point struct {
	X,Y,Z float32
}

func NewPoint(X,Y,Z float32) *Point {
	return &Point{
		X: X,
		Y: Y,
		Z: Z,
	}
}

func (p *Point) Length() float32 {
	return float32(math.Sqrt(float64(p.X*p.X + p.Y*p.Y + p.Z*p.Z)))
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
func (p *Point) Mult(k float32) *Point {
	return &Point{
		X: p.X*k,
		Y: p.Y*k,
		Z: p.Z*k,
	}
}

func (p *Point) Normalized() *Point {
	return p.Mult(1.0/p.Length())
}

func (p *Point) Elements() (float32,float32,float32)  {
	return p.X, p.Y, p.Z
}

func Zero() *Point {
	return &Point{
		X: 0.0,
		Y: 0.0,
		Z: 0.0,
	}
}

func (p *Point) Product(p2 *Point) float32 {
	return p.X*p2.X + p.Y*p2.Y + p.Z*p2.Z
}
