package sysfs

import (
	"os"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestNewI2cDeviceClose(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})

	SetFilesystem(fs)
	SetSyscall(&MockSyscall{})

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i.Close(), nil)
}

func TestNewI2cDevice(t *testing.T) {
	fs := NewMockFilesystem([]string{})
	SetFilesystem(fs)

	i, err := NewI2cDevice(os.DevNull)
	gobottest.Assert(t, err.Error(), " : /dev/null: No such file.")

	fs = NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err = NewI2cDevice("/dev/i2c-1")
	gobottest.Assert(t, err, nil)

	SetSyscall(&MockSyscall{})

	i, err = NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, i.SetAddress(0xff), nil)

	buf := []byte{0x01, 0x02, 0x03}

	n, err := i.Write(buf)

	gobottest.Assert(t, n, len(buf))
	gobottest.Assert(t, err, nil)

	buf = make([]byte, 4)

	n, err = i.Read(buf)

	gobottest.Assert(t, n, 3)
	gobottest.Assert(t, err, nil)
}

func TestNewI2cDeviceReadByteNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadByte()
	gobottest.Assert(t, err.Error(), "SMBus read byte not supported")
}

func TestNewI2cDeviceWriteByteNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteByte(0x01)
	gobottest.Assert(t, err.Error(), "SMBus write byte not supported")
}

func TestNewI2cDeviceReadByteDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadByteData(0x01)
	gobottest.Assert(t, err.Error(), "SMBus read byte data not supported")
}

func TestNewI2cDeviceWriteByteDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteByteData(0x01, 0x01)
	gobottest.Assert(t, err.Error(), "SMBus write byte data not supported")
}

func TestNewI2cDeviceReadWordDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadWordData(0x01)
	gobottest.Assert(t, err.Error(), "SMBus read word data not supported")
}

func TestNewI2cDeviceWriteWordDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteWordData(0x01, 0x01)
	gobottest.Assert(t, err.Error(), "SMBus write word data not supported")
}

func TestNewI2cDeviceWrite(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	buf := []byte{0x01, 0x02, 0x03}

	n, err := i.Write(buf)

	gobottest.Assert(t, n, len(buf))
	gobottest.Assert(t, err, nil)
}
