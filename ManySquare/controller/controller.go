package controller

import (
	"../basic"
	"../gfx"
	"github.com/go-gl/gl/v4.1-core/gl"
	"math"
)

type Controllable interface {
	Update()
	Draw()
}

const (
	PenetrationEffect = 4e3
	PenetrationEffectConst = 14
	VelocityEffect = 2e3
)

type controller struct {
	squares []*square
	walls []*face

	Controllable
}

func NewController() *controller {
	shaderinit()

	walls := make([]*face, 4)

	squares := make([]*square, 0)
	squares = append(squares, newSquare(basic.NewPoint(0.5, 0.0, 0.0)))

	p0 := &basic.Point{X: -1.0, Y: -1.0}
	p1 := &basic.Point{X: 1.0, Y: -1.0}

	for i := 0; i < 4; i++ {
		walls[i] = &face{
			p0: p0.Rotation2D(float64(i)*math.Pi/2.0),
			p1: p1.Rotation2D(float64(i)*math.Pi/2.0),
			depth: 0.4,
		}
	}

	return &controller{squares: squares, walls: walls}
}

func (c *controller) Update() {
	for _, s := range c.squares {
		s.force = &basic.Point{Y: -9.8*s.mass}
	}

	c.collideV2()

	for _, s := range c.squares {
		s.accelerate(0.01)
	}
}

func (c *controller) AddSquare(x, y float64) {
	c.squares = append(c.squares, newSquare(basic.NewPoint(x, y, 0.0)))
}

func (c *controller) collideV2() {
	for i := range c.squares {
		iSquare := c.squares[i]
		iFaces := make([]*face, 4)
		iFaces[0], iFaces[1], iFaces[2], iFaces[3] = iSquare.faces()
		center := iSquare.position

		for _, iFace := range iFaces {
			for _, wall := range c.walls {
				area, hitCenter, _ := wall.detectCollisionV2(iFace.p0, iFace.p1)
				if area != 0 {
					velocity := iSquare.velocity

					penaltyForce := wall.normal().Mult(PenetrationEffectConst+PenetrationEffect*area - VelocityEffect*velocity.Dot(wall.normal()))

					penaltyTorque := (hitCenter.Sub(center)).Cross(wall.normal())
					penaltyTorque = penaltyTorque.Mult(PenetrationEffectConst+PenetrationEffect*area)

					iSquare.addFource(penaltyForce)
					iSquare.addTorque(penaltyTorque)
				}
			}
			for j := i + 1; j < len(c.squares); j++ {
				jSquare := c.squares[j]
				jFaces := make([]*face, 4)
				jFaces[0], jFaces[1], jFaces[2], jFaces[3] = jSquare.faces()


				for _, jFace := range jFaces {
					area, hitCenter, _ := jFace.detectCollisionV2(iFace.p0, iFace.p1)
					if area != 0 {
						velocity := iSquare.velocity.Sub(jSquare.velocity)

						penaltyForce := jFace.normal().Mult(PenetrationEffectConst+PenetrationEffect*area - VelocityEffect*velocity.Sub(jSquare.velocity).Dot(jFace.normal()))

						penaltyTorque := (hitCenter.Sub(center)).Cross(jFace.normal())
						penaltyTorque = penaltyTorque.Mult(PenetrationEffect*area)

						iSquare.addFource(penaltyForce)
						iSquare.addTorque(penaltyTorque)

						jpenaltyTorque := jFace.normal().Cross(hitCenter.Sub(jSquare.position))
						jpenaltyTorque = jpenaltyTorque.Mult(PenetrationEffect*area)


						jSquare.addFource(penaltyForce.Mult(-1))
						jSquare.addTorque(jpenaltyTorque)
					}
				}
			}
		}
	}
}

func (c *controller) Draw() {
	array := make([]float32, 3*7*2*len(c.squares))

	for k, square := range c.squares {

		p0, p1, p2, p3 := square.points()

		points := make([][]*basic.Point, 2)
		points[0] = []*basic.Point{p0, p1, p2}
		points[1] = []*basic.Point{p3, p0, p2}

		offset := 2*3*7*k

		for i, ps := range points {
			for j, p := range ps {
				array[offset+(i*3+j)*7+0], array[offset+(i*3+j)*7+1], array[offset+(i*3+j)*7+2] = p.Elements()
				array[offset+(i*3+j)*7+3], array[offset+(i*3+j)*7+4], array[offset+(i*3+j)*7+5], array[offset+(i*3+j)*7+6] = square.colorElements()
			}
		}
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*7*2*int32(len(c.squares)))
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
