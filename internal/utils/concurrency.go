package utils

func PollThread[T any](out <-chan T) (T, bool) {
	select {
	case a := <-out:
		return a, true
	default:
		var zero T
		return zero, false
	}
}
