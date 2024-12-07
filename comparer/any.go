package comparer

// Any checks whether any element of slice satisfies the given predicate.
func Any[E any](slice []E, pred func(elem E) bool) bool {
	for _, elem := range slice {
		if pred(elem) {
			return true
		}
	}
	return false
}
