package sysfs

import (
	"syscall"
	"unsafe"
)

// SystemCaller represents a Syscall
// Prevent unsafe call, since go 1.15, see "Pattern 4" in: https://go101.org/article/unsafe.html
type SystemCaller interface {
	Syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno)
}

// NativeSyscall represents the native Syscall
type NativeSyscall struct{}

// MockSyscall represents the mock Syscall used for unit tests
type MockSyscall struct {
	lastTrap   uintptr
	lastFile   File
	lastSignal uintptr
	devAddress uintptr
	smbus      *i2cSmbusIoctlData
	dataSlice  []byte
	Impl       func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno)
}

var sys SystemCaller = &NativeSyscall{}

// SetSyscall sets the Syscall implementation
func SetSyscall(s SystemCaller) {
	sys = s
}

// Syscall calls either the NativeSyscall or user defined Syscall
func Syscall(trap uintptr, f File, a2 uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno) {
	return sys.Syscall(trap, f, a2, payload)
}

// Syscall calls the native syscall.Syscall, implements the SystemCaller interface
func (sys *NativeSyscall) Syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno) {
	return syscall.Syscall(trap, f.Fd(), signal, uintptr(payload))
}

// Syscall calls the user defined implementation, used for tests, implements the SystemCaller interface
func (sys *MockSyscall) Syscall(trap uintptr, f File, signal uintptr, payload unsafe.Pointer) (r1, r2 uintptr, err syscall.Errno) {
	sys.lastTrap = trap     // points to the used syscall (e.g. "SYS_IOCTL")
	sys.lastFile = f        // a character device file (e.g. file to path "/dev/i2c-1")
	sys.lastSignal = signal // points to used function type (e.g. I2C_SMBUS, I2C_RDWR)

	if signal == I2C_SLAVE {
		// in this case the uintptr corresponds the address
		sys.devAddress = uintptr(payload)
	}

	if signal == I2C_SMBUS {
		// set the I2C smbus data object reference to payload and fill with some data
		sys.smbus = (*i2cSmbusIoctlData)(payload)

		// get the data object payload as byte slice
		if sys.smbus.readWrite == I2C_SMBUS_WRITE {
			if sys.smbus.data != nil {
				sys.dataSlice = unsafe.Slice((*byte)(unsafe.Pointer(sys.smbus.data)), sys.smbus.size-1)
			}
		}

		// fill data object with data from given slice to simulate reading
		if sys.smbus.readWrite == I2C_SMBUS_READ {
			if sys.dataSlice != nil {
				dataSize := sys.smbus.size - 1
				if sys.smbus.size == I2C_SMBUS_BYTE {
					dataSize = 1
				}
				slc := unsafe.Slice((*byte)(unsafe.Pointer(sys.smbus.data)), dataSize)
				copy(slc, sys.dataSlice)
			}
		}
	}

	// call mock implementation
	if sys.Impl != nil {
		return sys.Impl(trap, f.Fd(), signal, uintptr(unsafe.Pointer(payload)))
	}
	return 0, 0, 0
}
