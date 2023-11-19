package converter

// ConvertToPointer converts any type (E) to pointer (*E)
func ConvertToPointer[E any](toConvert E) *E {
	return &toConvert
}
