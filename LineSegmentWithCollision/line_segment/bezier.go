package line_segment

import (
	"../basic"
	"../gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/lucasb-eyer/go-colorful"
)

type Bezier struct {
	p0, p1, p2, p3 *basic.Point
	q0, q1, q2, q3 *basic.Point
	b1, b2 *bead
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

	b1 := NewBead(p0, p0.Add(p0.Sub(p1).Normalized().Mult(0.01)), &colorful.Color{1.0, 0.0, 0.0})
	b2 := NewBead(p3, p3.Add(p3.Sub(p2).Normalized().Mult(0.01)), &colorful.Color{0.0, 0.0, 1.0})
	return &Bezier{p0, p1, p2, p3, q0, q1, q2, q3, b1, b2}
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

	b.b1.Draw()
	b.b2.Draw()
}

func (bez *Bezier) modify(b *bead) {
	c := b.current

	a0 := 3*bez.q3.Product(bez.q3)
	a1 := 5*bez.q3.Product(bez.q2)
	a2 := 4*bez.q3.Product(bez.q1) + 2*bez.q2.Product(bez.q2)
	a3 := 3*bez.q2.Product(bez.q1) + 3*bez.q3.Product(bez.q0.Sub(c))
	a4 := bez.q1.Product(bez.q1) + 2*bez.q2.Product(bez.q0.Sub(c))
	a5 := bez.q1.Product(bez.q0.Sub(c))


	s := NewSturm(float64(a0), float64(a1), float64(a2), float64(a3), float64(a4), float64(a5))
	roots := s.Root(0.0, 1.0, make([]float64, 0, 5))

	length := func(t float32) float32 {
		return bez.Point(t).Sub(c).Product(bez.Point(t).Sub(c))
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

	b.current = bez.Point(minT)

}

func (b *Bezier) Update() {
	b.b1.Update(0.01)
	b.modify(b.b1)
	b.b2.Update(0.01)
	b.modify(b.b2)

	dir := b.b2.current.Sub(b.b1.current)

	if dir.Length() <= radius*2 {
		d := radius*2 - dir.Length()
		b.b1.current = b.b1.current.Add(dir.Normalized().Mult(-d/2.0))
		b.b2.current = b.b2.current.Add(dir.Normalized().Mult(d/2.0))

		b1Impuls := dir.Normalized().Mult(dir.Normalized().Product(b.b1.velocity()))
		b2Impuls := dir.Normalized().Mult(dir.Normalized().Product(b.b2.velocity()))

		b.b1.prev = b.b1.current.Sub(b.b1.velocity().Add(b2Impuls).Sub(b1Impuls))
		b.b2.prev = b.b2.current.Sub(b.b2.velocity().Add(b1Impuls).Sub(b2Impuls))
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
