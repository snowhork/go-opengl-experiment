package controller

import (
	"../gfx"
	"../basic"
	"github.com/go-gl/gl/v4.1-core/gl"
	"gonum.org/v1/gonum/mat"
	"math"
)

type Controllable interface {
	Update()
	Draw()
}

const (
	N = 128
	g = 0.05
	dt = 1.0/8.0
)

type controller struct {
	ux *velocities
	uy *velocities
	hs *heights
	params *gridParams
	hsum float64

	cnt int
	Controllable
}

func NewController() *controller {
	shaderinit()

	//walls := make([]*face, 4)
	//
	//p0 := &basic.Point{X: -1.0, Y: -1.0}
	//p1 := &basic.Point{X: 1.0, Y: -1.0}
	//
	//for i := 0; i < 4; i++ {
	//	walls[i] = &face{
	//		p0: p0.Rotation2D(float64(i)*math.Pi/2.0),
	//		p1: p1.Rotation2D(float64(i)*math.Pi/2.0),
	//		depth: 0.1,
	//	}
	//}

	params := NewGridParams(-1,-1,1,1,N,N)
	hgridParams := NewGridParams(-1.0+params.dx(),-1.0+params.dy(),1-params.dx(),1-params.dy(),N-1,N-1)

	uValues := make([][]*mat.VecDense, N+1)
	for i := 0; i < N+1; i++ {
		uValues[i] = make([]*mat.VecDense, N+1)
		for j := 0; j < N+1; j++ {
			uValues[i][j] = basic.ZeroVec()
		}
	}
	ux, uy := NewVelocities(N, params)

	hsum := 0.0
	hValues := make([][]float64, N)
	for i := 0; i < N; i++ {
		hValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			x, y := hgridParams.pos(i, j)
			hValues[i][j] = 0.1*math.Pow((1-x)*(1+x)*(1-y)*(1+y), 4)

			hsum += hValues[i][j]
		}
	}

	//log.Println(hgridParams.pos(N, 0))


	hs := &heights{hValues, hgridParams}

	return &controller{ux: ux, uy: uy, hs: hs, params: params, hsum: hsum}
}

func (c *controller) SetWave(inx, iny float64) {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			x, y := c.params.pos(i, j)
			x = inx + x
			y = iny + y

			if x >= 1 || x <= -1 || y >= 1 || y <= -1 {
				continue
			}
			cal := 0.1*math.Pow((1-x)*(1+x)*(1-y)*(1+y), 4)
			cal = math.Max(0, cal)

			c.hs.values[i][j] = math.Max(cal, c.hs.values[i][j])
			c.ux.values[i][j] = 0
			c.uy.values[i][j] = 0
		}
	}

}

func (c *controller) advect() {
	uxValues := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		uxValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			rx, ry := c.ux.gridParams.pos(i, j)

			uij_x := c.ux.at(rx, ry)
			uij_y := c.uy.at(rx, ry)

			ux := c.ux.at(rx-dt*uij_x, ry-dt*uij_y)

			if i == 0 || i == N {
				ux = 0
			}
			uxValues[i][j] = ux
		}
	}

	uyValues := make([][]float64, N)
	for i := 0; i < N; i++ {
		uyValues[i] = make([]float64, N+1)
		for j := 0; j <= N; j++ {
			rx, ry := c.uy.gridParams.pos(i, j)

			uij_x := c.ux.at(rx, ry)
			uij_y := c.uy.at(rx, ry)

			uy := c.uy.at(rx-dt*uij_x, ry-dt*uij_y)
			if j == 0 || j == N {
				uy = 0
			}
			uyValues[i][j] = uy
		}
	}

	hValues := make([][]float64, N)
	newhsum := 0.0
	for i := 0; i < N; i++ {
		hValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			rx, ry := c.hs.gridParams.pos(i, j)
			uij_x := c.ux.at(rx, ry)
			uij_y := c.uy.at(rx, ry)
			h := c.hs.at(rx-dt*uij_x, ry-dt*uij_y)

			//log.Println(h, c.hs.ix(i,j))

			//if j == 0 {
			//	log.Println(h)
			//}
			hValues[i][j] = h
			newhsum += hValues[i][j]
		}
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			hValues[i][j] *= c.hsum/newhsum
		}
	}

	c.hs.values = hValues
	c.ux.values = uxValues
	c.uy.values = uyValues

	c.cnt += 1
}

func (c *controller) Update() {
	dx := c.params.dx()
	dy := c.params.dy()

	//ix, iy, _ := c.params.indexAt(0.4, 0.4)

	//log.Println("h", c.hs.ix(N+1,N+1))
	//if c.cnt >= 200 {
	//	return
	//} else {
	//	hsum := 0.0
	//	hs := make([][]float64, N+1)
	//	log.Println(c.cnt, "**********************")
	//
	//	for i := 0; i < N; i++ {
	//		hs[i] = make([]float64, N)
	//		for j := 0; j < N; j++ {
	//			hs[i][j] = float64(c.ux.values[i][j])
	//			hsum += c.hs.values[i][j]
	//		}
	//		//log.Println(i, hs[i])
	//	}
	//
	//	log.Println(hsum)
	//
	//}

	c.advect()
	uxValues := make([][]float64, N+1)
	for i := 0; i <= N; i++ {
		uxValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			rx, ry := c.ux.gridParams.pos(i, j)
			ux := c.ux.at(rx, ry) - g*(c.hs.at(rx+dx, ry)-c.hs.at(rx-dx, ry))/(2*dx)*dt

			if i == 0 || i == N {
				ux = 0
			}
			uxValues[i][j] = ux
		}
	}

	uyValues := make([][]float64, N)
	for i := 0; i < N; i++ {
		uyValues[i] = make([]float64, N+1)
		for j := 0; j <= N; j++ {
			rx, ry := c.uy.gridParams.pos(i, j)
			uy := c.uy.at(rx, ry) - g*(c.hs.at(rx, ry+dy)-c.hs.at(rx, ry-dy))/(2*dy)*dt

			if j == 0 || j == N {
				uy = 0
			}
			uyValues[i][j] = uy
		}
	}

	hValues := make([][]float64, N)
	newhsum := 0.0
	for i := 0; i < N; i++ {
		hValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			rx, ry := c.hs.gridParams.pos(i, j)
			//uij_x := c.ux.at(rx, ry)
			//uij_y := c.uy.at(rx, ry)
			h := c.hs.at(rx, ry)


			//log.Println(h, c.hs.ix(i,j))

			//if j == 0 {
			//	log.Println(h)
			//}
			hValues[i][j] = h - h*dt*(
				(c.ux.at(rx+dx, ry)-c.ux.at(rx-dy, ry))/(2*dx) +
				(c.uy.at(rx, ry+dy)-c.uy.at(rx, ry-dy))/(2*dy))

			newhsum += hValues[i][j]
		}
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			hValues[i][j] *= c.hsum/newhsum
		}
	}

	c.hs.values = hValues
	c.ux.values = uxValues
	c.uy.values = uyValues

	c.cnt += 1
}


func (c *controller) Draw() {
	array := make([]float32, (N+2)*(N+2)*3*7*2)

	for i := 0; i <= N+1; i++ {
		iOffset := 3*7*2*(N+2)*i
		for j := 0; j <= N+1; j++ {
			x0, y0 := c.hs.gridParams.pos(i-1, j-1)
			h0 := c.hs.ix(i-1, j-1)
			x1, y1 := c.hs.gridParams.pos(i, j-1)
			h1 := c.hs.ix(i, j-1)
			x2, y2 := c.hs.gridParams.pos(i, j)
			h2 := c.hs.ix(i, j)
			x3, y3 := c.hs.gridParams.pos(i-1, j)
			h3 := c.hs.ix(i-1, j)
			//log.Println(x2)

			xs := []float32{float32(x0), float32(x1), float32(x3), float32(x1), float32(x2), float32(x3)}
			ys := []float32{float32(y0), float32(y1), float32(y3), float32(y1), float32(y2), float32(y3)}
			hs := []float32{float32(h0), float32(h1), float32(h3), float32(h1), float32(h2), float32(h3)}

			//xs := []float32{float32(x0), float32(x3), float32(x3), float32(x1), float32(x3), float32(x2)}
			//ys := []float32{float32(y0), float32(y1), float32(y3), float32(y1), float32(y3), float32(y2)}
			//hs := []float32{float32(h0), float32(h1), float32(h3), float32(h1), float32(h3), float32(h2)}

			for k := 0; k < 6; k++ {
				//log.Println(iOffset + 7*3*2*j + 7*k + 0)
				array[iOffset + 7*(k+3*2*(j)) + 0], array[iOffset + 7*(k+3*2*(j)) + 1], array[iOffset + 7*(k+3*2*(j)) + 2] = xs[k], hs[k], ys[k]
				array[iOffset + 7*(k+3*2*(j)) + 3], array[iOffset + 7*(k+3*2*(j)) + 4], array[iOffset + 7*(k+3*2*(j)) + 5], array[iOffset + 7*(k+3*2*(j)) + 6] = hs[k]*10, hs[k]*6 + 0.4,1.0,1.0
			}
		}
	}
	VAO := makeVao(array)
	gl.BindVertexArray(VAO)

	gl.DrawArrays(gl.TRIANGLES, 0, (N+2)*(N+2)*3*7*2)
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
