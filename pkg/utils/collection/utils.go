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

// MapWithError maps the array to another array
// with elements that were converted by provided function.
// If the provided function returns an error,
// MapWithError returns nil array and this error as a result.
//
// Example: Map([1, 2, 3], f(x) -> x*2) -> [2, 4, 6]
func MapWithError[T, U any](arr []T, f func(x T) (U, error)) ([]U, error) {
	result := make([]U, 0, len(arr))

	for _, x := range arr {
		mapped, err := f(x)
		if err != nil {
			return nil, err
		}

		result = append(result, mapped)
	}

	return result, nil
}

// Filter filters array by provided predicate.
//
// Example: Filter([1, 2, 3], isOdd) -> [1, 3]
func Filter[T any](arr []T, f func(x T) bool) []T {
	result := make([]T, 0, len(arr))

	for _, x := range arr {
		if f(x) {
			result = append(result, x)
		}
	}

	return result
}

// Keys returns a slice of keys retrieved from provided map.
//
// Example: Keys(map[1:1, 2:2, 3:4]) -> [1, 2, 3]
func Keys[T comparable, U any](items map[T]U) []T {
	result := make([]T, 0, len(items))

	for k := range items {
		result = append(result, k)
	}

	return result
}
