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
	PenetrationEffect = 3.5e4
	VelocityEffect = 9e3
	//PenetrationEffect = 1400
	//VelocityEffect = 1000
)

type controller struct {
	s *square
	walls []*face

	Controllable
}

func NewController() *controller {
	shaderinit()

	walls := make([]*face, 4)

	p0 := &basic.Point{X: -1.0, Y: -1.0}
	p1 := &basic.Point{X: 1.0, Y: -1.0}

	for i := 0; i < 4; i++ {
		walls[i] = &face{
			p0: p0.Rotation2D(float64(i)*math.Pi/2.0),
			p1: p1.Rotation2D(float64(i)*math.Pi/2.0),
			depth: 0.1,
		}
	}

	square := newSquare(basic.Zero())
	return &controller{s: square, walls: walls}
}

func (c *controller) SetSquare(x, y float64) {
	square := newSquare(basic.NewPoint(x, y, 0.0))
	c.s = square
}


func (c *controller) Update() {
	c.s.force = &basic.Point{Y: -9.8*c.s.mass}

	//c.collide()
	c.collideV2()

	c.s.accelerate(0.01)
}


func (c *controller) collide() {
	faces := make([]*face, 4)
	faces[0], faces[1], faces[2], faces[3] = c.s.faces()

	for _, wall := range c.walls {
		for _, f := range faces {
			x, d := wall.detectCollision(f.p0)
			if x != nil {
				k := PenetrationEffect*d - VelocityEffect*c.s.velocity.Dot(wall.normal())
				penalty := wall.normal().Mult(k)

				c.s.addFource(penalty)
				c.s.addTorqueToPoint(penalty, f.p0)

			}
		}
	}
}

func (c *controller) collideV2() {
	faces := make([]*face, 4)
	faces[0], faces[1], faces[2], faces[3] = c.s.faces()
	center := c.s.position

	for _, wall := range c.walls {
		for _, f := range faces {
			u0, u1, d := wall.detectCollisionV2(f.p0, f.p1)
			if u0 != nil {
				penaltyForce := wall.normal().Mult(d*PenetrationEffect*(u1.Sub(u0).Length())/2.0 - VelocityEffect*c.s.velocity.Dot(wall.normal()))

				penaltyTorque := u0.Sub(center).Mult(1.0/3.0).Add(u1.Sub(center).Mult(1.0/6.0)).Cross(wall.normal())
				penaltyTorque = penaltyTorque.Mult(d*PenetrationEffect*(u1.Sub(u0).Length()))


				//log.Println(u0)
				c.s.addFource(penaltyForce)
				c.s.addTorque(penaltyTorque)

				//log.Println(c.s.force)
			}
		}
	}
}

func (c *controller) Draw() {
	array := make([]float32, 3*7*2)

	offset := 0

	p0, p1, p2, p3 := c.s.points()

	points := make([][]*basic.Point, 2)
	points[0] = []*basic.Point{p0, p1, p2}
	points[1] = []*basic.Point{p3, p0, p2}

	for i, ps := range points {
		for j, p := range ps {
			array[offset+(i*3+j)*7+0], array[offset+(i*3+j)*7+1], array[offset+(i*3+j)*7+2] = p.Elements()
			array[offset+(i*3+j)*7+3], array[offset+(i*3+j)*7+4], array[offset+(i*3+j)*7+5], array[offset+(i*3+j)*7+6] = c.s.colorElements()
		}
	}

	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, 3*7*2)
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
