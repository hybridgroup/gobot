package sysfs

import (
	"os"
)

type File interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (ret int, err error)
	Sync() (err error)
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Fd() uintptr
	Close() error
}

type Filesystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (file File, err error)
}

// Filesystem that opens real files on the host.
type NativeFilesystem struct{}

// Default to the host filesystem.
var fs Filesystem = &NativeFilesystem{}

// Override the default filesystem.
func SetFilesystem(f Filesystem) {
	fs = f
}

func (fs *NativeFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return os.OpenFile(name, flag, perm)
}

// Open a file the same as os.OpenFile().
func OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return fs.OpenFile(name, flag, perm)
}
