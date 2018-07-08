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
	DistanceThreshold = 0.5
	largeCount = 500
)
func NewGaraxy(smallCount int) *Galaxy {

	stars := make([]*star, largeCount+smallCount)

	for i := 0; i < largeCount ; i++  {
		r := rand.Float64()*0.5 + 0.000001
		theta := rand.Float64()*2*math.Pi
		eps := 0.5

		V := 0.012
		V = 0.0

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
		//log.Println(p)
		d := p.Sub(s.Current)
		force := d.Mult(G*s.mass*m/(math.Pow(d.Length(), 3)+ETA))
		//log.Println(force)

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
		//log.Println(other)
		return
	}
	if (n.xMax - n.xMin)/n.balance.Sub(s.Current).Length() <= DistanceThreshold {
		//calcforce(n.balance, n.mass)
	} else {
		for _, child := range n.children {
			search(child, s)
		}
	}
}


func (g *Galaxy) Update() {
	//for i, starFrom := range g.stars {
	//	for j, starTo := range g.stars {
	//		if i == j {
	//			continue
	//		}
	//		d := starFrom.Current.Sub(starTo.Current)
	//		force := d.Mult(G*starFrom.mass*starTo.mass/(math.Pow(d.Length(), 3)+ETA))
	//		//log.Println(force)
	//
	//		starTo.force = starTo.force.Add(force)
	//	}
	//}

	//log.Println(g.stars[0].force.Length(), g.stars[1].force.Length())

	//log.Println(g.stars[0].prev, g.stars[1].prev)

	root := g.Tree()

	for _, star := range g.stars {
		star.force = basic.Zero()
		Cnt = 0
		search(root, star)
		//log.Println(Cnt)
	}

	//log.Println(g.stars[0].force.Length(), g.stars[1].force.Length())
	//log.Println("----------------")
	//log.Println("current", g.stars[0].Current, g.stars[1].Current)
	//log.Println("force", g.stars[0].force, g.stars[1].force)
	//log.Println("mass", g.stars[0].mass, g.stars[1].mass)
	//log.Println("prev", g.stars[0].prev, g.stars[1].prev)


	for _, star := range g.stars {
		star.accelerate(dt)
	}
	//log.Println("next", g.stars[0].Current, g.stars[1].Current)

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

func (g *Galaxy) DrawDebug() {
	//nodes := g.Tree()
	//
	//points := make([]*basic.Point, 0)
	//children := make([]*node, 0)

	//n := nodes[400]
	//
	//points = append(points, n.balance)
	//
	//
	//beforeParent := n
	//currentParent := n.parent
	//
	//cnt := 0
	//depth := 0
	//for ; currentParent != nil; {
	//	depth += 1
	//	for _, child := range currentParent.children {
	//		if child == beforeParent {
	//			continue
	//		}
	//		cnt += 1
	//		children = append(children, child)
	//		points = append(points, child.balance)
	//	}
	//	beforeParent = currentParent
	//	currentParent = currentParent.parent
	//}
	//
	////log.Println(depth, cnt)
	//
	//array := make([]float32, (VertexCount*3)*7*len(g.stars))
	//
	//for i, p := range points {
	//	col := float32(i+1)/float32(len(points))
	//	col = 1.0
	//	array = append(array, pointArray(p, col)...)
	//}
	//
	//VAO := makeVao(array)
	//gl.BindVertexArray(VAO)
	//
	//for i := range g.stars {
	//	gl.DrawArrays(gl.TRIANGLES, int32((VertexCount*3)*7*i), VertexCount*3*7)
	//}
	//
	//balanceArray := make([]float32, 0)
	//
	//for _, c := range children {
	//	balanceArray = append(balanceArray, pointLineArray(c)...)
	//}
	//
	//VAO2 := makeVao(balanceArray)
	//gl.BindVertexArray(VAO2)
	//
	//for i := range children {
	//	gl.DrawArrays(gl.LINE_LOOP, int32(4*i), 4)
	//}

	//log.Println(len(points))

}

//func pointLineArray(n *node) []float32 {
//	array := make([]float32, 28)
//	array[0], array[1], array[2] = n.xMin, n.yMin, 0.0
//	array[3], array[4], array[5], array[6] = 1.0, 1.0, 1.0, 1.0
//
//	array[7*1+0], array[7*1+1], array[7*1+2] = n.xMin, n.yMax, 0.0
//	array[7*1+3], array[7*1+4], array[7*1+5], array[7*1+6] = 1.0, 1.0, 1.0, 1.0
//
//	array[7*2+0], array[7*2+1], array[7*2+2] = n.xMax, n.yMax, 0.0
//	array[7*2+3], array[7*2+4], array[7*2+5], array[7*2+6] = 1.0, 1.0, 1.0, 1.0
//
//	array[7*3+0], array[7*3+1], array[7*3+2] = n.xMax, n.yMin, 0.0
//	array[7*3+3], array[7*3+4], array[7*3+5], array[7*3+6] = 1.0, 1.0, 1.0, 1.0
//
//	return array
//}

func pointArray(p *basic.Point, col float32) []float32 {
	array := make([]float32, (VertexCount*3)*7)

	for i := 0; i < VertexCount; i++ {
		r := 0.02
		theta := math.Pi*2.0*float64(i)/ VertexCount

		array[(i*3)*7+0], array[(i*3)*7+1], array[(i*3)*7+2] = p.Elements()
		array[(i*3)*7+3], array[(i*3)*7+4], array[(i*3)*7+5], array[(i*3)*7+6] = col, float32(1.0), float32(1.0), float32(1.0)

		array[(i*3+1)*7+0], array[(i*3+1)*7+1], array[(i*3+1)*7+2] = p.Add(
			&basic.Point{
				X: r*math.Cos(theta),
				Y: r*math.Sin(theta),
			}).Elements()
		array[(i*3+1)*7+3], array[(i*3+1)*7+4], array[(i*3+1)*7+5], array[(i*3+1)*7+6] = col, float32(1.0), float32(1.0), float32(1.0)

		theta2 := math.Pi*2.0*float64(i+1)/ VertexCount
		array[(i*3+2)*7+0], array[(i*3+2)*7+1], array[(i*3+2)*7+2] = p.Add(
			&basic.Point{
				X: r*math.Cos(theta2),
				Y: r*math.Sin(theta2),
			}).Elements()
		array[(i*3+2)*7+3], array[(i*3+2)*7+4], array[(i*3+2)*7+5], array[(i*3+2)*7+6] = col, float32(1.0), float32(1.0), float32(1.0)
	}

	return array

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
