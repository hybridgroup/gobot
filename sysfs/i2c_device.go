package sysfs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

const I2CSlave = 0x0703

// NewI2cDevice creates a new io.ReadWriteCloser with the proper ioctrl given an i2c bus location and device address
func NewI2cDevice(location string, address byte) (io.ReadWriteCloser, error) {
	file, err := OpenFile(location, os.O_RDWR, os.ModeExclusive)

	if err != nil {
		return nil, err
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		I2CSlave,
		uintptr(address),
	)

	if errno != 0 {
		return nil, errors.New(fmt.Sprintf("Failed with syscall.Errno %v", errno))
	}

	return file, nil
}
