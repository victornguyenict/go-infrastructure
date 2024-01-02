package fileutil

import "path/filepath"

// GetAbsolutePath returns the absolute path of the specified file or directory.
func GetAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
