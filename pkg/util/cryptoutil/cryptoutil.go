package cryptoutil

import (
	"crypto/hmac"
	"crypto/sha256"
)

// CreateHMAC creates a HMAC hash using SHA256 and returns it.
func CreateHMAC(secret, message []byte) ([]byte, error) {
	hmacHash := hmac.New(sha256.New, secret)

	_, err := hmacHash.Write(message)
	if err != nil {
		return nil, err
	}

	return hmacHash.Sum(nil), nil
}
