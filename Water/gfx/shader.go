package gfx

import (
  "io/ioutil"
  "strings"
  "fmt"

  "github.com/go-gl/gl/v4.1-core/gl"
  "github.com/go-gl/mathgl/mgl32"
)

type Shader struct {
  handle uint32
}

type Program struct {
  handle uint32
  shaders []*Shader
}

func (shader *Shader) Delete() {
  gl.DeleteShader(shader.handle)
}

func (prog *Program) Delete() {
  for _, shader := range prog.shaders {
    shader.Delete()
  }
  gl.DeleteProgram(prog.handle)
}

func (prog *Program) Attach(shaders ...*Shader) {
  for _, shader := range shaders {
    gl.AttachShader(prog.handle, shader.handle)
    prog.shaders = append(prog.shaders, shader)
  }
}

func (prog *Program) Use() {
  gl.UseProgram(prog.handle)
}

const (
  windowWidth  = 500
  windowHeight = 500
  )

func (prog *Program) SetProjectionMat() {
  projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
  projectionUniform := gl.GetUniformLocation(prog.handle, gl.Str("projection\x00"))
  gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])
}

func (prog *Program) SetCamera() {
  camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
  cameraUniform := gl.GetUniformLocation(prog.handle, gl.Str("camera\x00"))
  gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
}

func (prog *Program) SetModel() {
  model := mgl32.Ident4()
  modelUniform := gl.GetUniformLocation(prog.handle, gl.Str("model\x00"))
  gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])
}

func (prog *Program) Link() error {
  gl.LinkProgram(prog.handle)
  return getGlError(prog.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
    "PROGRAM::LINKING_FAILURE")
}

func (prog *Program) GetUniformLocation(name string) int32 {
  return gl.GetUniformLocation(prog.handle, gl.Str(name + "\x00"))
}

func NewProgram(shaders ...*Shader) (*Program, error) {
  prog := &Program{handle:gl.CreateProgram()}
  prog.Attach(shaders...)

  if err := prog.Link(); err != nil {
    return nil, err
  }

  return prog, nil
}

func NewShader(src string, sType uint32) (*Shader, error) {

  handle := gl.CreateShader(sType)
  glSrc := gl.Str(src + "\x00")
  gl.ShaderSource(handle, 1, &glSrc, nil)
  gl.CompileShader(handle)
  err := getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
    "SHADER::COMPILE_FAILURE::")
  if err != nil {
    return nil, err
  }
  return &Shader{handle:handle}, nil
}

func NewShaderFromFile(file string, sType uint32) (*Shader, error) {
  src, err := ioutil.ReadFile(file)
  if err != nil {
    return nil, err
  }

  handle := gl.CreateShader(sType)
  glSrc, free := gl.Strs(string(src) + "\x00")
  gl.ShaderSource(handle, 1, glSrc, nil)
  free()
  gl.CompileShader(handle)
  err = getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
                    "SHADER::COMPILE_FAILURE::" + file)
  if err != nil {
    return nil, err
  }
  return &Shader{handle:handle}, nil
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
  getObjInfoLogFn getObjInfoLog, failMsg string) error {

  var success int32
  getObjIvFn(glHandle, checkTrueParam, &success)

  if success == gl.FALSE {
    var logLength int32
    getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

    log := gl.Str(strings.Repeat("\x00", int(logLength)))
    getObjInfoLogFn(glHandle, logLength, nil, log)

    return fmt.Errorf("%s: %s", failMsg, gl.GoStr(log))
  }

  return nil
}
