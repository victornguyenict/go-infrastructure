package cryptoutil

import "encoding/base64"

// Base64Encode returns the base64 encoding of input bytes.
func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode returns the bytes represented by the base64 string.
func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}
