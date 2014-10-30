package mocks

import (
	"errors"
	"github.com/hybridgroup/gobot/internal"
	"os"
)

// A mock filesystem of simple files.
type Filesystem struct {
	Seq   int // Increases with each write or read.
	Files map[string]*File
}

// A simple mock file that contains a single string.  Any write
// overwrites, and any read returns from the start.
type File struct {
	Contents string
	Seq      int // When this file was last written or read.
	Opened   bool
	Closed   bool

	fs *Filesystem
}

var _ internal.File = (*File)(nil)
var _ internal.Filesystem = (*Filesystem)(nil)

func (f *File) Write(b []byte) (n int, err error) {
	return f.WriteString(string(b))
}

func (f *File) WriteString(s string) (ret int, err error) {
	f.Contents = s
	f.Seq = f.fs.next()
	return len(s), nil
}

func (f *File) Sync() (err error) {
	return nil
}

func (f *File) Read(b []byte) (n int, err error) {
	count := len(b)
	if len(f.Contents) < count {
		count = len(f.Contents)
	}
	copy(b, []byte(f.Contents)[:count])
	f.Seq = f.fs.next()

	return count, nil
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return f.Read(b)
}

func (f *File) Fd() uintptr {
	panic("Not implemented.")
}

func (f *File) Close() error {
	return nil
}

func NewFilesystem() *Filesystem {
	return &Filesystem{
		Files: make(map[string]*File),
	}
}

func (fs *Filesystem) OpenFile(name string, flag int, perm os.FileMode) (file internal.File, err error) {
	f, ok := fs.Files[name]
	if ok {
		f.Opened = true
		f.Closed = false
		return f, nil
	} else {
		return nil, errors.New("No such file.")
	}
}

func (fs *Filesystem) Add(name string) *File {
	f := &File{
		Seq: -1,
		fs:  fs,
	}
	fs.Files[name] = f
	return f
}

func (fs *Filesystem) next() int {
	fs.Seq++
	return fs.Seq
}
