package hex

import (
	"../gfx"
	"math"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Hex struct {
	points []float32
	program *gfx.Program
	count int
}

func NewHex() *Hex {
	vertShader, err := gfx.NewShaderFromFile("hex/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("hex/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()

	triangle := &Hex{
		program: program,
	}

	triangle.setPoint(0)

	return triangle
}

func (triangle *Hex) setPoint(delta float64) {
	sin := func(x float64) float32 { return float32(math.Sin(x)) }
	cos := func(x float64) float32 { return float32(math.Cos(x)) }
	Pi := math.Pi

	// x, y, z
	// r, g, b, a

	white := []float32{
		0.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 1.0,
	}

	red := []float32{
		cos(Pi/6.0 + delta), sin(Pi/6.0 + delta), 0.0,
		1.0, 0.0, 0.0, 1.0,
	}

	yellow := []float32{
		cos(Pi/2.0 + delta), sin(Pi/2.0 + delta), 0.0,
		1.0, 1.0, 0.0, 1.0,
	}

	green := []float32{
		cos(Pi*5.0/6.0 + delta), sin(Pi*5.0/6.0 + delta), 0.0,
		0.0, 1.0, 0.0, 1.0,
	}

	cyan := []float32{
		cos(Pi*7.0/6.0 + delta), sin(Pi*7.0/6.0 + delta), 0.0,
		0.0, 1.0, 1.0, 1.0,
	}

	blue := []float32{
		cos(Pi*3.0/2.0 + delta), sin(Pi*3.0/2.0 + delta), 0.0,
		0.0, 0.0, 1.0, 1.0,
	}

	magenta := []float32{
		cos(Pi*11.0/6.0 + delta), sin(Pi*11.0/6.0 + delta), 0.0,
		1.0, 0.0, 1.0, 1.0,
	}


	points := white
	points = append(points, red...)
	points = append(points, yellow...)
	points = append(points, white...)
	points = append(points, yellow...)
	points = append(points, green...)
	points = append(points, white...)
	points = append(points, green...)
	points = append(points, cyan...)
	points = append(points, white...)
	points = append(points, cyan...)
	points = append(points, blue...)
	points = append(points, white...)
	points = append(points, blue...)
	points = append(points, magenta...)
	points = append(points, white...)
	points = append(points, magenta...)
	points = append(points, red...)

	triangle.points = points
}

func (triangle *Hex) Update() {
	triangle.count += 1

	triangle.setPoint(float64(triangle.count)*0.03)

// do nothing
}

func (triangle *Hex) Draw() {
	//drawable := gfx.MakeVao(triangle.points)
	//gl.BindVertexArray(drawable)

	VAO := makeVao(triangle.points)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*6)
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
