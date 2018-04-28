package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"./basic"
	"./line_segment"
)

const (
	width  = 500
	height = 500
)

type Drawable interface {
	Draw()
}

func main() {
	runtime.LockOSThread()

	window := initGlfw()
	defer glfw.Terminate()


	initOpenGL()

	points := make([]*basic.Point, 4, 4)
	count := 0
	bez := line_segment.NewBezier(basic.Zero(), basic.Zero(), basic.Zero(), basic.Zero())

	mouseCallback := func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
		if count == 4 {
			return
		}
		if action != glfw.Press {
			return
		}
		x, y := window.GetCursorPos()
		w, h := window.GetSize()

		x -= float64(w)/2.0
		y -= float64(h)/2.0
		x /= float64(w)/2.0
		y /= -float64(h)/2.0

		points[count] = basic.NewPoint(float32(x), float32(y), 0)
		count += 1

		if count == 4 {
			bez = line_segment.NewBezier(points[0], points[1], points[2], points[3])
			//bez = line_segment.NewBezier(
			//	&basic.Point{X: -0.7},
			//	&basic.Point{X: -0.25},
			//	&basic.Point{X: 0.25},
			//	&basic.Point{X: 0.8},
			//)
			count = 0
		}
	}

	window.SetMouseButtonCallback(mouseCallback)
	for !window.ShouldClose() {
		bez.Update()
		draw(bez, window)
	}
}

func draw(drawable Drawable, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	drawable.Draw()

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
