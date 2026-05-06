package maths

import "math/rand/v2"

type RandomFloat struct {
	Min float64
	Max float64
}

func (n *RandomFloat) Eval() float64 {
	return Lerp(n.Min, n.Max, rand.Float64())
}

func NewRandomFloat(min, max float64) RandomFloat {
	return RandomFloat{
		Min: min,
		Max: max,
	}
}

type RandomInt struct {
	Min int
	Max int
}

func (n *RandomInt) Eval() int {
	r := rand.Float64()
	return int(float64(n.Min)*r + float64(n.Max)*(1-r))
}

func NewRandomInt(min, max int) RandomInt {
	return RandomInt{
		Min: min,
		Max: max,
	}
}
