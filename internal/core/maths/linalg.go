package maths

import "math"

// Represents the matrix
//
// | a  b |
//
// | c  d |
type Mat2x2 struct {
	a, b, c, d float64
}

func (A *Mat2x2) TimesMat(B Mat2x2) Mat2x2 {
	return Mat2x2{
		a: A.a*B.a + A.b*B.c,
		b: A.a*B.b + A.b*B.d,
		c: A.c*B.a + A.d*B.c,
		d: A.c*B.b + A.d*B.d,
	}
}

func (A *Mat2x2) Scale(k float64) Mat2x2 {
	return Mat2x2{
		a: A.a * k,
		b: A.b * k,
		c: A.c * k,
		d: A.d * k,
	}
}

func (A *Mat2x2) TimesVec(V Vec2) Vec2 {
	return Vec2{
		X: A.a*V.X + A.b*V.Y,
		Y: A.c*V.X + A.d*V.Y,
	}
}

func Eye() Mat2x2 {
	return Mat2x2{
		a: 1,
		b: 0,
		c: 0,
		d: 1,
	}
}

func Mat2x2FromRot(theta float64) Mat2x2 {
	cos := math.Cos(theta)
	sin := math.Sin(theta)
	return Mat2x2{
		a: cos,
		b: -sin,
		c: sin,
		d: cos,
	}
}

func NewMatrix(a, b, c, d float64) Mat2x2 {
	return Mat2x2{
		a: a,
		b: b,
		c: c,
		d: d,
	}
}

type Vec2 struct {
	X, Y float64
}

func (v *Vec2) Plus(u Vec2) Vec2 {
	return Vec2{
		X: v.X + u.X,
		Y: v.Y + u.Y,
	}
}

func (v *Vec2) Dot(u Vec2) float64 {
	return v.X*u.X + v.Y*u.Y
}

func (v *Vec2) Ortho() Vec2 {
	return NewVec2(-v.Y, v.X)
}

func NewVec2(x, y float64) Vec2 {
	return Vec2{
		X: x,
		Y: y,
	}
}
