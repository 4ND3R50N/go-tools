package converter

import "golang.org/x/exp/constraints"

// ToPointer converts any type (E) to pointer (*E)
func ToPointer[E any](toConvert E) *E {
	return &toConvert
}

// ToValue converts any pointer (*E) to value (E).
// Panics if toConvert is a nil pointer.
func ToValue[E any](toConvert *E) E {
	return *toConvert
}

// ToValueOrZero converts like ToValue but returns a zero value if toConvert is nil.
func ToValueOrZero[T constraints.Ordered](toConvert *T) T {
	if toConvert == nil {
		var zero T
		return zero
	}
	return ToValue(toConvert)
}

// ToPointerOrNil converts like ToPointer but returns nil if toConvert is a zero value.
func ToPointerOrNil[T constraints.Ordered](toConvert T) *T {
	var zero T
	if toConvert == zero {
		return nil
	}
	return ToPointer(toConvert)
}
