package line_segment

import (
	"math"
)

type Sturm struct {
	f0, f1, f2, f3, f4, f5 Poly
}

func NewSturm(a0, a1, a2, a3, a4, a5 float64) *Sturm {
	f0 := &Poly5{a0, a1, a2, a3, a4, a5}
	f1 := &Poly4{a0*5, a1*4, a2*3, a3*2, a4}

	q0 := f0.a0/f1.a0
	q1 := (f0.a1-q0*f1.a1)/f1.a0
	f2 := &Poly3{
		a0: -(f0.a2 - q1*f1.a1 - q0*f1.a2),
		a1: -(f0.a3 - q1*f1.a2 - q0*f1.a3),
		a2: -(f0.a4 - q1*f1.a3 - q0*f1.a4),
		a3: -(f0.a5 - q1*f1.a4),
	}

	q0 = f1.a0/f2.a0
	q1 = (f1.a1-q0*f2.a1)/f2.a0

	f3 := &Poly2{
		a0: -(f1.a2 - q1*f2.a1 - q0*f2.a2),
		a1: -(f1.a3 - q1*f2.a2 - q0*f2.a3),
		a2: -(f1.a4 - q1*f2.a3),
	}

	q0 = f2.a0/f3.a0
	q1 = (f2.a1-q0*f3.a1)/f3.a0

	f4 := &Poly1{
		a0: -(f2.a2 - q1*f3.a1 - q0*f3.a2),
		a1: -(f2.a3 - q1*f3.a2),
	}

	q0 = f3.a0/f4.a0
	q1 = (f3.a1-q0*f4.a1)/f4.a0

	f5 := &Poly0{
		a0: -(f3.a2 - q1*f4.a1),
	}

	return &Sturm{
		f0, f1, f2, f3, f4, f5,
	}
}

func (s *Sturm) RootsNumber(t0, t1 float64) int {
	fs := []Poly{s.f1, s.f2, s.f3, s.f4, s.f5}

	N := func(x float64) int {
		count := 0
		a := s.f0.calc(x)
		for _, f := range fs {
			if a*f.calc(x) < 0 {
				count += 1
				a = f.calc(x)
			}
		}
		return count
	}

	return N(t0) - N(t1)
}

func (s *Sturm) Root(t0, t1 float64, res []float64) []float64 {
	BiEPS := 1e-4
	NewtonEPS := 1e-6

	num := s.RootsNumber(t0, t1)
	if num == 1 {
		if t1-t0 > BiEPS {
			res = s.Root(t0, (t0+t1)/2.0, res)
			res = s.Root((t0+t1)/2.0, t1, res)
		} else {
			t := (t0 + t1) / 2.0
			for i := 0; i < 50; i++ {
				tNext := s.newton(t)
				if math.Abs(t-tNext) <= NewtonEPS {
					res = append(res, tNext)
					return res
				}
				t = tNext
			}
		}

	}
	if num >= 2 {
		res = s.Root(t0, (t0+t1)/2.0, res)
		res = s.Root((t0+t1)/2.0, t1, res)
	}
	return res
}

func (s *Sturm) newton(x float64) float64 {
	return x - s.f0.calc(x)/s.f1.calc(x)
}