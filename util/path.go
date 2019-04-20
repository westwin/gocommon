package util

import (
	"strings"
)

// Join to append path to base path
func Join(base, path string) string {
	b := strings.HasSuffix(base, "/")
	p := strings.HasPrefix(path, "/")
	switch {
	case b && p:
		return base + path[1:]
	case b || p:
		return base + path
	default:
		return base + "/" + path
	}
}
