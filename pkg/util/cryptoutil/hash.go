package cryptoutil

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// HashStringMD5 takes an input string and returns its MD5 hash.
func HashStringMD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// HashStringSHA256 takes an input string and returns its SHA-256 hash.
func HashStringSHA256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
