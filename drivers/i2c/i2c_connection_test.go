//go:build !windows
// +build !windows

package i2c

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

const dev = "/dev/i2c-1"

func getSyscallFuncImpl(errorMask byte) func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err system.SyscallErrno) {
	// bit 0: error on function query
	// bit 1: error on set address
	// bit 2: error on command
	return func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err system.SyscallErrno) {
		// function query
		if (trap == system.Syscall_SYS_IOCTL) && (a2 == system.I2C_FUNCS) {
			if errorMask&0x01 == 0x01 {
				return 0, 0, 1
			}

			var funcPtr *uint64 = (*uint64)(unsafe.Pointer(a3))
			*funcPtr = system.I2C_FUNC_SMBUS_READ_BYTE | system.I2C_FUNC_SMBUS_READ_BYTE_DATA |
				system.I2C_FUNC_SMBUS_READ_WORD_DATA |
				system.I2C_FUNC_SMBUS_WRITE_BYTE | system.I2C_FUNC_SMBUS_WRITE_BYTE_DATA |
				system.I2C_FUNC_SMBUS_WRITE_WORD_DATA
		}
		// set address
		if (trap == system.Syscall_SYS_IOCTL) && (a2 == system.I2C_SLAVE) {
			if errorMask&0x02 == 0x02 {
				return 0, 0, 1
			}
		}
		// command
		if (trap == system.Syscall_SYS_IOCTL) && (a2 == system.I2C_SMBUS) {
			if errorMask&0x04 == 0x04 {
				return 0, 0, 1
			}
		}
		// Let all operations succeed
		return 0, 0, 0
	}
}

func initI2CDevice() gobot.I2cSystemDevicer {
	a := system.NewAccesser()
	a.UseMockFilesystem([]string{dev})
	msc := a.UseMockSyscall()
	msc.Impl = getSyscallFuncImpl(0x00)

	d, _ := a.NewI2cDevice(dev)
	return d
}

func initI2CDeviceAddressError() gobot.I2cSystemDevicer {
	a := system.NewAccesser()
	a.UseMockFilesystem([]string{dev})
	msc := a.UseMockSyscall()
	msc.Impl = getSyscallFuncImpl(0x02)

	d, _ := a.NewI2cDevice(dev)
	return d
}

func TestI2CAddress(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x66)
	assert.Equal(t, 0x66, c.address)
}

func TestI2CClose(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	assert.NoError(t, c.Close())
}

func TestI2CRead(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	i, _ := c.Read([]byte{})
	assert.Equal(t, 0, i)
}

func TestI2CReadAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	_, err := c.Read([]byte{})
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CWrite(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	i, _ := c.Write([]byte{0x01})
	assert.Equal(t, 1, i)
}

func TestI2CWriteAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	_, err := c.Write([]byte{0x01})
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CReadByte(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadByte()
	assert.Equal(t, uint8(0xFC), v)
}

func TestI2CReadByteAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	_, err := c.ReadByte()
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CReadByteData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadByteData(0x01)
	assert.Equal(t, uint8(0xFD), v)
}

func TestI2CReadByteDataAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	_, err := c.ReadByteData(0x01)
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CReadWordData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	v, _ := c.ReadWordData(0x01)
	assert.Equal(t, uint16(0xFFFE), v)
}

func TestI2CReadWordDataAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	_, err := c.ReadWordData(0x01)
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CWriteByte(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteByte(0x01)
	assert.NoError(t, err)
}

func TestI2CWriteByteAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	err := c.WriteByte(0x01)
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CWriteByteData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteByteData(0x01, 0x01)
	assert.NoError(t, err)
}

func TestI2CWriteByteDataAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	err := c.WriteByteData(0x01, 0x01)
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CWriteWordData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteWordData(0x01, 0x01)
	assert.NoError(t, err)
}

func TestI2CWriteWordDataAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	err := c.WriteWordData(0x01, 0x01)
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func TestI2CWriteBlockData(t *testing.T) {
	c := NewConnection(initI2CDevice(), 0x06)
	err := c.WriteBlockData(0x01, []byte{0x01, 0x02})
	assert.NoError(t, err)
}

func TestI2CWriteBlockDataAddressError(t *testing.T) {
	c := NewConnection(initI2CDeviceAddressError(), 0x06)
	err := c.WriteBlockData(0x01, []byte{0x01, 0x02})
	assert.ErrorContains(t, err, "Setting address failed with syscall.Errno operation not permitted")
}

func Test_setBit(t *testing.T) {
	var expectedVal uint8 = 129
	actualVal := setBit(1, 7)
	assert.Equal(t, actualVal, expectedVal)
}

func Test_clearBit(t *testing.T) {
	var expectedVal uint8
	actualVal := clearBit(128, 7)
	assert.Equal(t, actualVal, expectedVal)
}
