package arrays

// Given an array of structs, return an array of the
func Filter[V any](input []V, f func(V) bool) (result []V) {
	for _, item := range input {
		if f(item) {
			result = append(result, item)
		}
	}
	return
}

func MapSlice[V any, W any](input []V, f func(V) W) (result []W) {
	for _, item := range input {
		result = append(result, f(item))
	}
	return
}

func MapMap[V any, W any, X comparable](input map[X]V, f func(V) W) map[X]W {
	result := make(map[X]W, 0)
	for key, item := range input {
		result[key] = f(item)
	}
	return result
}

func MapToArray[V comparable, W any](input map[V]W) (valueSlice []W) {
	for _, value := range input {
		valueSlice = append(valueSlice, value)
	}
	return
}
