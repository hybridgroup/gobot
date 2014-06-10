package beaglebone

import (
	"os"
	"syscall"
)

const I2CSlave = 0x0703

type i2cDevice struct {
	i2cDevice   *os.File
	address     byte
	i2cLocation string
}

func newI2cDevice(i2cLocation string, address byte) *i2cDevice {
	d := new(i2cDevice)
	d.i2cLocation = i2cLocation
	d.address = address
	return d
}

func (i *i2cDevice) start() {
	var err error
	i.i2cDevice, err = os.OpenFile(i.i2cLocation, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	_, _, errCode := syscall.Syscall(syscall.SYS_IOCTL, i.i2cDevice.Fd(), I2CSlave, uintptr(i.address))
	if errCode != 0 {
		panic(err)
	}

	i.write([]byte{0})
}

func (i *i2cDevice) write(data []byte) {
	i.i2cDevice.Write(data)
}

func (i *i2cDevice) read(len uint) []byte {
	buf := make([]byte, len)
	i.i2cDevice.Read(buf)
	return buf
}
