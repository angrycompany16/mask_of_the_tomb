package maths

import (
	"image/color"
	"math"
	"math/rand/v2"
)

type Direction int

const (
	DirNone Direction = iota - 1
	DirUp
	DirDown
	DirLeft
	DirRight
)

func Opposite(dir Direction) Direction {
	switch dir {
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

func ToRadians(dir Direction) float64 {
	switch dir {
	case DirUp:
		return 0
	case DirDown:
		return math.Pi
	case DirRight:
		return math.Pi / 2
	case DirLeft:
		return 3 * math.Pi / 2
	}
	return 0
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

func Mix(a, b color.Color, t float64) color.Color {
	rA, gA, bA, aA := a.RGBA()
	rB, gB, bB, aB := b.RGBA()
	return color.RGBA64{
		uint16(Lerp(float64(rA), float64(rB), t)),
		uint16(Lerp(float64(gA), float64(gB), t)),
		uint16(Lerp(float64(bA), float64(bB), t)),
		uint16(Lerp(float64(aA), float64(aB), t)),
	}
}

func RandomRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
