package maths

import "math/rand/v2"

type RandomFloat64 struct {
	Min float64 `yaml:"Min"`
	Max float64 `yaml:"Max"`
}

func (n *RandomFloat64) Eval() float64 {
	return Lerp(n.Min, n.Max, rand.Float64())
}
