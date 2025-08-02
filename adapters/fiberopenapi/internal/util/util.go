package util

import "regexp"

// ConvertPath converts ":param" to "{param}" anywhere in the path.
func ConvertPath(path string) string {
	re := regexp.MustCompile(`:([a-zA-Z0-9_]+)`)
	return re.ReplaceAllString(path, `{$1}`)
}
