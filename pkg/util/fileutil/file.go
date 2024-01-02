package fileutil

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// ReadFile reads the content of the file specified by the filePath.
func ReadFile(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// WriteFile writes the data to the file specified by the filePath.
func WriteFile(filePath string, data string) error {
	return ioutil.WriteFile(filePath, []byte(data), 0644)
}

// FileExists checks if a file exists and is not a directory before we try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// DeleteFile deletes the specified file.
func DeleteFile(filePath string) error {
	return os.Remove(filePath)
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are the same, then return success.
// Otherwise, attempt to create a copy of the file.
func CopyFile(src, dst string) error {
	srcFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcFileStat.Mode().IsRegular() {
		return ErrNonRegularSourceFile(srcFileStat)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

// MoveFile moves a file from src to dst. Can be used to rename files or move them to a different directory.
func MoveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// GetFileSize returns the size of the file specified.
func GetFileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// ReadBinaryFile reads the entire file into memory as bytes.
func ReadBinaryFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// WriteBinaryFile writes data to a file as bytes.
func WriteBinaryFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

// ChangeFilePermissions changes the permissions of the specified file.
func ChangeFilePermissions(filename string, perms os.FileMode) error {
	return os.Chmod(filename, perms)
}

// FindFiles walks the directory tree and finds files that match a certain pattern.
func FindFiles(rootDir string, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// BatchRemoveFiles removes multiple files in a single operation.
func BatchRemoveFiles(files []string) error {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	return nil
}

// CompareFiles checks if two files have the same content.
func CompareFiles(file1, file2 string) (bool, error) {
	bytes1, err := ioutil.ReadFile(file1)
	if err != nil {
		return false, err
	}
	bytes2, err := ioutil.ReadFile(file2)
	if err != nil {
		return false, err
	}
	return bytes.Equal(bytes1, bytes2), nil
}

// WalkFilesWithFunc applies a given function to all files in a directory tree.
func WalkFilesWithFunc(root string, fileFunc filepath.WalkFunc) error {
	return filepath.Walk(root, fileFunc)
}

// GetFileNamesInDir returns the filenames (not including path) of all files in a directory.
func GetFileNamesInDir(dirPath string) ([]string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames, nil
}

// DiffFiles compares two files and returns true if they are different.
// This function does a line-by-line comparison.
func DiffFiles(file1Path, file2Path string) (bool, error) {
	file1, err := os.Open(file1Path)
	if err != nil {
		return false, err
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		return false, err
	}
	defer file2.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	for scanner1.Scan() && scanner2.Scan() {
		if scanner1.Text() != scanner2.Text() {
			return true, nil // Files are different
		}
	}

	// Check if files have the same length
	if scanner1.Scan() || scanner2.Scan() {
		return true, nil // One file had more lines than the other
	}

	return false, nil // Files are the same
}

// TouchFile updates the access and modification times of the file,
// similar to the Unix command `touch`.
func TouchFile(filePath string) error {
	now := time.Now()
	return os.Chtimes(filePath, now, now)
}
