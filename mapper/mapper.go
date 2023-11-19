package mapper

// Map can be used to convert any list from MemberType E to T.
// It takes a list of E (E can also be a *E) and a function that is used for each element of E to T.
// It returns a []T after all.
//
// Mapping: [-][-][-] -> [+][+][+]
func Map[E any, T any](from []E, fn func(fromEntry E) T) []T {
	converted := make([]T, len(from))
	for i, fromEntry := range from {
		toEntry := fn(fromEntry)
		converted[i] = toEntry
	}
	return converted
}

// FlatMap can be used to convert any list from MemberType E to T.
// It takes a list of E (E can also be a *E) and a function that is used for each element of E to []T.
// Instead of returning [][]T, the resulting list gets flattened to []T.
// Mapping: [[-][-][-]] [[-][-][-]] -> [+][+][+][+][+][+]
func FlatMap[E any, T any](from []E, fn func(fromEntry E) []T) []T {
	var converted []T
	for _, fromEntry := range from {
		toEntry := fn(fromEntry)
		converted = append(converted, toEntry...)
	}
	return converted
}

// MapWithErr behaves like PointerMap, but stops mapping when fn returns an error.
// The first error that occurs is returned.
func MapWithErr[E any, T any](from []E, fn func(fromEntry E) (T, error)) ([]T, error) {
	converted := make([]T, len(from))
	for i, fromEntry := range from {
		toEntry, err := fn(fromEntry)
		if err != nil {
			return converted, err
		}
		converted[i] = toEntry
	}
	return converted, nil
}
