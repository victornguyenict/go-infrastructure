package fileutil

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ErrNonRegularSourceFile error returned when source file is not a regular file.
func ErrNonRegularSourceFile(fileStat os.FileInfo) error {
	// Implement detailed error message or custom error type
	return os.ErrInvalid
}

// CreateTempFile creates a new temporary file in the directory dir
// with a name beginning with prefix, opens the file for reading and writing,
// and returns the resulting *os.File.
func CreateTempFile(dir, prefix string) (*os.File, error) {
	tempFile, err := ioutil.TempFile(dir, prefix)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

// WalkDirectory traverses the directory tree rooted at root, calling walkFn for each file or directory in the tree, including root.
func WalkDirectory(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}

// GetExecutablePath returns the absolute path of the current executable.
func GetExecutablePath() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(executable)
}

// ListEnvironmentPaths returns a slice of paths from the PATH environment variable.
func ListEnvironmentPaths() []string {
	pathVar := os.Getenv("PATH")
	paths := filepath.SplitList(pathVar)
	return paths
}

// IsHidden checks if a given file or directory is hidden.
func IsHidden(filePath string) (bool, error) {
	fileName := filepath.Base(filePath)
	if fileName != "." && strings.HasPrefix(fileName, ".") {
		return true, nil
	}
	return false, nil
}

// RenameFilesWithPattern renames all files in a directory that match a certain pattern.
func RenameFilesWithPattern(dirPath, pattern, replaceWith string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, file := range files {
		if strings.Contains(file.Name(), pattern) {
			newName := strings.Replace(file.Name(), pattern, replaceWith, -1)
			oldPath := filepath.Join(dirPath, file.Name())
			newPath := filepath.Join(dirPath, newName)
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
		}
	}
	return nil
}

// FindFilesMatchingPattern returns a list of file paths that match the given pattern in the specified directory.
func FindFilesMatchingPattern(dirPath, pattern string) ([]string, error) {
	var matches []string
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if matched, _ := filepath.Match(pattern, file.Name()); matched {
			matches = append(matches, filepath.Join(dirPath, file.Name()))
		}
	}
	return matches, nil
}

// GrepFile searches for a string pattern in a file and returns the lines containing the pattern.
func GrepFile(filePath string, pattern string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), pattern) {
			lines = append(lines, scanner.Text())
		}
	}
	return lines, scanner.Err()
}

// CountLines returns the number of lines in the specified file.
func CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// ConcatenateFiles concatenates multiple files into one destination file.
func ConcatenateFiles(destPath string, filePaths []string) error {
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	for _, filePath := range filePaths {
		srcFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(destFile, srcFile)
		srcFile.Close() // close srcFile on each loop
		if err != nil {
			return err
		}
	}
	return nil
}

// ListFilesRecursive lists all files in a directory and its subdirectories.
func ListFilesRecursive(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ConvertUnixToDos converts Unix line endings (LF) to DOS/Windows line endings (CRLF).
func ConvertUnixToDos(filePath string) error {
	// omitted for brevity - implement conversion from LF to CRLF
	return nil
}

// ConvertDosToUnix converts DOS/Windows line endings (CRLF) to Unix line endings (LF).
func ConvertDosToUnix(filePath string) error {
	// omitted for brevity - implement conversion from CRLF to LF
	return nil
}

// ChownRecursive changes the owner and group of the specified directory and all its sub-content.
func ChownRecursive(path string, uid, gid int) error {
	return filepath.Walk(path, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(name, uid, gid)
	})
}

// RemoveEmptyDirs recursively removes all empty subdirectories in a given directory.
func RemoveEmptyDirs(dirPath string) error {
	dirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}
	for _, d := range dirs {
		if d.IsDir() {
			fullDirPath := filepath.Join(dirPath, d.Name())
			err := RemoveEmptyDirs(fullDirPath) // recursive call
			if err != nil {
				return err
			}
			subDirs, _ := ioutil.ReadDir(fullDirPath)
			if len(subDirs) == 0 {
				os.Remove(fullDirPath) // remove empty dir
			}
		}
	}
	return nil
}

// SimplifyPath cleans up a file path, evaluating any . or .. elements and returning a clean, absolute path.
func SimplifyPath(path string) (string, error) {
	return filepath.Abs(filepath.Clean(path))
}

// CopyIfNewer copies the source file to destination only if the source is newer than the destination.
func CopyIfNewer(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	dstInfo, err := os.Stat(dst)
	// If destination doesn't exist or source is newer, perform the copy
	if os.IsNotExist(err) || srcInfo.ModTime().After(dstInfo.ModTime()) {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}
	return nil
}
