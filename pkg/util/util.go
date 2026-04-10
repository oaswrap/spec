package util

import (
	"path"
)

// Optional returns the first value from the provided values or the default value if no values are provided.
func Optional[T any](defaultValue T, value ...T) T {
	if len(value) > 0 {
		return value[0]
	}
	return defaultValue
}

// PtrOf returns a pointer to the provided value.
func PtrOf[T any](value T) *T {
	return &value
}

// JoinPath joins multiple path segments into a single path, ensuring proper formatting.
func JoinPath(paths ...string) string {
	if len(paths) == 0 {
		return ""
	}

	lastElement := paths[len(paths)-1]
	if len(lastElement) > 0 && lastElement[len(lastElement)-1] == '/' {
		return path.Join(paths...) + "/"
	}
	return path.Join(paths...)
}
