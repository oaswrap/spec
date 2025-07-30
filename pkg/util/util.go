package util

func Optional[T any](defaultValue T, value ...T) T {
	if len(value) > 0 {
		return value[0]
	}
	return defaultValue
}

func PtrOf[T any](value T) *T {
	return &value
}
