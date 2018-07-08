package galaxy

import (
	"../gfx"
	"../basic"
	"github.com/go-gl/gl/v4.1-core/gl"
	"math/rand"
	"math"
)

type Galaxy struct {
	stars []*star
}


const (
	G = 0.00001
	ETA = 0.0001
	dt = 0.0001
	DistanceThreshold = 0.5
	largeCount = 1000
)

func NewGaraxy(smallCount int) *Galaxy {

	stars := make([]*star, largeCount+smallCount)

	for i := 0; i < largeCount ; i++  {
		r := rand.Float64()*0.5 + 0.000001
		theta := rand.Float64()*2*math.Pi
		eps := 0.5

		V := 0.013

		p := &basic.Point{
			X: r*math.Cos(theta),
			Y: eps*r*math.Sin(theta)}
		v := &basic.Point{
			X: -eps*V*r*math.Sin(theta),
			Y: V*r*math.Cos(theta)}
		stars[i] = newStar(p, v,
			1.0/(r*r), i)
	}

	vertShader, err := gfx.NewShaderFromFile("galaxy/shaders/basic.vert", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragShader, err := gfx.NewShaderFromFile("galaxy/shaders/basic.frag", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	program, err := gfx.NewProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}

	program.Use()

	return &Galaxy{
		stars: stars,
	}
}

var Cnt = 0

func search(n *node, s *star) {
	calcforce := func(p *basic.Point, m float64) {
		d := p.Sub(s.Current)
		force := d.Mult(G*s.mass*m/(math.Pow(d.Length(), 3)+ETA))

		s.force = s.force.Add(force)

		Cnt += 1
	}

	if len(n.stars) == 0 {
		return
	}
	if len(n.stars) == 1 {
		if s.number == n.stars[0].number {
			return
		}
		other := n.stars[0]
		calcforce(other.Current, other.mass)
		return
	}
	if (n.xMax - n.xMin)/n.balance.Sub(s.Current).Length() <= DistanceThreshold {
		calcforce(n.balance, n.mass)
	} else {
		for _, child := range n.children {
			search(child, s)
		}
	}
}


func (g *Galaxy) Update() {
	root := g.Tree()

	for _, star := range g.stars {
		star.force = basic.Zero()
		Cnt = 0
		search(root, star)
	}

	for _, star := range g.stars {
		star.accelerate(dt)
	}
}



func (g *Galaxy) Draw() {

	points := make([]float32, (VertexCount*3)*7*len(g.stars))


	for _, star := range g.stars {
		points = append(points, star.array()...)
	}

	VAO := makeVao(points)
	gl.BindVertexArray(VAO)

	for i := range g.stars {
		gl.DrawArrays(gl.TRIANGLES, int32((VertexCount*3)*7*i), VertexCount*3*7)
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
