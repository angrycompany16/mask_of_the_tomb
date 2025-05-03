package timeutil

import "time"

var start time.Time

func Init() {
	start = time.Now()
}

func GetTime() float64 {
	return float64(time.Since(start).Microseconds())
}
