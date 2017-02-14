package i2c

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

func initI2CDevice() sysfs.I2cDevice {
	sysfs.SetSyscall(&sysfs.MockSyscall{})
	i, _ := sysfs.NewI2cDevice("/dev/i2c-1")
	return i
}

func TestI2CAddress(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x66)
	gobottest.Assert(t, c.address, 0x66)
}

func TestI2CRead(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	i, _ := c.Read([]byte{})
	gobottest.Assert(t, i, 0)
}
