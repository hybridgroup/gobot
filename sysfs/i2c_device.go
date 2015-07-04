package sysfs

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	I2C_SLAVE                = 0x0703
	I2C_SMBUS                = 0x0720
	I2C_SMBUS_WRITE          = 0
	I2C_SMBUS_READ           = 1
	I2C_SMBUS_I2C_BLOCK_DATA = 8
)

type i2cSmbusIoctlData struct {
	readWrite byte
	command   byte
	size      uint32
	data      uintptr
}

type I2cDevice interface {
	io.ReadWriteCloser
	SetAddress(int) error
}

type i2cDevice struct {
	file File
}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location and device address
func NewI2cDevice(location string, address int) (d *i2cDevice, err error) {
	d = &i2cDevice{}

	if d.file, err = OpenFile(location, os.O_RDWR, os.ModeExclusive); err != nil {
		return
	}

	err = d.SetAddress(address)

	return
}

func (d *i2cDevice) SetAddress(address int) (err error) {
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SLAVE,
		uintptr(byte(address)),
	)

	if errno != 0 {
		err = fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return
}

func (d *i2cDevice) Close() (err error) {
	return d.file.Close()
}

func (d *i2cDevice) Read(b []byte) (n int, err error) {
	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_READ,
		command:   0,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return n, fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) Write(b []byte) (n int, err error) {
	if len(b) <= 2 {
		return d.file.Write(b)
	}

	command := byte(b[0])
	buf := b[1:]

	data := make([]byte, len(buf)+1)
	data[0] = byte(len(buf))

	copy(data[1:], buf)

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_WRITE,
		command:   command,
		size:      I2C_SMBUS_I2C_BLOCK_DATA,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		err = fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return len(b), err
}
