package cryptoutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// SignData signs the data using a private key and returns the signature.
func SignData(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hash[:], nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// VerifySignature verifies the signature of the data using a public key.
func VerifySignature(publicKey *rsa.PublicKey, data []byte, signature []byte) error {
	hash := sha256.Sum256(data)
	return rsa.VerifyPSS(publicKey, crypto.SHA256, hash[:], signature, nil)
}
