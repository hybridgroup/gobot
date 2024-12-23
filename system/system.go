package system

import (
	"fmt"
	"os"
	"unsafe"

	"gobot.io/x/gobot/v2"
)

const systemDebug = false

type digitalPinAccesserType int

const (
	digitalPinAccesserTypeGpiod digitalPinAccesserType = iota
	digitalPinAccesserTypeSysfs
)

// A File represents basic IO interactions with the underlying file system
type File interface {
	Write(b []byte) (n int, err error)
	WriteString(s string) (ret int, err error)
	Sync() error
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
// For go vet false positives, see: https://github.com/golang/go/issues/41205
type systemCaller interface {
	syscall(
		trap uintptr,
		f File,
		signal uintptr,
		payload unsafe.Pointer,
		address uint16,
	) (r1, r2 uintptr, err SyscallErrno)
}

// digitalPinAccesser represents unexposed interface to allow the switch between different implementations and
// a mocked one
type digitalPinAccesser interface {
	isSupported() bool
	createPin(chip string, pin int, o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner
	setFs(fs filesystem)
	isType(accesserType digitalPinAccesserType) bool
}

// spiAccesser represents unexposed interface to allow the switch between different implementations and a mocked one
type spiAccesser interface {
	isSupported() bool
	createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error)
}

// Accesser provides access to system calls, filesystem, implementation for digital pin and SPI
type Accesser struct {
	sys              systemCaller
	fs               filesystem
	digitalPinAccess digitalPinAccesser
	spiAccess        spiAccesser
}

// NewAccesser returns a accesser to native system call, native file system and the chosen digital pin access.
// Digital pin accesser can be empty or "sysfs", otherwise it will be automatically chosen.
func NewAccesser(options ...func(Optioner)) *Accesser {
	a := &Accesser{
		sys: &nativeSyscall{},
		fs:  &nativeFilesystem{},
	}
	a.spiAccess = &periphioSpiAccess{fs: a.fs}
	a.digitalPinAccess = &gpiodDigitalPinAccess{fs: a.fs}
	for _, option := range options {
		if option == nil {
			continue
		}
		option(a)
	}
	return a
}

// UseMockDigitalPinAccess sets the digital pin handler accesser to the chosen one. Used only for tests.
func (a *Accesser) UseMockDigitalPinAccess() *mockDigitalPinAccess {
	dpa := newMockDigitalPinAccess(a.digitalPinAccess)
	a.digitalPinAccess = dpa
	return dpa
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

// UseMockSpi sets the SPI implementation of the accesser to the mocked one. Used only for tests.
func (a *Accesser) UseMockSpi() *MockSpiAccess {
	msc := &MockSpiAccess{}
	a.spiAccess = msc
	return msc
}

// NewDigitalPin returns a new system digital pin, according to the given pin number.
func (a *Accesser) NewDigitalPin(chip string, pin int,
	o ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	return a.digitalPinAccess.createPin(chip, pin, o...)
}

// IsSysfsDigitalPinAccess returns whether the used digital pin accesser is a sysfs one.
func (a *Accesser) IsSysfsDigitalPinAccess() bool {
	return a.digitalPinAccess.isType(digitalPinAccesserTypeSysfs)
}

// IsGpiodDigitalPinAccess returns whether the used digital pin accesser is a sysfs one.
func (a *Accesser) IsGpiodDigitalPinAccess() bool {
	return a.digitalPinAccess.isType(digitalPinAccesserTypeGpiod)
}

// NewPWMPin returns a new system PWM pin, according to the given pin number.
func (a *Accesser) NewPWMPin(path string, pin int, polNormIdent string, polInvIdent string) gobot.PWMPinner {
	sfa := &sysfsFileAccess{fs: a.fs, readBufLen: 200}
	return newPWMPinSysfs(sfa, path, pin, polNormIdent, polInvIdent)
}

func (a *Accesser) NewAnalogPin(path string, w bool, readBufLen uint16) gobot.AnalogPinner {
	r := readBufLen > 0
	if readBufLen == 0 {
		readBufLen = 32 // max. count of characters for int value is 20
	}
	return newAnalogPinSysfs(&sysfsFileAccess{fs: a.fs, readBufLen: readBufLen}, path, r, w)
}

// NewSpiDevice returns a new connection to SPI with the given parameters.
func (a *Accesser) NewSpiDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error) {
	return a.spiAccess.createDevice(busNum, chipNum, mode, bits, maxSpeed)
}

// NewOneWireDevice returns a new 1-wire device with the given parameters.
// note: this is a basic implementation without using the possibilities of bus controller
// it depends on automatic device search, see https://www.kernel.org/doc/Documentation/w1/w1.generic
func (a *Accesser) NewOneWireDevice(familyCode byte, serialNumber uint64) (gobot.OneWireSystemDevicer, error) {
	sfa := &sysfsFileAccess{fs: a.fs, readBufLen: 200}
	deviceID := fmt.Sprintf("%02x-%012x", familyCode, serialNumber)
	return newOneWireDeviceSysfs(sfa, deviceID), nil
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
