package fire_flower

import (
	"../gfx"
	"../basic"
	"math"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lucasb-eyer/go-colorful"
)

type FireFlower struct {
	lines []line
	faces []*face
	program *gfx.Program
	count int
	vertexCount int
}

func NewFireFlower(
	position *basic.Point,
	speed float32,
	divideDepth int,
	vertexCount int) *FireFlower {

	phi := float32((1+math.Sqrt2)/2.0)

	points := []*basic.Point{
		basic.NewPoint(1, phi, 0),
		basic.NewPoint(1, -phi, 0),
		basic.NewPoint(-1, -phi, 0),
		basic.NewPoint(-1, phi, 0),
		basic.NewPoint(0, 1, phi),
		basic.NewPoint(0, 1, -phi),
		basic.NewPoint(0, -1, -phi),
		basic.NewPoint(0, -1, phi),
		basic.NewPoint(phi, 0, 1),
		basic.NewPoint(-phi, 0, 1),
		basic.NewPoint(-phi, 0, -1),
		basic.NewPoint(phi, 0, -1),
	}

	baseFace := []*face{
		newFace(points[0], points[3], points[5]),
		newFace(points[3], points[10], points[5]),
		newFace(points[0], points[5], points[11]),
		newFace(points[5], points[6], points[11]),
		newFace(points[5], points[10], points[6]),
		newFace(points[11], points[6], points[1]),
		newFace(points[6], points[10], points[2]),
		newFace(points[6], points[2], points[1]),

		newFace(points[3], points[0], points[4]),
		newFace(points[3], points[4], points[9]),
		newFace(points[0], points[8], points[4]),
		newFace(points[4], points[7], points[9]),
		newFace(points[4], points[8], points[7]),
		newFace(points[9], points[7], points[2]),
		newFace(points[7], points[8], points[1]),
		newFace(points[7], points[1], points[2]),

		newFace(points[3], points[9], points[10]),
		newFace(points[0], points[11], points[8]),

		newFace(points[2], points[10], points[9]),
		newFace(points[1], points[8], points[11]),
	}

	faces := baseFace

	for depth := 0; depth < divideDepth; depth++ {
		nextFace := make([]*face, 0, len(faces)*4)
		for _, face := range faces {
			nextFace = append(nextFace, face.Subdivide()...)
		}
		faces = nextFace
	}

	lines := make([]line, len(faces), len(faces))

	for i := 0; i < len(faces); i++ {
		lines[i] = *newLine(
			position,
			faces[i].Balance().Mult(speed),
			&colorful.Color{1.0, 0.0, 0.0},
			1.0,
			vertexCount)
	}

	vertShader, err := gfx.NewShaderFromFile("fire_flower/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("fire_flower/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()

	flower := &FireFlower{
		lines: lines,
		faces: faces,
		program: program,
		vertexCount: vertexCount,
	}
	return flower
}


func (flower *FireFlower) Update() {
	for i := range flower.lines {
		flower.lines[i].update()
	}
}

func (flower *FireFlower) Draw() {
	points := make([]float32, 7*len(flower.lines)*flower.vertexCount, 7*len(flower.lines)*flower.vertexCount)

	for i := range flower.lines {
		for j := 0; j < 7*flower.vertexCount; j++ {
			points[i*7*flower.vertexCount + j] = flower.lines[i].points[j]
		}
	}

	VAO := makeVao(points)
	gl.BindVertexArray(VAO)

	for i := range flower.lines {
		gl.DrawArrays(gl.LINE_STRIP, int32(i*flower.vertexCount), int32(flower.vertexCount))
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
