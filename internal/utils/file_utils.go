package utils

import "strings"

func HasAllowedExt(filename string, exts []string) bool {
	f := strings.ToLower(filename)
	for _, ext := range exts {
		if strings.HasSuffix(f, ext) {
			return true
		}
	}
	return false
}
