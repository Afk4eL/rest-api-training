package csrf

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	CSRFHeader = "X-CSRF-Token"

	tokenLength = 32
)

// TODO:logger
func MakeToken() (string, error) {
	bytes := make([]byte, tokenLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", nil
	}

	return base64.URLEncoding.EncodeToString(bytes)[:tokenLength], nil
}

func ValidateToken() {

}
