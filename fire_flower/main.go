package main

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"./basic"
	"./fire_flower"
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

	flowers := make([]*fire_flower.FireFlower, 0, 20)

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

		log.Println(x, y)
		flower := fire_flower.NewFireFlower(basic.NewPoint(float32(x), float32(y), 0),
			0.002, 3, 50)
		flowers = append(flowers, flower)

	}

	window.SetMouseButtonCallback(mouseCallback)
	for !window.ShouldClose() {
		for _, flower := range flowers  {
			flower.Update()
		}
		draw(flowers, window)
	}
}

func draw(drawables []*fire_flower.FireFlower, window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, drawable := range drawables {
		drawable.Draw()
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
