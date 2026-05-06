package maths

import "math/rand/v2"

type RandomFloat64 struct {
	Min float64
	Max float64
}

func (n *RandomFloat64) Eval() float64 {
	return Lerp(n.Min, n.Max, rand.Float64())
}

type RandomInt struct {
	Min int
	Max int
}

func (n *RandomInt) Eval() int {
	r := rand.Float64()
	return int(float64(n.Min)*r + float64(n.Max)*(1-r))
}
