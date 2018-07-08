package controller

type velocities struct {
	values [][]float64
	xN int
	yN int
	*gridParams
}

func NewVelocities(N int, params *gridParams) (*velocities, *velocities) {
	dx := params.dx()
	dy := params.dy()
	uxParams := NewGridParams(params.xmin, params.ymin+dy/2, params.xmax, params.ymax-dy/2, N, N-1)
	uyParams := NewGridParams(params.xmin+dx/2, params.ymin, params.xmax-dx/2, params.ymax, N-1, N)

	uxValues := make([][]float64, N+1)
	for i := 0; i < N+1; i++ {
		uxValues[i] = make([]float64, N)
		for j := 0; j < N; j++ {
			uxValues[i][j] = 0
		}
	}
	uyValues := make([][]float64, N)
	for i := 0; i < N; i++ {
		uyValues[i] = make([]float64, N+1)
		for j := 0; j < N+1; j++ {
			uyValues[i][j] = 0
		}
	}

	return &velocities{uxValues, N+1, N, uxParams},
	&velocities{uyValues, N, N+1,uyParams}
}

func (us *velocities) at(x,y float64) float64 {
	ix, iy, areas := us.gridParams.indexAt(x, y)

	return us.ix(ix, iy)*areas[0] +
		us.ix(ix+1, iy)*areas[1] +
		us.ix(ix, iy+1)*areas[2] +
		us.ix(ix+1, iy+1)*areas[3]
}

func (us *velocities) ix(i,j int) float64 {
	if (i <= -1 || i >= us.xN) && (j <= -1 || j >= us.yN) {
		return 0
	}
	if i <= -1 {
		return us.values[0][j]
	}
	if i >= us.xN {
		return us.values[us.xN-1][j]
	}
	if j <= -1 {
		return us.values[i][0]
	}
	if j >= us.yN {
		return us.values[i][us.yN-1]
	}
	return us.values[i][j]
}

