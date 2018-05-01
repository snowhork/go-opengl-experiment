package controller

import (
	"../basic"
	"github.com/go-gl/gl/v4.1-core/gl"

	"math"
	"math/rand"
	"github.com/lucasb-eyer/go-colorful"
)

type Controllable interface {
	Update()
	Draw()
}

var damping = float32(0.99)

type controller struct {
	beads []*bead
	beadsNum int

	Controllable
}

func NewController() *controller {
	beadsNum := 20

	beads := make([]*bead, beadsNum)

	for i := 0; i < beadsNum ; i++ {
		beads[i] = NewBead(&basic.Point{X: rand.Float32()*2.0-1.0, Y: rand.Float32()}, rand.Float32()*0.1+0.05, &colorful.Color{R: 0.7, G: 0.2})
	}

	return &controller{beads: beads, beadsNum: beadsNum}
}

func (c *controller) AddBead(x, y float32) {
	bead := NewBead(&basic.Point{X: x, Y: y}, rand.Float32()*0.1+0.05, &colorful.Color{R: 0.7, G: 0.2})
	c.beads = append(c.beads, bead)
	c.beadsNum += 1
}

func (c *controller) Update() {
	for i := 0; i < c.beadsNum ; i++ {
		c.beads[i].accelerate(0.01)
	}

	c.collide(false)

	for i := 0; i < c.beadsNum ; i++ {
		c.beads[i].borderCollide(false)
	}

	for i := 0; i < c.beadsNum ; i++ {
		c.beads[i].inertia()
	}

	c.collide(true)

	for i := 0; i < c.beadsNum ; i++ {
		c.beads[i].borderCollide(true)
	}

}

func (c *controller) collide(preserveImpulse bool) {
	for i := 0; i < c.beadsNum ; i++ {
		for j := i+1; j < c.beadsNum ; j++ {
			b1 := c.beads[i]
			b2 := c.beads[j]
			dir := b2.current.Sub(b1.current)

			if dir.Length() < b1.radius+b2.radius {
				d := b1.radius+b2.radius - dir.Length()

				v1 := b1.velocity()
				v2 := b2.velocity()

				b1.current = b1.current.Add(dir.Normalized().Mult(-d/2.0))
				b2.current = b2.current.Add(dir.Normalized().Mult(d/2.0))

				if preserveImpulse {
					b1Impuls := dir.Normalized().Mult(dir.Normalized().Product(v1)).Mult(damping)
					b2Impuls := dir.Normalized().Mult(dir.Normalized().Product(v2)).Mult(damping)

					b1.prev = b1.current.Sub(v1.Add(b2Impuls).Sub(b1Impuls))
					b2.prev = b2.current.Sub(v2.Add(b1Impuls).Sub(b2Impuls))
				}
			}
		}
	}
}

func (c *controller) Draw() {
	array := make([]float32, (VertexCount*3)*7*c.beadsNum)

	for j, b := range c.beads {
		offset := 21*VertexCount*j
		for i := 0; i < VertexCount; i++ {
			theta := math.Pi * 2.0 * float64(i) / VertexCount

			array[offset+(i*3)*7+0], array[offset+(i*3)*7+1], array[offset+(i*3)*7+2] = b.current.Elements()
			array[offset+(i*3)*7+3], array[offset+(i*3)*7+4], array[offset+(i*3)*7+5], array[offset+(i*3)*7+6] = b.colorRGBA()

			array[offset+(i*3+1)*7+0], array[offset+(i*3+1)*7+1], array[offset+(i*3+1)*7+2] = b.current.Add(
				&basic.Point{
					X: b.radius*float32(math.Cos(theta)),
					Y: b.radius*float32(math.Sin(theta)),
				}).Elements()
			array[offset+(i*3+1)*7+3], array[offset+(i*3+1)*7+4], array[offset+(i*3+1)*7+5], array[offset+(i*3+1)*7+6] = b.colorRGBA()

			theta2 := math.Pi * 2.0 * float64(i+1) / VertexCount
			array[offset+(i*3+2)*7+0], array[offset+(i*3+2)*7+1], array[offset+(i*3+2)*7+2] = b.current.Add(
				&basic.Point{
					X: b.radius*float32(math.Cos(theta2)),
					Y: b.radius*float32(math.Sin(theta2)),
				}).Elements()
			array[offset+(i*3+2)*7+3], array[offset+(i*3+2)*7+4], array[offset+(i*3+2)*7+5], array[offset+(i*3+2)*7+6] = b.colorRGBA()
		}
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	for j := range c.beads {
		gl.DrawArrays(gl.TRIANGLES, int32(3*7*VertexCount*j), 3*7*VertexCount)
	}
}

func makeVao(array []float32) uint32 {
	var vbo uint32

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(array), gl.Ptr(array), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 7*4, gl.PtrOffset(0))
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 7*4, gl.PtrOffset(3*4))

	return vao
}
