package system

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

var (
	_ File       = (*MockFile)(nil)
	_ filesystem = (*MockFilesystem)(nil)
)

// MockFilesystem represents a filesystem of mock files.
type MockFilesystem struct {
	Seq            int // Increases with each write or read.
	Files          map[string]*MockFile
	WithReadError  bool
	WithWriteError bool
	WithCloseError bool
	numCallsWrite  int
	numCallsRead   int
}

// A MockFile represents a mock file that contains a single string.  Any write
// overwrites, and any read returns from the start.
type MockFile struct {
	Contents           string
	Seq                int // When this file was last written or read.
	Opened             bool
	Closed             bool
	fd                 uintptr
	simulateWriteError error
	simulateReadError  error

	fs *MockFilesystem
}

var (
	errRead  = fmt.Errorf("read error")
	errWrite = fmt.Errorf("write error")
	errClose = fmt.Errorf("close error")
)

// Write writes string(b) to f.Contents
func (f *MockFile) Write(b []byte) (int, error) {
	if f.fs.WithWriteError {
		return 0, errWrite
	}

	return f.WriteString(string(b))
}

// Seek seeks to a specific offset in a file
func (f *MockFile) Seek(offset int64, whence int) (int64, error) {
	return offset, nil
}

// WriteString writes s to f.Contents
func (f *MockFile) WriteString(s string) (int, error) {
	f.Contents = s
	f.Seq = f.fs.next()
	f.fs.numCallsWrite++
	return len(s), f.simulateWriteError
}

// Sync implements the File interface Sync function
func (f *MockFile) Sync() error {
	return nil
}

// Read copies b bytes from f.Contents
func (f *MockFile) Read(b []byte) (int, error) {
	if f.fs.WithReadError {
		return 0, errRead
	}

	count := len(b)
	if len(f.Contents) < count {
		// note: if length of content is smaller than given b than b will not completely overridden, e.g. if content
		// is empty, than b is not changed at all. This is fine for this MockFile, but in real world, on different
		// length, an error will be created.
		count = len(f.Contents)
	}
	copy(b, []byte(f.Contents)[:count])
	f.Seq = f.fs.next()

	f.fs.numCallsRead++
	return count, f.simulateReadError
}

// ReadAt calls MockFile.Read
func (f *MockFile) ReadAt(b []byte, off int64) (int, error) {
	return f.Read(b)
}

// Fd returns a random uintprt based on the time of the MockFile creation
func (f *MockFile) Fd() uintptr {
	return f.fd
}

// Close implements the File interface Close function
func (f *MockFile) Close() error {
	if f != nil {
		f.Opened = false
		f.Closed = true
		if f.fs != nil && f.fs.WithCloseError {
			f.Closed = false
			return errClose
		}
	}
	return nil
}

// newMockFilesystem returns a new MockFilesystem given a list of file and folder paths
func newMockFilesystem(items []string) *MockFilesystem {
	m := &MockFilesystem{
		Files: make(map[string]*MockFile),
	}

	for _, item := range items {
		m.Add(item)
	}

	return m
}

// OpenFile opens file name from fs.Files, if the file does not exist it returns an os.PathError
func (fs *MockFilesystem) openFile(name string, _ int, _ os.FileMode) (File, error) {
	f, ok := fs.Files[name]
	if ok {
		f.Opened = true
		f.Closed = false
		return f, nil
	}

	return (*MockFile)(nil), &os.PathError{Err: fmt.Errorf("%s: no such file", name)}
}

// Stat returns a generic FileInfo for all files in fs.Files.
// If the file does not exist it returns an os.PathError
func (fs *MockFilesystem) stat(name string) (os.FileInfo, error) {
	_, ok := fs.Files[name]
	if ok {
		// return file based mock FileInfo, CreateTemp don't like "/" in between
		tmpFile, err := os.CreateTemp("", strings.ReplaceAll(name, "/", "_"))
		if err != nil {
			log.Println("A")
			return nil, err
		}
		defer func() {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
		}()

		return os.Stat(tmpFile.Name())
	}

	dirName := name + "/"
	for path := range fs.Files {
		if strings.HasPrefix(path, dirName) {
			// return dir based mock FileInfo, TempDir don't like "/" in between
			tmpDir, err := os.MkdirTemp("", strings.ReplaceAll(name, "/", "_"))
			if err != nil {
				return nil, err
			}
			defer os.RemoveAll(tmpDir)

			return os.Stat(tmpDir)
		}
	}

	return nil, &os.PathError{Err: fmt.Errorf("%s: no such file", name)}
}

// Find returns all items (files or folders) below the given directory matching the given pattern.
func (fs *MockFilesystem) find(baseDir string, pattern string) ([]string, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	var found []string
	for name := range fs.Files {
		if !strings.HasPrefix(name, baseDir) {
			continue
		}
		item := strings.TrimPrefix(name[len(baseDir):], "/")

		firstItem := strings.Split(item, "/")[0]
		if reg.MatchString(firstItem) {
			found = append(found, path.Join(baseDir, firstItem))
		}
	}
	return found, nil
}

// readFile returns the contents of the given file. If the file does not exist it returns an os.PathError
func (fs *MockFilesystem) readFile(name string) ([]byte, error) {
	if fs.WithReadError {
		return nil, errRead
	}

	f, ok := fs.Files[name]
	if !ok {
		return nil, &os.PathError{Err: fmt.Errorf("%s: no such file", name)}
	}

	f.fs.numCallsRead++
	return []byte(f.Contents), nil
}

func (fs *MockFilesystem) next() int {
	fs.Seq++
	return fs.Seq
}

// Add adds a new file to fs.Files given a name, and returns the newly created file
func (fs *MockFilesystem) Add(name string) *MockFile {
	f := &MockFile{
		Seq: -1,
		fd:  uintptr(time.Now().UnixNano() & 0xffff),
		fs:  fs,
	}
	fs.Files[name] = f
	return f
}
