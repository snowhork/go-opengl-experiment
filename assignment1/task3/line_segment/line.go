package line_segment

import (
	"../basic"
	"../gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type line struct {
	p0, p1 *basic.Point
}

func newLine(p0, p1 *basic.Point) *line {
	vertShader, err := gfx.NewShaderFromFile("line_segment/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("line_segment/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()
	return &line{p0, p1}
}

func (l *line) color() (float32, float32, float32, float32) {
	return 1.0, 1.0, 0.0, 1.0
}

func (l *line) Draw() {
	points := make([]float32, 7*2, 7*2)

	points[0*7 + 0], points[0*7 + 1], points[0*7 + 2] = l.p0.Elements()
	points[0*7 + 3], points[0*7 + 4], points[0*7 + 5], points[0*7 + 6] = l.color()

	points[1*7 + 0], points[1*7 + 1], points[1*7 + 2] = l.p1.Elements()
	points[1*7 + 3], points[1*7 + 4], points[1*7 + 5], points[1*7 + 6] = l.color()

	VAO := makeVao(points)
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

