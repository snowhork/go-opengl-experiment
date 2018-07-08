package controller

import (
	"../gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Controllable interface {
	Update()
	Draw()
}

const (
	PenetrationEffect = 3.5e4
	VelocityEffect = 9e3
	//PenetrationEffect = 1400
	//VelocityEffect = 1000
)

type controller struct {
	s *square
	//walls []*face

	Controllable
}

func NewController() *controller {
	shaderinit()

	square := newSquare(0.0, 0.0)
	return &controller{s: square}
}


func (c *controller) Update() {
	c.s.update(0.01)
}

func (c *controller) Draw() {
	array := make([]float32, 3*7*2*(Level-1)*(Level-1))

	for i := 0; i < Level-1; i += 1 {
		for j := 0; j < Level-1; j += 1 {
			p0, p1, p2, p3 := c.s.particles[i][j], c.s.particles[i][j+1], c.s.particles[i+1][j+1], c.s.particles[i+1][j]

			points := make([][]*particle, 2)
			points[0] = []*particle{p0, p1, p2}
			points[1] = []*particle{p3, p0, p2}

			offset := (i*(Level-1)+j)*6*7

			for k, ps := range points {
				for l, p := range ps {
					array[offset+(k*3+l)*7+0], array[offset+(k*3+l)*7+1] = p.elements()
					array[offset+(k*3+l)*7+3], array[offset+(k*3+l)*7+4], array[offset+(k*3+l)*7+5], array[offset+(k*3+l)*7+6] = 0.0, 0.8, 0.8, 1.0
				}
			}
		}
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*7*2*(Level-1)*(Level-1))
}

func shaderinit()  {
	vertShader, err := gfx.NewShaderFromFile("shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()
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
