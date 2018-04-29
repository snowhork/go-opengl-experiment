package line_segment

import (
	"math"
	"github.com/go-gl/gl/v4.1-core/gl"

	"../basic"
	"github.com/lucasb-eyer/go-colorful"
)

type bead struct {
	prev *basic.Point
	current *basic.Point
	color *colorful.Color
}

const VertexCount = 32
const radius = 0.02

func NewBead(prev, current *basic.Point, color *colorful.Color) *bead {
	b := &bead{
		prev: prev,
		current: current,
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


func (b *bead) Update(dt float32) {
	g := &basic.Point{Y: -9.8}
	next := b.current.Add(b.current).Sub(b.prev).Add(g.Mult(dt*dt))
	b.prev = b.current
	b.current = next
}

func (b *bead) Draw() {
	array := make([]float32, (VertexCount*3)*7)

	for i := 0; i < VertexCount; i++ {
		theta := math.Pi*2.0*float64(i)/ VertexCount

		array[(i*3)*7+0], array[(i*3)*7+1], array[(i*3)*7+2] = b.current.Elements()
		array[(i*3)*7+3], array[(i*3)*7+4], array[(i*3)*7+5], array[(i*3)*7+6] = b.colorRGBA()

		array[(i*3+1)*7+0], array[(i*3+1)*7+1], array[(i*3+1)*7+2] = b.current.Add(
			&basic.Point{
				X: float32(radius *math.Cos(theta)),
				Y: float32(radius *math.Sin(theta)),
			}).Elements()
		array[(i*3+1)*7+3], array[(i*3+1)*7+4], array[(i*3+1)*7+5], array[(i*3+1)*7+6] = b.colorRGBA()

		theta2 := math.Pi*2.0*float64(i+1)/ VertexCount
		array[(i*3+2)*7+0], array[(i*3+2)*7+1], array[(i*3+2)*7+2] = b.current.Add(
			&basic.Point{
				X: float32(radius *math.Cos(theta2)),
				Y: float32(radius *math.Sin(theta2)),
			}).Elements()
		array[(i*3+2)*7+3], array[(i*3+2)*7+4], array[(i*3+2)*7+5], array[(i*3+2)*7+6] = b.colorRGBA()
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*7*VertexCount)
}
