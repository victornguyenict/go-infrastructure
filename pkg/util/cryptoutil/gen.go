package cryptoutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"io"
)

// GenerateRandomBytes returns securely generated random bytes.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString generates a secure random string of the desired length.
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

// GenerateRSAKeys generates a new RSA private key and returns it along with its public key.
func GenerateRSAKeys(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// GenerateSalt creates a cryptographically secure random salt.
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)

	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}
