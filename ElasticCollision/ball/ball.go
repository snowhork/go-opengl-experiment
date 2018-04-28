package ball

import (
	"../gfx"
	"math"
	"github.com/go-gl/gl/v4.1-core/gl"

	"../basic"
)

type Ball struct {
	prev *basic.Point
	current *basic.Point
}

const VertexCount = 32

func NewBall(pos *basic.Point) *Ball {
	vertShader, err := gfx.NewShaderFromFile("ball/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("ball/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()

	b := &Ball{
		prev: pos,
		current: pos,
	}

	return b
}


func (b *Ball) color() (float32, float32, float32, float32) {
	return 1.0, 0.0, 0.0, 1.0
}


func (b *Ball) Update(dt float32) {
	g := &basic.Point{Y: -9.8}
	next := b.current.Add(b.current).Sub(b.prev).Add(g.Mult(dt*dt))
	b.prev = b.current
	b.current = next

	if b.current.Y < -1.0 {

		b.prev.Y = -1.0 - (b.prev.Y + 1.0)
		b.current.Y = -1.0 - (b.current.Y + 1.0)
	}
}

func (b *Ball) Draw() {
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
