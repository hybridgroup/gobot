package system

import (
	"fmt"
	"os"
	"unsafe"

	"gobot.io/x/gobot/v2"
)

type digitalPinAccesserType int

const (
	digitalPinAccesserTypeCdev digitalPinAccesserType = iota
	digitalPinAccesserTypeSysfs
)

type spiBusAccesserType int

const (
	spiBusAccesserTypePeriphio spiBusAccesserType = iota
	spiBusAccesserTypeGPIO
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
	isType(accesserType digitalPinAccesserType) bool
	isSupported() bool
	createPin(chip string, pin int, o ...func(gobot.DigitalPinOptioner) bool) gobot.DigitalPinner
	setFs(fs filesystem)
}

// spiAccesser represents unexposed interface to allow the switch between different implementations and a mocked one
type spiAccesser interface {
	isType(accesserType spiBusAccesserType) bool
	isSupported() bool
	createDevice(busNum, chipNum, mode, bits int, maxSpeed int64) (gobot.SpiSystemDevicer, error)
}

type accesserConfiguration struct {
	debug           bool
	debugSpi        bool
	debugDigitalPin bool
	useGpioSysfs    *bool
	spiGpioConfig   *spiGpioConfig
}

// Accesser provides access to system calls, filesystem, implementation for digital pin and SPI
type Accesser struct {
	accesserCfg      *accesserConfiguration
	sys              systemCaller
	fs               filesystem
	digitalPinAccess digitalPinAccesser
	spiAccess        spiAccesser
}

// NewAccesser returns a accesser to native system call, native file system and the chosen digital pin access.
// Digital pin accesser can be empty or "sysfs", otherwise it will be automatically chosen.
func NewAccesser(options ...AccesserOptionApplier) *Accesser {
	a := &Accesser{
		accesserCfg: &accesserConfiguration{},
	}

	for _, o := range options {
		if o == nil {
			continue
		}
		o.apply(a.accesserCfg)
	}
	return a
}

// AddAnalogSupport adds the support to access the analog features of the system, usually by sysfs.
func (a *Accesser) AddAnalogSupport() {
	if a.fs == nil {
		a.fs = &nativeFilesystem{} // for sysfs access
	}
}

// AddPWMSupport adds the support to access the PWM features of the system, usually by sysfs.
func (a *Accesser) AddPWMSupport() {
	if a.fs == nil {
		a.fs = &nativeFilesystem{} // for sysfs access
	}
}

// AddDigitalPinSupport adds the support to access the GPIO features of the system. Usually by character device or
// sysfs Kernel API. Related options can be applied here.
func (a *Accesser) AddDigitalPinSupport(options ...AccesserOptionApplier) {
	for _, o := range options {
		if o == nil {
			continue
		}
		o.apply(a.accesserCfg)
	}

	if a.fs == nil {
		a.fs = &nativeFilesystem{} // for sysfs access or check for /dev/gpiochip* in cdev
	}

	if a.accesserCfg.useGpioSysfs == nil || !*a.accesserCfg.useGpioSysfs {
		dpa := &cdevDigitalPinAccess{fs: a.fs}

		if dpa.isSupported() || a.accesserCfg.useGpioSysfs == nil {
			a.digitalPinAccess = dpa

			if a.accesserCfg.debug || a.accesserCfg.debugDigitalPin {
				fmt.Printf("use cdev driver for digital pins with this chips: %v\n", dpa.chips)
			}

			return
		}

		if a.accesserCfg.debug || a.accesserCfg.debugDigitalPin {
			fmt.Println("cdev driver not supported, fallback to sysfs driver")
		}
	}

	// currently sysfs is supported by all Kernels
	dpa := &sysfsDigitalPinAccess{sfa: &sysfsFileAccess{fs: a.fs, readBufLen: 2}}
	a.digitalPinAccess = dpa
	if a.accesserCfg.debug || a.accesserCfg.debugDigitalPin {
		fmt.Println("use sysfs driver for digital pins")
	}
}

// HasDigitalPinSysfsAccess returns whether the used digital pin accesser is a sysfs one.
// If no digital pin accesser is defined, returns false.
func (a *Accesser) HasDigitalPinSysfsAccess() bool {
	return a.digitalPinAccess != nil && a.digitalPinAccess.isType(digitalPinAccesserTypeSysfs)
}

// HasDigitalPinCdevAccess returns whether the used digital pin accesser is a sysfs one.
// If no digital pin accesser is defined, returns false.
func (a *Accesser) HasDigitalPinCdevAccess() bool {
	return a.digitalPinAccess != nil && a.digitalPinAccess.isType(digitalPinAccesserTypeCdev)
}

// AddI2CSupport adds the support to access the I2C features of the system, usually by syscall with character device.
func (a *Accesser) AddI2CSupport() {
	if a.fs == nil {
		a.fs = &nativeFilesystem{} // for access to the i2c character device, e.g. /dev/i2c-2
	}

	a.sys = &nativeSyscall{}
}

// AddSPISupport adds the support to access the SPI features of the system, usually by character device or GPIOs.
// Related options can be applied here.
func (a *Accesser) AddSPISupport(options ...AccesserOptionApplier) {
	for _, o := range options {
		if o == nil {
			continue
		}
		o.apply(a.accesserCfg)
	}

	if a.fs == nil {
		a.fs = &nativeFilesystem{} // to check for "/dev/spidev*" or access by GPIO (see AddDigitalPinSupport())
	}

	if a.accesserCfg.spiGpioConfig != nil {
		// currently GPIO SPI access is always supported
		a.accesserCfg.spiGpioConfig.debug = a.accesserCfg.debugSpi
		a.spiAccess = &gpioSpiAccess{cfg: *a.accesserCfg.spiGpioConfig}

		if a.accesserCfg.debug || a.accesserCfg.debugSpi {
			fmt.Printf("use gpio driver for SPI with this config: %s\n", a.accesserCfg.spiGpioConfig.String())
		}

		return
	}

	gsa := &periphioSpiAccess{fs: a.fs}
	if !gsa.isSupported() {
		if a.accesserCfg.debug || a.accesserCfg.debugSpi {
			fmt.Println("periphio driver not supported for SPI, please activate SPI or try to use GPIOs")
		}
		return
	}

	a.spiAccess = gsa
	if a.accesserCfg.debug || a.accesserCfg.debugSpi {
		fmt.Println("use periphio driver for SPI")
	}
}

// HasSpiPeriphioAccess returns whether the used SPI accesser is periphio based.
// If SPI accesser is defined, returns false.
func (a *Accesser) HasSpiPeriphioAccess() bool {
	return a.spiAccess != nil && a.spiAccess.isType(spiBusAccesserTypePeriphio)
}

// HasSpiGpioAccess returns whether the used SPI accesser is GPIO based.
// If SPI accesser is defined, returns false.
func (a *Accesser) HasSpiGpioAccess() bool {
	return a.spiAccess != nil && a.spiAccess.isType(spiBusAccesserTypeGPIO)
}

// AddOneWireSupport adds the support to access the one wire features of the system, usually by sysfs.
func (a *Accesser) AddOneWireSupport() {
	if a.fs == nil {
		a.fs = &nativeFilesystem{} // for sysfs access
	}
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
	if a.digitalPinAccess != nil {
		a.digitalPinAccess.setFs(fs)
	}

	return fs
}

// UseMockSpi sets the SPI implementation of the accesser to the mocked one. Used only for tests.
func (a *Accesser) UseMockSpi() *MockSpiAccess {
	msc := newMockSpiAccess(a.spiAccess)
	a.spiAccess = msc
	return msc
}

// NewDigitalPin returns a new system digital pin, according to the given pin number.
func (a *Accesser) NewDigitalPin(chip string, pin int,
	options ...func(gobot.DigitalPinOptioner) bool,
) gobot.DigitalPinner {
	return a.digitalPinAccess.createPin(chip, pin, options...)
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
