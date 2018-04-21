package fire_flower

import (
	"github.com/lucasb-eyer/go-colorful"
	"../basic"
)

type line struct {
	count int32
	points []float32
	vertexCount int
}



func newLine(
	position *basic.Point,
	velocity *basic.Point,
	color *colorful.Color,
	alpha float32,
	vertexCount int) *line {

	points := make([]float32, 7*vertexCount, 7*vertexCount)

	for i := 0; i <= vertexCount-2; i++ {
		points[i*7+0] = position.X
		points[i*7+1] = position.Y
		points[i*7+2] = position.Z
		points[i*7+3] = float32(color.R)
		points[i*7+4] = float32(color.G)
		points[i*7+5] = float32(color.B)
		points[i*7+6] = alpha
	}

	points[(vertexCount-1)*7+0] = position.X+velocity.X
	points[(vertexCount-1)*7+1] = position.Y+velocity.Y
	points[(vertexCount-1)*7+2] = position.Z+velocity.Z
	points[(vertexCount-1)*7+3] = float32(color.R)
	points[(vertexCount-1)*7+4] = float32(color.G)
	points[(vertexCount-1)*7+5] = float32(color.B)
	points[(vertexCount-1)*7+6] = alpha

	return &line{
		count: 0,
		points: points,
		vertexCount: vertexCount,
	}
}

//func (line *Line) Draw(makeVao func([]float32) uint32) {
//	gl.Uniform3f(line.program.GetUniformLocation("objectColor"), float32(line.color.R), float32(line.color.G), float32(line.color.B))
//
//	drawable := makeVao(line.points)
//	gl.BindVertexArray(drawable)
//	gl.DrawArrays(gl.LINE_STRIP, 0, int32(vertexCount-1))//
//}

func (line *line) update() {
	line.count += 1

	speed := float32(1.0/900.0)
	g := -float32(9.8)

	oldX     := float32(line.points[(line.vertexCount-2)*7+0])
	currentX := float32(line.points[(line.vertexCount-1)*7+0])
	oldY     := float32(line.points[(line.vertexCount-2)*7+1])
	currentY := float32(line.points[(line.vertexCount-1)*7+1])
	currentZ := float32(line.points[(line.vertexCount-1)*7+2])
	currentR := float32(line.points[(line.vertexCount-1)*7+3])
	currentG := float32(line.points[(line.vertexCount-1)*7+4])
	currentB := float32(line.points[(line.vertexCount-1)*7+5])
	currentA := float32(line.points[(line.vertexCount-1)*7+6])

	newX := currentX + (currentX - oldX)
	newY := currentY + (currentY - oldY) + speed*speed*g

	nextPoints := []float32{
		newX, newY, currentZ, currentR, currentG, currentB, currentA,
	}

	newPoints := make([]float32, 7*line.vertexCount, 7*line.vertexCount)
	copy(newPoints, append(line.points[7:], nextPoints...))

	for i := 0; i < line.vertexCount ; i++  {
		newPoints[i*7+4] = float32(i)/float32(line.vertexCount-1)
	}
	line.points = newPoints
}

