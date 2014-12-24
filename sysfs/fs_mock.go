package sysfs

import (
	"errors"
	"os"
	"time"
)

// A mock filesystem of simple files.
type MockFilesystem struct {
	Seq   int // Increases with each write or read.
	Files map[string]*MockFile
}

// A simple mock file that contains a single string.  Any write
// overwrites, and any read returns from the start.
type MockFile struct {
	Contents string
	Seq      int // When this file was last written or read.
	Opened   bool
	Closed   bool
	fd       uintptr

	fs *MockFilesystem
}

var _ File = (*MockFile)(nil)
var _ Filesystem = (*MockFilesystem)(nil)

func (f *MockFile) Write(b []byte) (n int, err error) {
	return f.WriteString(string(b))
}

func (f *MockFile) WriteString(s string) (ret int, err error) {
	f.Contents = s
	f.Seq = f.fs.next()
	return len(s), nil
}

func (f *MockFile) Sync() (err error) {
	return nil
}

func (f *MockFile) Read(b []byte) (n int, err error) {
	count := len(b)
	if len(f.Contents) < count {
		count = len(f.Contents)
	}
	copy(b, []byte(f.Contents)[:count])
	f.Seq = f.fs.next()

	return count, nil
}

func (f *MockFile) ReadAt(b []byte, off int64) (n int, err error) {
	return f.Read(b)
}

func (f *MockFile) Fd() uintptr {
	return f.fd
}

func (f *MockFile) Close() error {
	return nil
}

func NewMockFilesystem(files []string) *MockFilesystem {
	m := &MockFilesystem{
		Files: make(map[string]*MockFile),
	}

	for i := range files {
		m.Add(files[i])
	}

	return m
}

func (fs *MockFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	f, ok := fs.Files[name]
	if ok {
		f.Opened = true
		f.Closed = false
		return f, nil
	} else {
		return (*MockFile)(nil), &os.PathError{Err: errors.New(name + ": No such file.")}
	}
}

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
