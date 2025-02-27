package utils

import (
	"fmt"
	"log"
	"math"
	"path/filepath"
	"runtime"
)

type Direction int

const (
	DirNone Direction = iota - 1
	DirUp
	DirDown
	DirLeft
	DirRight
)

// This should ONLY be used when you are almost completely certain that the function
// being calles will not throw an error
func HandleLazy(err error) {
	pc, file, no, ok := runtime.Caller(1)
	funcDetails := runtime.FuncForPC(pc)
	if err != nil {
		if ok {
			fmt.Println(no)
			fmt.Printf("Failure from %s, file %s, line number %d\n", filepath.Base((funcDetails.Name())), file, no)
		}
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
