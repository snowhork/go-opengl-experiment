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