package sysfs

import (
	"os"
)

// A File represents basic IO interactions with the underlying file system
type File interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (ret int, err error)
	Sync() (err error)
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Fd() uintptr
	Close() error
}

// Filesystem opens files and returns either a native file system or user defined
type Filesystem interface {
	OpenFile(name string, flag int, perm os.FileMode) (file File, err error)
}

// NativeFilesystem represents the native file system implementation
type NativeFilesystem struct{}

// Default to the host filesystem.
var fs Filesystem = &NativeFilesystem{}

// SetFilesystem sets the filesystem implementation.
func SetFilesystem(f Filesystem) {
	fs = f
}

// OpenFile calls os.OpenFile().
func (fs *NativeFilesystem) OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return os.OpenFile(name, flag, perm)
}

// OpenFile calls either the NativeFilesystem or user defined OpenFile
func OpenFile(name string, flag int, perm os.FileMode) (file File, err error) {
	return fs.OpenFile(name, flag, perm)
}
