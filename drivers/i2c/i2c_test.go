package i2c

import (
	"testing"

	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

func initI2CDevice() sysfs.I2cDevice {
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)

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

func TestI2CWrite(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	i, _ := c.Write([]byte{0x01})
	gobottest.Assert(t, i, 1)
}

func TestI2CReadByte(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadByte()
	gobottest.Assert(t, v, uint8(0))
}

func TestI2CReadByteData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadByteData(0x01)
	gobottest.Assert(t, v, uint8(0))
}

func TestI2CReadWordData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadWordData(0x01)
	gobottest.Assert(t, v, uint16(0))
}

func TestI2CWriteByte(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteByte(0x01)
	gobottest.Assert(t, err, nil)
}

func TestI2CWriteByteData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteByteData(0x01, 0x01)
	gobottest.Assert(t, err, nil)
}

func TestI2CWriteWordData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteWordData(0x01, 0x01)
	gobottest.Assert(t, err, nil)
}

func TestI2CWriteBlockData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteBlockData(0x01, []byte{0x01, 0x02})
	gobottest.Assert(t, err, nil)
}
