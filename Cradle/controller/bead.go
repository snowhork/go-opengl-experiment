package controller

import (
	"../basic"
	"../gfx"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lucasb-eyer/go-colorful"
	"math"
)

type bead struct {
	pTheta float64
	theta  float64
	center *basic.Point
	radius float64
	l float64
	m float64
	color  *colorful.Color
}

const VertexCount = 32
const Radius = 0.05

func NewBead(point *basic.Point, theta float64, color *colorful.Color) *bead {
	vertShader, err := gfx.NewShaderFromFile("controller/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("controller/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()

	b := &bead{
		pTheta: 0.0,
		theta:  theta,
		radius: Radius,
		l: 0.5,
		m: 1.0,
		center: point,
		color:  color,
	}
	return b
}


func (b *bead) colorRGBA() (float32, float32, float32, float32) {
	return float32(b.color.R), float32(b.color.G), float32(b.color.B), 1.0
}


func (b *bead) accelerate(dt float64) {
	b.pTheta += dt*(-b.m*9.8*b.l)*math.Sin(b.theta)
	b.theta += dt*b.pTheta/(b.m*b.l*b.l)
}

func (b *bead) pos() *basic.Point {
	return &basic.Point{
		X: b.center.X+b.l*math.Sin(b.theta),
		Y: b.center.Y-b.l*math.Cos(b.theta),
	}
}

func (b *bead) velocity() *basic.Point {
	k := b.pTheta/b.m

	return &basic.Point{
		X: b.center.X+k*math.Cos(b.theta),
		Y: b.center.X+k*math.Sin(b.theta),
	}
}

func (b *bead) fixPosition(pos *basic.Point) {
	d := pos.Sub(b.center).Normalized()
	cos := d.Product(basic.Down())
	theta := math.Acos(cos)
	if d.X > 0 {
		b.theta = theta
	} else {
		b.theta = -theta
	}
}

func (b *bead) fixVelocity(v *basic.Point) {
	t := &basic.Point{X: math.Cos(b.theta), Y: math.Sin(b.theta)}

	speed := v.Product(t)

	b.pTheta = b.m*speed
}