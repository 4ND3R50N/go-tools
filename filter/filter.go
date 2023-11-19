package filter

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
