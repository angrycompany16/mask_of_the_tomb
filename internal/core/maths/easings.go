package maths

import "math"

func SineInOut(x float64) float64 {
	return -(math.Cos(math.Pi*x) - 1) / 2
}

func ExpInOut(x float64) float64 {
	if x == 0 {
		return 0
	} else if x == 1 {
		return 1
	}

	if x < 0.5 {
		return math.Pow(2, 20*x-10) / 2
	}

	return (2 - math.Pow(2, -20*x+10)) / 2
}

func QuadInOut(x float64) float64 {
	if x < 0.5 {
		return 2 * x * x
	}
	return 1 - math.Pow(-2*x+2, 2)/2
}

func QuartInOut(x float64) float64 {
	if x < 0.5 {
		return 8 * math.Pow(x, 4)
	}
	return 1 - math.Pow(-2*x+2, 4)/2
}

func QuadIn(x float64) float64 {
	return math.Pow(x, 2)
}

func QuadOut(x float64) float64 {
	return 1 - math.Pow(1-x, 2)
}
