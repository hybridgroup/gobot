package arietta

import (
	"fmt"
	"github.com/hybridgroup/gobot/internal"
	"os"
	"syscall"
)

const (
	// Imported from /usr/include/linux/i2c-dev.h.
	I2C_SLAVE = 0x0703 // Change slave address.
)

// i2cDevice is a I2C device on a bus accessed through the Linux char
// and I2C bus interface.
type i2cDevice struct {
	path    string
	address byte
	file    internal.File
}

func wrapErr(method string, err error) error {
	panic(fmt.Sprintf("%s: core %v", method, err))
}

func newI2cDevice(path string, address byte) *i2cDevice {
	d := &i2cDevice{
		path:    path,
		address: address,
	}
	return d
}

func (i *i2cDevice) finalize() {
	if i.file != nil {
		i.file.Close()
		i.file = nil
	}
}

func (i *i2cDevice) start() {
	// TODO(michaelh): change to a device that belongs to a bus so
	// we can have more than one device on a bus.

	i.file = openOrDie(os.O_RDWR, i.path)

	// Set the device address.
	_, _, errCode := syscall.Syscall(syscall.SYS_IOCTL, i.file.Fd(), I2C_SLAVE, uintptr(i.address))
	if errCode != 0 {
		panic(errCode.Error())
	}
}

func (i *i2cDevice) write(data []byte) {
	// TODO(michaelh): error handling.
	i.file.Write(data)
}

func (i *i2cDevice) read(len uint) []byte {
	// TODO(michaelh): error handling.
	buf := make([]byte, len)
	i.file.Read(buf)
	return buf
}
