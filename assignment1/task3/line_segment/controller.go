package line_segment

import (
	"../basic"
)

type Controller struct {
	b *bead
	l *line
}

func NewController(p0, p1 *basic.Point) *Controller {
	l := newLine(p0, p1)
	b := newBead(p0, p0)

	return &Controller{l: l, b: b}
}

func (con *Controller) Draw() {
	con.l.Draw()
	con.b.Draw()
}

func (con *Controller) Update(dt float32) {
	g := &basic.Point{Y: -9.8}
	next := con.b.current.Add(con.b.current).Sub(con.b.prev).Add(g.Mult(dt*dt))
	con.b.prev = con.b.current

	t := next.Sub(con.l.p0).Product(con.l.p1.Sub(con.l.p0))/con.l.p1.Sub(con.l.p0).Length2()

	if t < 0 {
		t = 0
	}

	if t > 1 {
		t = 1
	}

	con.b.current = con.l.p1.Mult(t).Add(con.l.p0.Mult(1-t))
}