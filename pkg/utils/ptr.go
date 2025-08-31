package utils

func SPtr(v string) *string {
	return &v
}

func Ptr[T any](input T) *T {
	return &input
}
