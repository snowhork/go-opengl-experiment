package controller

import (
	"../basic"
	"github.com/go-gl/gl/v4.1-core/gl"

	"math"
	"github.com/lucasb-eyer/go-colorful"
)

type Controllable interface {
	Update()
	Draw()
}

type controller struct {
	beads []*bead
	beadsNum int

	Controllable
}

func NewController() *controller {
	beadsNum := 5

	beads := make([]*bead, beadsNum)


	beads[0] = NewBead(&basic.Point{
		X: -Radius*2*float64(beadsNum)/2.0,
	}, -math.Pi/2.0, &colorful.Color{R: 0.7, G: 0.2})

	for i := 1; i < beadsNum ; i++ {
		pos := &basic.Point{
			X: -Radius*2*float64(beadsNum)/2.0+float64(i)*2*Radius,
		}

		beads[i] = NewBead(pos, 0, &colorful.Color{R: 0.7, G: 0.2})
	}

	return &controller{beads: beads, beadsNum: beadsNum}
}


func (c *controller) Update() {
	for i := 0; i < c.beadsNum; i++ {
		c.beads[i].accelerate(0.01)
	}

	c.collide()
}

func (c *controller) collide() {
	for i := 0; i < c.beadsNum ; i++ {
		for j := i+1; j < c.beadsNum ; j++ {
			b1 := c.beads[i]
			b2 := c.beads[j]
			dir := b2.pos().Sub(b1.pos())

			if dir.Length() < b1.radius+b2.radius {
				d := b1.radius+b2.radius - dir.Length()
				// push back
				b1Pos := b1.pos().Add(dir.Normalized().Mult(-d/2.0))
				b2Pos := b2.pos().Add(dir.Normalized().Mult(d/2.0))

				// fix theta for constraint
				b1.fixPosition(b1Pos)
				b2.fixPosition(b2Pos)

				// swap pTheta
				b1.pTheta, b2.pTheta = b2.pTheta, b1.pTheta
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

			array[offset+(i*3)*7+0], array[offset+(i*3)*7+1], array[offset+(i*3)*7+2] = b.pos().Elements()
			array[offset+(i*3)*7+3], array[offset+(i*3)*7+4], array[offset+(i*3)*7+5], array[offset+(i*3)*7+6] = b.colorRGBA()

			array[offset+(i*3+1)*7+0], array[offset+(i*3+1)*7+1], array[offset+(i*3+1)*7+2] = b.pos().Add(
				&basic.Point{
					X: b.radius*math.Cos(theta),
					Y: b.radius*math.Sin(theta),
				}).Elements()
			array[offset+(i*3+1)*7+3], array[offset+(i*3+1)*7+4], array[offset+(i*3+1)*7+5], array[offset+(i*3+1)*7+6] = b.colorRGBA()

			theta2 := math.Pi * 2.0 * float64(i+1) / VertexCount
			array[offset+(i*3+2)*7+0], array[offset+(i*3+2)*7+1], array[offset+(i*3+2)*7+2] = b.pos().Add(
				&basic.Point{
					X: b.radius*math.Cos(theta2),
					Y: b.radius*math.Sin(theta2),
				}).Elements()
			array[offset+(i*3+2)*7+3], array[offset+(i*3+2)*7+4], array[offset+(i*3+2)*7+5], array[offset+(i*3+2)*7+6] = b.colorRGBA()
		}
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	for j := range c.beads {
		gl.DrawArrays(gl.TRIANGLES, int32(3*7*VertexCount*j), 3*7*VertexCount)
	}

	c.DrawLine()
	c.DrawBar()
}

func (c *controller) DrawLine() {
	array := make([]float32, 2*7*c.beadsNum)

	for j, b := range c.beads {
		offset := j*14
		array[offset+0], array[offset+1], array[offset+2] = b.pos().Elements()
		array[offset+3], array[offset+4], array[offset+5], array[offset+6] = b.colorRGBA()

		array[offset+7+0], array[offset+7+1], array[offset+7+2] = b.center.Elements()
		array[offset+7+3], array[offset+7+4], array[offset+7+5], array[offset+7+6] = b.colorRGBA()
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	for j := range c.beads {
		gl.DrawArrays(gl.LINE_STRIP, int32(2*j), 2)
	}
}

func (c *controller) DrawBar() {
	array := make([]float32, 2*7)

	array[0+0], array[0+1] = -1.0, 0.0
	array[0+3], array[0+4], array[0+5], array[0+6] = 1.0, 1.0, 1.0, 1.0
	array[7+0], array[7+1] = 1.0, 0.0
	array[7+3], array[7+4], array[7+5], array[7+6] = 1.0, 1.0, 1.0, 1.0

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.LINE_STRIP, 0, 2)
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
