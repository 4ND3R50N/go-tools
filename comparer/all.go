package comparer

// All checks whether all elements of slice satisfy the given predicate.
func All[E any](slice []E, pred func(elem E) bool) bool {
	for _, elem := range slice {
		if !pred(elem) {
			return false
		}
	}
	return true
}
