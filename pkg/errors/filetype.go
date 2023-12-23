package errors

import (
	"fmt"
	"github.com/eliofery/golang-image/pkg/static"
	"io"
	"net/http"
	"path/filepath"
)

type FileError struct {
	Issue string
}

func (ft FileError) Error() string {
	return fmt.Sprintf("не верный тип файла: %s", ft.Issue)
}

func checkContentType(r io.ReadSeeker, allowedTypes []string) error {
	op := "filetype.checkContentType"

	testBytes := make([]byte, 512)
	_, err := r.Read(testBytes)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	contentType := http.DetectContentType(testBytes)
	for _, allowedType := range allowedTypes {
		if allowedType == contentType {
			return nil
		}
	}

	return FileError{
		Issue: fmt.Sprintf("не верный формат файла: %s", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if static.HasExtension(filename, allowedExtensions) {
		return nil
	}

	return FileError{
		Issue: fmt.Sprintf("не верное расширение файла: %s", filepath.Ext(filename)),
	}
}
