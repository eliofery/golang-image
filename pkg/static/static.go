package static

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GalleriesDir(id uint) string {
	return filepath.Join(os.Getenv("STATIC_DIR"), fmt.Sprintf("galleries/%d", id))
}

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
