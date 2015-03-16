package sysfs

import (
	"errors"
	"os"
	"time"
)

var _ File = (*MockFile)(nil)
var _ Filesystem = (*MockFilesystem)(nil)

// MockFilesystem represents  a filesystem of mock files.
type MockFilesystem struct {
	Seq   int // Increases with each write or read.
	Files map[string]*MockFile
}

// A MockFile represents a mock file that contains a single string.  Any write
// overwrites, and any read returns from the start.
type MockFile struct {
	Contents string
	Seq      int // When this file was last written or read.
	Opened   bool
	Closed   bool
	fd       uintptr

	fs *MockFilesystem
}

// Write writes string(b) to f.Contents
func (f *MockFile) Write(b []byte) (n int, err error) {
	return f.WriteString(string(b))
}

// WriteString writes s to f.Contents
func (f *MockFile) WriteString(s string) (ret int, err error) {
	f.Contents = s
	f.Seq = f.fs.next()
	return len(s), nil
}

// Sync implements the File interface Sync function
func (f *MockFile) Sync() (err error) {
	return nil
}

// Read copies b bytes from f.Contents
func (f *MockFile) Read(b []byte) (n int, err error) {
	count := len(b)
	if len(f.Contents) < count {
		count = len(f.Contents)
	}
	copy(b, []byte(f.Contents)[:count])
	f.Seq = f.fs.next()

	return count, nil
}

// ReadAt calls MockFile.Read
func (f *MockFile) ReadAt(b []byte, off int64) (n int, err error) {
	return f.Read(b)
}

// Fd returns a random uintprt based on the time of the MockFile creation
func (f *MockFile) Fd() uintptr {
	return f.fd
}

// Close implements the File interface Close function
func (f *MockFile) Close() error {
	return nil
}

// NewMockFilesystem returns a new MockFilesystem given a list of file paths
func NewMockFilesystem(files []string) *MockFilesystem {
	m := &MockFilesystem{
		Files: make(map[string]*MockFile),
	}

	for i := range files {
		m.Add(files[i])
	}

	return m
}

// OpenFile opens file name from fs.Files, if the file does not exist it returns an os.PathError
func (fs *MockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	f, ok := fs.Files[name]
	if ok {
		f.Opened = true
		f.Closed = false
		return f, nil
	}
	return (*MockFile)(nil), &os.PathError{Err: errors.New(name + ": No such file.")}
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

func (fs *MockFilesystem) next() int {
	fs.Seq++
	return fs.Seq
}
