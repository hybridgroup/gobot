package i2c

import (
	"testing"

	"syscall"
	"unsafe"

	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

func syscallImpl(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
	if (trap == syscall.SYS_IOCTL) && (a2 == sysfs.I2C_FUNCS) {
		var funcPtr *uint64 = (*uint64)(unsafe.Pointer(a3))
		*funcPtr = sysfs.I2C_FUNC_SMBUS_READ_BYTE | sysfs.I2C_FUNC_SMBUS_READ_BYTE_DATA |
			sysfs.I2C_FUNC_SMBUS_READ_WORD_DATA |
			sysfs.I2C_FUNC_SMBUS_WRITE_BYTE | sysfs.I2C_FUNC_SMBUS_WRITE_BYTE_DATA |
			sysfs.I2C_FUNC_SMBUS_WRITE_WORD_DATA
	}
	// Let all operations succeed
	return 0, 0, 0
}

func initI2CDevice() sysfs.I2cDevice {
	fs := sysfs.NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	sysfs.SetFilesystem(fs)

	sysfs.SetSyscall(&sysfs.MockSyscall{
		Impl: syscallImpl,
	})
	i, _ := sysfs.NewI2cDevice("/dev/i2c-1")
	return i
}

func TestI2CAddress(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x66)
	gobottest.Assert(t, c.address, 0x66)
}

func TestI2CClose(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	gobottest.Assert(t, c.Close(), nil)
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
