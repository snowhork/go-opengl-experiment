package controller

type heights struct {
	values [][]float64
	*gridParams
}

func (hs *heights) at(x,y float64) float64 {
	ix, iy, areas := hs.indexAt(x, y)

	return hs.ix(ix, iy)*areas[0] +
		hs.ix(ix+1, iy)*areas[1] +
		hs.ix(ix, iy+1)*areas[2] +
		hs.ix(ix+1, iy+1)*areas[3]
}

func (hs *heights) ix(i,j int) float64 {
	if i <= -1 {
		i = 0
	}
	if j <= -1 {
		j = 0
	}
	if i >= hs.xN {
		i = hs.xN
	}
	if j >= hs.yN {
		j = hs.yN
	}
	return hs.values[i][j]
}

