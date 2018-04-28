package line_segment

import (
	"../basic"
	"../gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
	"log"
)

type Bezier struct {
	p0, p1, p2, p3 *basic.Point
	q0, q1, q2, q3 *basic.Point
	b *bead
}

func NewBezier(p0, p1, p2, p3 *basic.Point) *Bezier {
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

	q0 := p0
	q1 := p0.Mult(-3).Add(p1.Mult(3))
	q2 := p0.Mult(3).Add(p1.Mult(-6)).Add(p2.Mult(3))
	q3 := p0.Mult(-1).Add(p1.Mult(3)).Add(p2.Mult(-3)).Add(p3)

	b := NewBead(p0, p0.Add(p0.Sub(p1).Normalized().Mult(0.01)))
	return &Bezier{p0, p1, p2, p3, q0, q1, q2, q3, b}
}

func (b *Bezier) Point(t float32) *basic.Point {
	return b.p0.Mult((1-t)*(1-t)*(1-t)).Add(b.p1.Mult(3*(1-t)*(1-t)*t)).Add(b.p2.Mult(3*(1-t)*t*t)).Add(b.p3.Mult(t*t*t))
}

func (b *Bezier) color() (float32, float32, float32, float32) {
	return 1.0, 1.0, 0.0, 1.0
}

func (b *Bezier) Draw() {
	VertexCount := 256
	points := make([]float32, 7*VertexCount, 7*VertexCount)

	for i := 0; i < VertexCount; i++ {
		t := float32(i)/float32(VertexCount-1)
		p := b.Point(t)

		points[i*7 + 0], points[i*7 + 1], points[i*7 + 2] = p.Elements()
		points[i*7 + 3], points[i*7 + 4], points[i*7 + 5], points[i*7 + 6] = b.color()
	}

	VAO := makeVao(points)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.LINE_STRIP, 0, int32(VertexCount))

	b.b.Draw()
}

func (b *Bezier) Update() {
	b.b.Update(0.01)
	c := b.b.current

	a0 := 3*b.q3.Product(b.q3)
	a1 := 5*b.q3.Product(b.q2)
	a2 := 4*b.q3.Product(b.q1) + 2*b.q2.Product(b.q2)
	a3 := 3*b.q2.Product(b.q1) + 3*b.q3.Product(b.q0.Sub(c))
	a4 := b.q1.Product(b.q1) + 2*b.q2.Product(b.q0.Sub(c))
	a5 := b.q1.Product(b.q0.Sub(c))


	s := NewSturm(float64(a0), float64(a1), float64(a2), float64(a3), float64(a4), float64(a5))
	roots := s.Root(0.0, 1.0, make([]float64, 0, 5))

	length := func(t float32) float32 {
		return b.Point(t).Sub(c).Product(b.Point(t).Sub(c))
	}

	minT := float32(0.0)
	if length(minT) > length(1.0) {
		minT = 1.0
	}

	for _, root := range roots {
		if length(minT) > length(float32(root)) {
			minT = float32(root)
		}
	}

	log.Println(minT)
	b.b.current = b.Point(minT)
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
