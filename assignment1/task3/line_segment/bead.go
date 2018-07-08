package line_segment

import (
	"math"
	"github.com/go-gl/gl/v4.1-core/gl"

	"../basic"
)

type bead struct {
	prev *basic.Point
	current *basic.Point
}

const VertexCount = 32

func newBead(prev, current *basic.Point) *bead {
	b := &bead{
		prev: prev,
		current: current,
	}

	return b
}


func (b *bead) color() (float32, float32, float32, float32) {
	return 1.0, 0.0, 1.0, 1.0
}

//
//func (b *bead) Update(dt float32) {
//	g := &basic.Point{Y: -9.8}
//	next := b.current.Add(b.current).Sub(b.prev).Add(g.Mult(dt*dt))
//	b.prev = b.current
//	b.current = next
//}

func (b *bead) Draw() {
	array := make([]float32, (VertexCount*3)*7)

	for i := 0; i < VertexCount; i++ {
		r := 0.02
		theta := math.Pi*2.0*float64(i)/ VertexCount

		array[(i*3)*7+0], array[(i*3)*7+1], array[(i*3)*7+2] = b.current.Elements()
		array[(i*3)*7+3], array[(i*3)*7+4], array[(i*3)*7+5], array[(i*3)*7+6] = b.color()

		array[(i*3+1)*7+0], array[(i*3+1)*7+1], array[(i*3+1)*7+2] = b.current.Add(
			&basic.Point{
				X: float32(r*math.Cos(theta)),
				Y: float32(r*math.Sin(theta)),
			}).Elements()
		array[(i*3+1)*7+3], array[(i*3+1)*7+4], array[(i*3+1)*7+5], array[(i*3+1)*7+6] = b.color()

		theta2 := math.Pi*2.0*float64(i+1)/ VertexCount
		array[(i*3+2)*7+0], array[(i*3+2)*7+1], array[(i*3+2)*7+2] = b.current.Add(
			&basic.Point{
				X: float32(r*math.Cos(theta2)),
				Y: float32(r*math.Sin(theta2)),
			}).Elements()
		array[(i*3+2)*7+3], array[(i*3+2)*7+4], array[(i*3+2)*7+5], array[(i*3+2)*7+6] = b.color()
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*7*VertexCount)
}
