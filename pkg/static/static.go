package static

import (
	"path/filepath"
	"strings"
)

func HasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)

		if filepath.Ext(file) == ext {
			return true
		}
	}

	return false
}
