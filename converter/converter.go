package converter

import (
	"encoding/json"
	"errors"
	"golang.org/x/exp/constraints"
)

var (
	ErrUnsupportedTypeConversion = errors.New("unsupported type conversion")
)

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

// FromInterfaceTo converts data from type interface{} to any type given.
// The used switch case shows which type conversions are supported.
// If the type is not supported, ErrUnsupportedTypeConversion is thrown.
func FromInterfaceTo[T any](data interface{}) (T, error) {
	var convertedData T
	switch d := data.(type) {
	case string:
		if err := json.Unmarshal([]byte(d), &convertedData); err != nil {
			return convertedData, err
		}
		return convertedData, nil
	default:
		return convertedData, ErrUnsupportedTypeConversion
	}
}
