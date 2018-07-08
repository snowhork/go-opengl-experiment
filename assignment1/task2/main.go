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
	var particle *ball.Ball = nil
	var pos *basic.Point = nil

	mouseCallback := func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if action != glfw.Press {
			return
		}
		x, y := window.GetCursorPos()
		w, h := window.GetSize()

		x -= float64(w)/2.0
		y -= float64(h)/2.0
		x /= float64(w)/2.0
		y /= -float64(h)/2.0

		pos = &basic.Point{X: float32(x), Y: float32(y)}

		particle = ball.NewBall(pos)
	}

	window.SetMouseButtonCallback(mouseCallback)

	for !window.ShouldClose() {
		if particle != nil {
			particle.Update(1.0/60.0)
		}
		draw(particle, window)
	}
}

func draw(b *ball.Ball, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	if b != nil {
		b.Draw()
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
