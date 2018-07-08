package controller

import (
	"math"
)

type gridParams struct {
	xmin,ymin,xmax,ymax float64
	xN, yN int
}

func NewGridParams(xmin,ymin,xmax,ymax float64, xN, yN int) *gridParams {
	return &gridParams{xmin: xmin,ymin: ymin,xmax: xmax,ymax: ymax, xN: xN, yN: yN}
}

func (params *gridParams) dx() float64 {
	return (params.xmax - params.xmin)/float64(params.xN)
}

func (params *gridParams) dy() float64 {
	return (params.ymax - params.ymin)/float64(params.yN)
}

func (params *gridParams) indexAt(x,y float64) (int, int, []float64){
	ix := int(math.Floor((x - params.xmin)/ params.dx()))
	xin := x  - params.xmin - float64(ix)*params.dx()

	iy := int(math.Floor((y - params.ymin)/ params.dy()))
	yin := y  - params.ymin - float64(iy)*params.dy()

	allArea := params.dx()*params.dy()

	return ix, iy,
	[]float64{
		(params.dx()-xin)*(params.dy()-yin)/allArea,
		xin*(params.dy()-yin)/allArea,
		(params.dx()-xin)*yin/allArea,
		xin*yin/allArea,
	}
}

func (params *gridParams) pos(ix, iy int) (float64, float64) {
	return params.xmin + float64(ix)*params.dx(), params.ymin + float64(iy)*params.dy()
}