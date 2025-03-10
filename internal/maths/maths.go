package maths

import "math"

type Direction int

const (
	DirNone Direction = iota - 1
	DirUp
	DirDown
	DirLeft
	DirRight
)

func OppositeDir(dir Direction) Direction {
	switch dir {
	case DirNone:
		return DirNone
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirRight:
		return DirLeft
	case DirLeft:
		return DirRight
	}
	return DirNone
}

func Clamp(val, min, max float64) float64 {
	if val < min {
		return min
	} else if val > max {
		return max
	}
	return val
}

// The humble lerp
func Lerp(a, b, t float64) float64 {
	return a*(1.0-t) + b*t
}

func Mod(x, m int) int {
	return (x%m + m) % m
}

func MinInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func MaxInt(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}
