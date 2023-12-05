package rand

import (
	"crypto/sha256"
	"encoding/base64"
)

const (
	SessionTokenBytes = 32
)

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}

func HashToken(token string) string {
	tokenHash := sha256.Sum256([]byte(token))

	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
