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
	dt = 0.001
)
func NewGaraxy(largeCount int, smallCount int) *Galaxy {

	stars := make([]*star, largeCount+smallCount)

	for i := 0; i < largeCount ; i++  {
		r := rand.Float64()*0.5 + 0.000001
		theta := rand.Float64()*2*math.Pi
		eps := 1.0

		V := 0.02

		p := &basic.Point{
			X: float32(r*math.Cos(theta)),
			Y: float32(eps*r*math.Sin(theta))}
		v := &basic.Point{
			X: -float32(eps*V*r*math.Sin(theta)),
			Y: float32(V*r*math.Cos(theta))}
		stars[i] = newStar(p, v,
			1.0/float32(r*r))
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


func (g *Galaxy) Update() {
	for i, starFrom := range g.stars {
		for j, starTo := range g.stars {
			if i == j {
				continue
			}
			d := starFrom.Current.Sub(starTo.Current)
			force := d.Mult(G*starFrom.mass*starTo.mass/(float32(math.Pow(float64(d.Length()), 3))+ETA))
			starTo.force = starTo.force.Add(force)
		}
	}

	//nodes := g.Tree()
	//
	//for _, node := range nodes {
	//	star := node.stars[0]
	//	beforeParent := node
	//	currentParent := node.parent
	//	i := 0
	//	j := 0
	//	for ; currentParent != nil; {
	//		j += 1
	//		for _, child := range currentParent.children {
	//			if child == beforeParent {
	//				i += 1
	//				continue
	//			}
	//			d := star.Current.Sub(child.balance)
	//			force := d.Mult(G*star.mass*child.mass/(float32(math.Pow(float64(d.Length()), 3))+ETA))
	//			star.force = star.force.Add(force)
	//		}
	//		beforeParent = currentParent
	//		currentParent = currentParent.parent
	//	}
	//
	//	//log.Println(i, j)
	//}

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
