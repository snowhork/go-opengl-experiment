package controller

import (
	"../basic"
	"../gfx"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lucasb-eyer/go-colorful"
)

type bead struct {
	prev *basic.Point
	current *basic.Point
	radius float32
	color *colorful.Color
}

const VertexCount = 32

func NewBead(point *basic.Point, radius float32, color *colorful.Color) *bead {
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
		prev: point,
		current: point,
		radius: radius,
		color: color,
	}
	return b
}


func (b *bead) colorRGBA() (float32, float32, float32, float32) {
	return float32(b.color.R), float32(b.color.G), float32(b.color.B), 1.0
}

func (b *bead) velocity() *basic.Point {
	return b.current.Sub(b.prev)
}

func (b *bead) accelerate(dt float32) {
	g := &basic.Point{Y: -9.8}
	b.current = b.current.Add(g.Mult(dt * dt))
}

func (b *bead) inertia() {
	next := b.current.Add(b.velocity())
	b.prev = b.current
	b.current = next
}

func (b *bead) borderCollide(preservable bool) {
	v := b.velocity()

	if b.current.Y - b.radius < -1.0 {
		b.current.Y = -1.0 + b.radius

		if preservable {
			b.prev.Y = b.current.Y + (v.Y)
		}
	}

	if b.current.Y + b.radius > 1.0 {
		b.current.Y = 1.0 - b.radius

		if preservable {
			b.prev.Y = b.current.Y + (v.Y)
		}
	}

	if b.current.X  - b.radius < -1.0 {
		b.current.X = -1.0 + b.radius

		if preservable {
			b.prev.X = b.current.X + (v.X)
		}
	}

	if b.current.X + b.radius > 1.0 {
		b.current.X = 1.0 - b.radius

		if preservable {
			b.prev.X = b.current.X + (v.X)
		}
	}
}