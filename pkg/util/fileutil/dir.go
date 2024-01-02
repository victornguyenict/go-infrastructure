package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// CreateDirectory creates a new directory with the specified permissions.
func CreateDirectory(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// RemoveDirectory deletes a directory and all its contents.
func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}

// ReadDir reads the directory named by dirname and returns a list of directory entries sorted by filename.
func ReadDir(dirname string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	return files, nil
}

// EnsureDir checks if a directory exists and if not, creates it with the specified permissions.
func EnsureDir(dirName string, perms os.FileMode) error {
	err := os.MkdirAll(dirName, perms)
	if err != nil {
		return err
	}
	return nil
}

// CopyDirectory recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDirectory(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return ErrSourceNotDirectory(src)
	}

	_, err = os.Stat(dst)
	if !os.IsNotExist(err) {
		return ErrDestinationExists(dst)
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CalculateDirSize returns the total size of files in the specified directory.
func CalculateDirSize(dirPath string) (int64, error) {
	var size int64
	err := filepath.Walk(dirPath, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// ErrSourceNotDirectory error returned when source is not a directory.
func ErrSourceNotDirectory(src string) error {
	// Implement detailed error message or custom error type
	return os.ErrInvalid
}

// ErrDestinationExists error returned when destination already exists.
func ErrDestinationExists(dst string) error {
	// Implement detailed error message or custom error type
	return os.ErrInvalid
}

// IsDirectoryWritable checks if a directory is writable.
func IsDirectoryWritable(path string) bool {
	testFile := filepath.Join(path, ".tmp_writable_check")
	f, err := os.Create(testFile)
	if err != nil {
		return false
	}
	f.Close()
	os.Remove(testFile) // Clean up
	return true
}
