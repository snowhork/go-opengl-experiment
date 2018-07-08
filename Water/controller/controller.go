package controller

import (
	"../gfx"
	"../basic"
	"github.com/go-gl/gl/v4.1-core/gl"
	"gonum.org/v1/gonum/mat"
	"math"
	"log"
)

type Controllable interface {
	Update()
	Draw()
}

const (
	D = 1

	n = 10 // N//3
	N = 1000

	h = 0.1
	hN = 11 // 1.0/hN

	nu = 0.01
	g = 0.1
	dt = 0.01
	rh0 = 1000

	stiffness = 0.1
)

type controller struct {
	particles []*particle
	neighbors [][][][]int

	rho0 float64

	Controllable
}

func NewController() *controller {
	shaderinit()

	particles := make([]*particle, N)

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				pos := mat.NewVecDense(3, []float64{float64(i) / n, float64(j) / n, float64(k) / n})
				mass := 4*math.Pow(D, 3)*math.Pi/(3*N)*rh0
				particles[i*n*n+j*n+k] = NewParticle(mass, pos)
			}
		}
	}

	c := &controller{particles: particles, rho0: rh0}
	c.setNeighbors()

	return c
}

func (c *controller) setNeighbors() {
	neighbors := make([][][][]int, hN)

	for i := 0; i < hN; i++ {
		neighbors[i] = make([][][]int, hN)
		for j := 0; j < hN; j++ {
			neighbors[i][j] = make([][]int, hN)
			for k := 0; k < hN; k++ {
				neighbors[i][j][k] = make([]int, hN)
			}

		}
	}


	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			for k := 0; k < n; k++ {
				p := c.particles[i*n*n+j*n+k]
				ix, iy, iz := p.getIndex()

				neighbors[ix][iy][iz] = append(neighbors[ix][iy][iz], i*n*n+j*n+k)
			}
		}
	}

	c.neighbors = neighbors
}

func (p *particle) getIndex() (ix,iy,iz int) {
	x := p.pos.At(0,0)
	y := p.pos.At(1,0)
	z := p.pos.At(2,0)
	x = math.Min(math.Max(x, 0.0), 1.0)
	y = math.Min(math.Max(y, 0.0), 1.0)
	z = math.Min(math.Max(z, 0.0), 1.0)
	ix = int(x/h)
	iy = int(y/h)
	iz = int(z/h)
	return
}

func (c *controller) getNeighbors(p *particle) []int {
	ix, iy, iz := p.getIndex()
	return c.neighbors[ix][iy][iz]
}

func DensityKernel(x, y *mat.VecDense) float64 {
	r2 := (x.At(0,0)-y.At(0,0))*(x.At(0,0)-y.At(0,0)) +
		  	(x.At(1,0)-y.At(1,0))*(x.At(1,0)-y.At(1,0)) +
			(x.At(2,0)-y.At(2,0))*(x.At(2,0)-y.At(2,0))

	if r2 > h*h {
		return 0
	}

	return 315.0/(64*math.Pi*math.Pow(h,9))*math.Pow(h*h-r2, 3)
}

func (c *controller) Density(i int) float64 {
	res := 0.0
	p := c.particles

	for _, j := range c.getNeighbors(p[i]) {

		res += p[i].mass*DensityKernel(p[i].pos, p[j].pos)
	}

	return res
}

func PressureGradKernel(x, y *mat.VecDense) *mat.VecDense {
	res := basic.ZeroVec()

	r2 := (x.At(0,0)-y.At(0,0))*(x.At(0,0)-y.At(0,0)) +
		(x.At(1,0)-y.At(1,0))*(x.At(1,0)-y.At(1,0)) +
		(x.At(2,0)-y.At(2,0))*(x.At(2,0)-y.At(2,0))
	if r2 > h*h {
		return res
	}

	r := math.Sqrt(r2)
	res.SubVec(x, y)
	res.ScaleVec(-45.0/(math.Pi*math.Pow(h, 6))*math.Pow(h-r, 2)/r, res)

	return res
}

func (c *controller) Pressure(i int) float64 {
	res := c.Density(i)-c.rho0
	if res <= 0 {
		return 0
	}
	return res
}


func (c *controller) PressureGrad(i int) *mat.VecDense {
	res := basic.ZeroVec()
	p := c.particles
	for _, j := range c.getNeighbors(p[i]) {
		coef := stiffness*(c.Pressure(i) + c.Pressure(j))/2.0*(p[i].mass/c.Density(i))
		res.AddScaledVec(PressureGradKernel(p[0].pos, p[1].pos), coef, res)
	}

	return res
}

func VelocityLaplaceKernel(x, y *mat.VecDense) float64 {
	res := 0.0

	r2 := (x.At(0,0)-y.At(0,0))*(x.At(0,0)-y.At(0,0)) +
		(x.At(1,0)-y.At(1,0))*(x.At(1,0)-y.At(1,0)) +
		(x.At(2,0)-y.At(2,0))*(x.At(2,0)-y.At(2,0))

	if r2 > h*h {
		return res
	}

	r := math.Sqrt(r2)

	return 45.0/(math.Pi*math.Pow(h, 6))*(h-r)
}

func (c *controller) VelocityLaplace(i int) *mat.VecDense {
	res := basic.ZeroVec()
	p := c.particles
	for _, j := range c.getNeighbors(p[i]) {
		res2 := basic.ZeroVec()
		res2.SubVec(p[j].vel, p[i].vel)
		res2.ScaleVec(VelocityLaplaceKernel(p[i].pos, p[j].pos), res2)

		res.AddVec(res2, res)
	}

	return res
}

func (c *controller) Force(i int) *mat.VecDense {
	return basic.NewVec(0,-c.particles[i].mass*g, 0)
}

func (c *controller) Update() {
	log.Println("frame")
	newVels := make([]*mat.VecDense, N)

	for i := 0; i < N; i++ {
		newVels[i] = basic.ZeroVec()
		v := newVels[i]
		p := c.particles[i]


		v.AddScaledVec(c.VelocityLaplace(i), nu, v)
		v.AddScaledVec(c.PressureGrad(i), -1.0/c.Density(i), v)
		v.AddVec(c.Force(i),v)
		v.ScaleVec(dt, v)
		v.AddVec(p.vel, v)
	}

	for i := 0; i < N; i++ {
		p := c.particles[i]
		p.pos.AddScaledVec(p.pos, dt, newVels[i])
		p.vel = newVels[i]
	}

	c.setNeighbors()

}


func (c *controller) Draw() {
	array := make([]float32, N*7)

	for i := 0; i < N; i++ {
		p := c.particles[i]

		array[i*7 + 0], array[i*7 + 1], array[i*7 + 2] = p.glpos()
		array[i*7 + 3], array[i*7 + 4], array[i*7 + 5], array[i*7 + 6] = 0.0, 0.5, 1.0, 1.0
	}
	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.POINTS, 0, N)
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

	program.SetProjectionMat()
	program.SetCamera()
	program.SetModel()
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
