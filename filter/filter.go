package filter

import "github.com/4ND3R50N/go-tools/comparer"

// Filter gets a slice and a function where you can define custom filter conditions to get a specific subset of the
// input slice.
func Filter[T any](slice []T, filter func(element T) bool) []T {
	var ret = make([]T, 0)
	for _, item := range slice {
		if filter(item) {
			ret = append(ret, item)
		}
	}
	return ret
}

// Distinct returns a slice of distinct values from the input slice.
// Example: It returns [1 2 3] for input [1 1 1 2 2 3].
func Distinct[T comparable](slice []T) []T {
	ret := make([]T, 0)
	for _, item := range slice {
		item := item
		if !comparer.Any(ret, func(ret_item T) bool {
			return ret_item == item
		}) {
			ret = append(ret, item)
		}
	}
	return ret
}
