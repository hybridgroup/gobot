package gobotBeaglebone

import (
	"os"
	"syscall"
)

const I2C_SLAVE = 0x0703

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

func (me *i2cDevice) start() {
	var err error
	me.i2cDevice, err = os.OpenFile(me.i2cLocation, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		panic(err)
	}
	_, _, errCode := syscall.Syscall(syscall.SYS_IOCTL, me.i2cDevice.Fd(), I2C_SLAVE, uintptr(me.address))
	if errCode != 0 {
		panic(err)
	}

	me.write([]byte{0})
}

func (me *i2cDevice) write(data []byte) {
	me.i2cDevice.Write(data)
}

func (me *i2cDevice) read(len byte) []byte {
	buf := make([]byte, len)
	me.i2cDevice.Read(buf)
	return buf
}
