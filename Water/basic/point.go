package basic

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func NewVec(X,Y,Z float64) *mat.VecDense {
	return mat.NewVecDense(3, []float64{X, Y, Z})
}

func ZeroVec() *mat.VecDense {
	return mat.NewVecDense(3, []float64{0, 0, 0})
}

func ZeroMat() *mat.Dense {
	return mat.NewDense(2 , 2, []float64{0, 0, 0, 0})
}

func Rotate2D(v *mat.VecDense, theta float64) *mat.VecDense {
	x := v.At(0,0)
	y := v.At(1,0)

	x_ := math.Cos(theta)*x - math.Sin(theta)*y
	y_ := math.Sin(theta)*x + math.Cos(theta)*y
	return mat.NewVecDense(2, []float64{x_, y_})
}

func Rotate2DAt(v *mat.VecDense, theta float64, at *mat.VecDense) *mat.VecDense {
	x := v.At(0,0) - at.At(0,0)
	y := v.At(1,0) - at.At(1,0)

	newX := math.Cos(theta)*x - math.Sin(theta)*y + at.At(0,0)
	newY := math.Sin(theta)*x + math.Cos(theta)*y + at.At(1,0)

	return mat.NewVecDense(2, []float64{newX, newY})
}
