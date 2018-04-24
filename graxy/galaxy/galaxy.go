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
	G = 0.05
	ETA = 0.00001
	GFACTOR = 3

	DISTANCE_FACTOR = 1

	MINMASS = 1e4
	MAXMASS = 1e10

)
func NewGaraxy(largeCount int, smallCount int) *Galaxy {

	stars := make([]*star, largeCount+smallCount)

	//randomLarge := func() float32 {
	//	return float32(rand.NormFloat64()*0.4 - 0.5)
	//}
	//
	//randomSmall := func() float32 {
	//	return float32(rand.NormFloat64()*0.2 + 0.5)
	//}


	for i := 0; i < largeCount ; i++  {
		r := rand.Float64()*0.5 + 0.0000001
		theta := rand.Float64()*2*math.Pi

		//eps := rand.Float64()*0.1 + 0.95
		eps := 0.5

		V := 0.003

		p := &basic.Point{
			X: float32(r*math.Cos(theta)),
			Y: float32(eps*r*math.Sin(theta))}
		v := &basic.Point{
			X: -float32(V*r*math.Sin(theta)),
			Y: float32(V*r*math.Cos(theta))}
		stars[i] = newStar(p, v,
			1.0/float32(r))
	}

	//for i := largeCount; i < largeCount+smallCount ; i++  {
	//	v := basic.Point{X: float32(rand.NormFloat64()*0.2), Y: float32(rand.NormFloat64()*0.2), Z: 0.0}
	//	stars[i] = newStar(v.Add(&basic.Point{X: 0.5, Y: 0.5}), v.Mult(0.01), 1.0)
	//}


	//stars[largeCount-3] = newStar(&basic.Point{
	//	X: float32(0.6*math.Cos(-math.Pi/6.0)),
	//	Y: float32(0.6*math.Sin(-math.Pi/6.0)),
	//	Z: 0.0},
	//	&basic.Point{},
	//	1e12)
	//
	//stars[largeCount-2] = newStar(&basic.Point{
	//	X: float32(0.6*math.Cos(math.Pi/2.0)),
	//	Y: float32(0.6*math.Sin(math.Pi/2.0)),
	//	Z: 0.0},
	//	&basic.Point{},
	//	1e12)

	//stars[largeCount-1] = newStar(&basic.Point{
		//X: float32(0.0*math.Cos(math.Pi*2.0/3.0)),
		//Y: float32(0.0*math.Sin(math.Pi*2.0/3.0)),
		//Z: 0.0},
		//&basic.Point{},
		//1e11)

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
			force := d.Mult(G*starFrom.mass*starTo.mass/(float32(math.Pow(float64(d.Length())*DISTANCE_FACTOR, GFACTOR))+ETA))
			starTo.force = starTo.force.Add(force)
		}
	}
	for _, star := range g.stars {
		star.accelerate(0.0001)
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
