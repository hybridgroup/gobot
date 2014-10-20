package edison

import (
	"os"
	"syscall"
)

const I2CSlave = 0x0703
const I2CLocation = "/dev/i2c-6"

type i2cDevice struct {
	file        *os.File
	address     byte
	i2cLocation string
}

/// newI2cDevice creates a new i2c device for address
func newI2cDevice(address byte) *i2cDevice {
	return &i2cDevice{
		i2cLocation: I2CLocation,
		address:     address,
	}
}

// start initializes i2c device.
// Panics on error.
func (i *i2cDevice) start() {
	var err error
	i.file, err = os.OpenFile(i.i2cLocation, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	_, _, errCode := syscall.Syscall(
		syscall.SYS_IOCTL,
		i.file.Fd(),
		I2CSlave, uintptr(i.address),
	)
	if errCode != 0 {
		panic(err)
	}

	i.write([]byte{0})
}

// write writes data to i2c file
func (i *i2cDevice) write(data []byte) {
	i.file.Write(data)
}

// read gets data from i2c file
func (i *i2cDevice) read(len uint) []byte {
	buf := make([]byte, len)
	i.file.Read(buf)
	return buf
}
