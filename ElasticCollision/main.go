package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"./basic"
	"./ball"
	"math/rand"
)

const (
	width  = 500
	height = 500
)

type Drawable interface {
	Draw()
}


const BallNum = 1
func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()

	initOpenGL()

	rand.Seed(123456)
	balls := make([]*ball.Ball, BallNum)

	for i := 0; i < BallNum; i++ {
		balls[i] = ball.NewBall(&basic.Point{X: -0.9+float32(i)*0.2, Y: rand.Float32()})
	}

	for !window.ShouldClose() {
		for i := 0; i < BallNum; i++ {
			balls[i].Update(0.01)
		}
		draw(balls, window)
	}
}

func draw(balls []*ball.Ball, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for i := 0; i < BallNum; i++ {
		balls[i].Draw()
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func initGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "Title", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
}
