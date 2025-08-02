package util

import (
	"path"
	"strings"
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

// JoinURL joins the base URL with the provided segments, ensuring proper formatting.
func JoinURL(base string, segments ...string) string {
	base = strings.TrimRight(base, "/")
	if len(segments) == 0 {
		return base
	}
	return base + "/" + strings.TrimLeft(JoinPath(segments...), "/")
}

// JoinPath joins multiple path segments into a single path, ensuring proper formatting.
func JoinPath(paths ...string) string {
	return path.Join(paths...)
}
