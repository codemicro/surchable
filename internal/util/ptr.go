package util

func Ptr[T any](x T) *T {
	return &x
}

func PtrNilIfDefault[T comparable](x T) *T {
	var n T
	if x == n {
		return nil
	}
	return &x
}
