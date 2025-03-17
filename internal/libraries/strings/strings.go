package strings

import (
	"fmt"
	"math"
	"strconv"
)

func Increment(s string) string {
	var total int64
	var numCount int
	for i := range len(s) {
		if IsAlphaNumeric(s[len(s)-i-1]) {
			num, _ := strconv.ParseInt(string(s[len(s)-i-1]), 10, 64)
			total += num * int64(math.Pow10(i))
			numCount++
		}
	}
	total += 1
	return fmt.Sprintf("%s%d", s[len(s)-numCount:], total)
}

func IsAlphaNumeric(c byte) bool {
	// Check if the byte value falls within the range of alphanumeric characters
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
