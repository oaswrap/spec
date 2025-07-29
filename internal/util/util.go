package util

import "regexp"

func Optional[T any](defaultValue T, value ...T) T {
	if len(value) > 0 {
		return value[0]
	}
	return defaultValue
}

// ConvertPath converts ":param" to "{param}" anywhere in the path.
func ConvertPath(path string) string {
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	return re.ReplaceAllString(path, `{$1}`)
}
