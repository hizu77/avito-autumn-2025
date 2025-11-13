package collection

// Map maps the array to another array
// with elements that were converted by provided function.
//
// Example: Map([1, 2, 3], f(x) -> x*2) -> [2, 4, 6]
func Map[T, U any](arr []T, f func(x T) U) []U {
	result := make([]U, 0, len(arr))

	for _, x := range arr {
		result = append(result, f(x))
	}

	return result
}
