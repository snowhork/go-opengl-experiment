package galaxy

import (
	"../basic"
)

type node struct {
	parent *node
	children []*node
	balance *basic.Point
	stars []*star
	mass float64
	xSum float64
	ySum float64
	pxSum float64
	pySum float64

	xMax float64
	xMin float64
	yMax float64
	yMin float64

}

func (g *Galaxy) Tree() *node {
	gx := 0.0
	gy := 0.0
	M := 0.0

	xMax, xMin, yMax, yMin := 0.0, 0.0, 0.0, 0.0

	for _, star := range g.stars {
		gx += star.Current.X*star.mass
		gy += star.Current.Y*star.mass
		M += star.mass

		if star.Current.X > xMax {
			xMax = star.Current.X
		}
		if star.Current.X < xMin {
			xMin = star.Current.X
		}
		if star.Current.Y > yMax {
			yMax = star.Current.Y
		}
		if star.Current.Y < yMin {
			yMin = star.Current.Y
		}


	}

	gx /= M
	gy /= M

	root := &node{
		parent:   nil,
		children: make([]*node, 0, 4),
		balance:  &basic.Point{X: gx, Y: gy},
		stars:    g.stars,
		mass:     M,
		xMax:     2.0,
		xMin:     -2.0,
		yMax:     2.0,
		yMin:     -2.0,
	}

	nodes := make([]*node, 0, len(g.stars))
	root.UpdateByTree(&nodes)

	return root
}

func newNode(parent *node, xMax, xMin, yMax, yMin float64) *node {
	count := 12

	n := &node{
		parent:     parent,
		children:   make([]*node, 0, 4),
		balance:    nil,
		stars:      make([]*star, 0, count),
		mass:       0,
		xMax:       xMax,
		xMin: 		xMin,
		yMax:       yMax,
		yMin:       yMin,
	}
	return n
}

func (n *node) appendStar(star *star) {
	n.stars = append(n.stars, star)
	n.mass += star.mass
	n.xSum += star.Current.X
	n.ySum += star.Current.Y
	n.pxSum += star.Current.X*star.mass
	n.pySum += star.Current.Y*star.mass
}

func (n *node) finalize() {
	if len(n.stars) == 0 {
		return
	}
	n.balance = &basic.Point{X: n.pxSum/n.mass, Y: n.pySum/n.mass}
	n.parent.children = append(n.parent.children, n)
}

func (n *node) UpdateByTree(nodes *[]*node) {
	if len(n.stars) == 0 {
		return
	}
	if len(n.stars) == 1 {
		*nodes = append(*nodes, n)
		return
	}

	gx, gy := (n.xMax+n.xMin)/2.0, (n.yMax+n.yMin)/2.0

	lowerLeft := newNode(n, gx, n.xMin, gy, n.yMin)
	upperLeft := newNode(n, gx, n.xMin, n.yMax, gy)
	lowerRight := newNode(n, n.xMax, gx, gy, n.yMin)
	upperRight := newNode(n, n.xMax, gx, n.yMax, gy)

	for _, star := range n.stars {
		if star.Current.X >= gx {
			if star.Current.Y >= gy {
				upperRight.appendStar(star)
			} else {
				lowerRight.appendStar(star)
			}
		} else {
			if star.Current.Y >= gy {
				upperLeft.appendStar(star)
			} else {
				lowerLeft.appendStar(star)
			}
		}
	}

	lowerLeft.finalize()
	upperLeft.finalize()
	lowerRight.finalize()
	upperRight.finalize()

	//if len(yMin.stars) > 1 {
	//	log.Println(yMin.stars[0].Current, yMin.stars[1].Current)
	//}
	//log.Println(len(yMin.stars), len(yMax.stars), len(lowerRight.stars), len(xMax.stars))

	lowerLeft.UpdateByTree(nodes)
	upperLeft.UpdateByTree(nodes)
	lowerRight.UpdateByTree(nodes)
	upperRight.UpdateByTree(nodes)
}
