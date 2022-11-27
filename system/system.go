package system

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"gobot.io/x/gobot"
)

const systemDebug = false

// A File represents basic IO interactions with the underlying file system
type File interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (ret int, err error)
	Sync() (err error)
	Read(b []byte) (n int, err error)
	ReadAt(b []byte, off int64) (n int, err error)
	Seek(offset int64, whence int) (ret int64, err error)
	Fd() uintptr
	Close() error
}

// filesystem is a unexposed interface to allow the switch between the native file system or a mocked implementation
type filesystem interface {
	openFile(name string, flag int, perm os.FileMode) (file File, err error)
	stat(name string) (os.FileInfo, error)
	find(baseDir string, pattern string) (dirs []string, err error)
	readFile(name string) (content []byte, err error)
}

// systemCaller represents unexposed Syscall interface to allow the switch between native and mocked implementation
// Prevent unsafe call, since go 1.15, see "Pattern 4" in: https://go101.org/article/unsafe.html
type systemCaller interface {
	syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno)
}

type digitalPinAccesser interface {
	isSupported() bool
	createPin(chip string, pin int, o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner
	setFs(fs filesystem)
}

// Accesser provides access to system calls and filesystem
type Accesser struct {
	sys              systemCaller
	fs               filesystem
	digitalPinAccess digitalPinAccesser
}

// NewAccesser returns a accesser to native system call, native file system and the chosen digital pin access.
// Digital pin accesser can be empty or "sysfs", otherwise it will be automatically chosen.
func NewAccesser(digitalPinAccess ...string) *Accesser {
	s := Accesser{
		sys: &nativeSyscall{},
		fs:  &nativeFilesystem{},
	}
	a := "sysfs"
	if len(digitalPinAccess) > 0 && digitalPinAccess[0] != "" {
		a = digitalPinAccess[0]
	}
	if a != "sysfs" {
		dpa := &gpiodDigitalPinAccess{fs: s.fs}
		if dpa.isSupported() {
			s.digitalPinAccess = dpa
			if systemDebug {
				fmt.Printf("use gpiod driver for digital pins with this chips: %v\n", dpa.chips)
			}
			return &s
		}
		if systemDebug {
			fmt.Println("gpiod driver not supported, fallback to sysfs")
		}
	}
	s.digitalPinAccess = &sysfsDigitalPinAccess{fs: s.fs}
	return &s
}

// UseDigitalPinAccessWithMockFs sets the digital pin handler accesser to the chosen one. Used only for tests.
func (a *Accesser) UseDigitalPinAccessWithMockFs(digitalPinAccess string, files []string) digitalPinAccesser {
	fs := newMockFilesystem(files)
	var dph digitalPinAccesser
	switch digitalPinAccess {
	case "sysfs":
		dph = &sysfsDigitalPinAccess{fs: fs}
	case "cdev":
		dph = &gpiodDigitalPinAccess{fs: fs}
	default:
		dph = &mockDigitalPinAccess{fs: fs}
	}
	a.fs = fs
	a.digitalPinAccess = dph
	return dph
}

// UseMockSyscall sets the Syscall implementation of the accesser to the mocked one. Used only for tests.
func (a *Accesser) UseMockSyscall() *mockSyscall {
	msc := &mockSyscall{}
	a.sys = msc
	return msc
}

// UseMockFilesystem sets the filesystem implementation of the accesser to the mocked one. Used only for tests.
func (a *Accesser) UseMockFilesystem(files []string) *MockFilesystem {
	fs := newMockFilesystem(files)
	a.fs = fs
	a.digitalPinAccess.setFs(fs)
	return fs
}

// NewDigitalPin returns a new system digital pin given the pin number and an optional pin label.
// If no label is supplied a default label will prepend to the pin number.
func (a *Accesser) NewDigitalPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner {
	return a.digitalPinAccess.createPin(chip, pin, o...)
}

// IsSysfsDigitalPinAccess returns whether the used digital pin accesser is a sysfs one.
func (a *Accesser) IsSysfsDigitalPinAccess() bool {
	if _, ok := a.digitalPinAccess.(*sysfsDigitalPinAccess); ok {
		return true
	}
	return false
}

// OpenFile opens file of given name from native or the mocked file system
func (a *Accesser) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	return a.fs.openFile(name, flag, perm)
}

// Stat returns a generic FileInfo, if the file with given name exists. It uses the native or the mocked file system.
func (a *Accesser) Stat(name string) (os.FileInfo, error) {
	return a.fs.stat(name)
}

// Find finds file from native or the mocked file system
func (a *Accesser) Find(baseDir string, pattern string) ([]string, error) {
	return a.fs.find(baseDir, pattern)
}

// ReadFile reads the named file and returns the contents. A successful call returns err == nil, not err == EOF.
// Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func (a *Accesser) ReadFile(name string) ([]byte, error) {
	return a.fs.readFile(name)
}
