package galaxy

import (
	"../basic"
)

type node struct {
	parent *node
	children []*node
	balance *basic.Point
	stars []*star
	mass float32
	xSum float32
	ySum float32
	pxSum float32
	pySum float32

}

func (g *Galaxy) Tree() []*node {
	n := &node{
		parent: nil,
		children: make([]*node, 0, 4),
		balance: nil,
		stars: g.stars,
		mass: 0,
	}

	nodes := make([]*node, 0, len(g.stars))
	n.UpdateByTree(&nodes)

	return nodes
}

func newNode(parent *node) *node {
	count := 12

	n := &node{
		parent: parent,
		children: make([]*node, 0, 4),
		balance: nil,
		stars: make([]*star, 0, count),
		mass: 0,
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
	n.pxSum /= n.mass
	n.pySum /= n.mass
	n.balance = &basic.Point{X: n.pxSum, Y: n.pySum}

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

	gx, gy := n.calcBalance()

	lowerLeft := newNode(n)
	upperLeft := newNode(n)
	lowerRight := newNode(n)
	upperRight := newNode(n)

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

	//if len(lowerLeft.stars) > 1 {
	//	log.Println(lowerLeft.stars[0].Current, lowerLeft.stars[1].Current)
	//}
	//log.Println(len(lowerLeft.stars), len(upperLeft.stars), len(lowerRight.stars), len(upperRight.stars))

	lowerLeft.UpdateByTree(nodes)
	upperLeft.UpdateByTree(nodes)
	lowerRight.UpdateByTree(nodes)
	upperRight.UpdateByTree(nodes)
}

func (n *node) calcBalance() (float32, float32) {
	xMean, yMean := float32(0.0), float32(0.0)

	for _, star := range n.stars {
		xMean += star.Current.X
		yMean += star.Current.Y
	}
	xMean /= float32(len(n.stars))
	yMean /= float32(len(n.stars))
	return xMean, yMean
}

