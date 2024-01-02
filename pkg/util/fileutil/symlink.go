package fileutil

import "os"

// CreateSymLink creates a symbolic link at linkPath pointing to the targetPath.
func CreateSymLink(targetPath, linkPath string) error {
	return os.Symlink(targetPath, linkPath)
}

// IsSymLink checks if the given path is a symbolic link.
func IsSymLink(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.Mode()&os.ModeSymlink != 0, nil
}

// ReadSymLink reads the target of a symbolic link.
func ReadSymLink(path string) (string, error) {
	return os.Readlink(path)
}
