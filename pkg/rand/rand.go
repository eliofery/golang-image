package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	op := "rand.Bytes"

	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return b, nil
}

func String(n int) (string, error) {
	b, err := Bytes(n)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
