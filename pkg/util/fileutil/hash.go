package fileutil

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// CalculateFileMD5 returns the MD5 hash of the file content.
func CalculateFileMD5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = fmt.Sprintf("%x", hashInBytes)

	return returnMD5String, nil
}

// CalculateFileSHA1 returns the SHA1 hash of the file content.
func CalculateFileSHA1(filePath string) (string, error) {
	var returnSHA1String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA1String, err
	}
	defer file.Close()

	hash := sha1.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA1String, err
	}

	hashInBytes := hash.Sum(nil)[:20]
	returnSHA1String = fmt.Sprintf("%x", hashInBytes)

	return returnSHA1String, nil
}

// CalculateFileSHA256 returns the SHA256 hash of the file content.
func CalculateFileSHA256(filePath string) (string, error) {
	var returnSHA256String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnSHA256String, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnSHA256String, err
	}

	hashInBytes := hash.Sum(nil)[:32]
	returnSHA256String = fmt.Sprintf("%x", hashInBytes)

	return returnSHA256String, nil
}
