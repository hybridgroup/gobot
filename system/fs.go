package system

import (
	"os"
	"path"
	"regexp"
)

// nativeFilesystem represents the native file system implementation
type nativeFilesystem struct{}

// openFile calls os.OpenFile().
func (fs *nativeFilesystem) openFile(name string, flag int, perm os.FileMode) (File, error) {
	return os.OpenFile(name, flag, perm)
}

// stat calls os.Stat()
func (fs *nativeFilesystem) stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// find returns all items (files or folders) below the given directory matching the given pattern.
func (fs *nativeFilesystem) find(baseDir string, pattern string) ([]string, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	items, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var found []string
	for _, item := range items {
		if reg.MatchString(item.Name()) {
			found = append(found, path.Join(baseDir, item.Name()))
		}
	}
	return found, nil
}

// readFile reads the named file and returns the contents. A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func (fs *nativeFilesystem) readFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}
