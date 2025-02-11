package utils

import (
	"log"
	"math"
)

// TODO: improve this so that it at least gives some info about where the panicking
// call came from. This is actually quite important
// This should ONLY be used when you are almost completely certain that the function
// being calles will not throw an error
func HandleLazy(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type FloatConvertible interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

func F64[F FloatConvertible](num F) float64 {
	return float64(num)
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
	return int(math.Min(F64(a), F64(b)))
}

func MaxInt(a, b int) int {
	return int(math.Max(F64(a), F64(b)))
}
