package tests

func ToP[T any](elem T) *T {
	return &elem
}
