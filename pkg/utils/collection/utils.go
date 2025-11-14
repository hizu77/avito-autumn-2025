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

// Reduce accumulates some value over the passed array.
// JS analogue: https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/reduce
//
// Example: Reduce([1, 2, 3, 4], f(accumulator, x) -> accumulator + x, 0) -> 10
func Reduce[T, U any](
	arr []T,
	reducer func(U, T) U,
	initialValue U,
) U {
	accumulator := initialValue

	for _, item := range arr {
		accumulator = reducer(accumulator, item)
	}

	return accumulator
}
