// Package utils provides BASIC FUCKING utilities 
package utils

func Select[T any, U any](input []T, mapper func(T) U) []U {
	output := make([]U, len(input))
	for i, item := range input {
		output[i] = mapper(item)
	}
	return output
}

func Filter[T any](input []T, predicate func(T) bool) []T {
	output := make([]T, 0)
	for _, item := range input {
		if predicate(item) {
			output = append(output, item)
		}
	}
	return output
}
