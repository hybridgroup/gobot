package sysfs

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
)

func TestNewI2cDeviceClose(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})

	SetFilesystem(fs)
	SetSyscall(&MockSyscall{})

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i.Close(), nil)
}

func TestNewI2cDeviceQueryFuncError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	SetSyscall(&MockSyscall{
		Impl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
			return 0, 0, 1
		},
	})

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, errors.New("Querying functionality failed with syscall.Errno operation not permitted"))
}

func TestNewI2cDevice(t *testing.T) {
	fs := NewMockFilesystem([]string{})
	SetFilesystem(fs)
	SetSyscall(&MockSyscall{})

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
	var _ i2c.I2cDevice = i

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

func TestNewI2cDeviceReadByte(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_READ_BYTE

	val, e := i.ReadByte()
	gobottest.Assert(t, val, byte(0))
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceReadByteError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	SetSyscall(&MockSyscall{
		Impl: func(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err syscall.Errno) {
			return 0, 0, 1
		},
	})

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_READ_BYTE

	_, e := i.ReadByte()
	gobottest.Refute(t, e, nil)
}

func TestNewI2cDeviceReadByteNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadByte()
	gobottest.Assert(t, err.Error(), "SMBus read byte not supported")
}

func TestNewI2cDeviceWriteByte(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_WRITE_BYTE

	e := i.WriteByte(0x01)
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceWriteByteNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteByte(0x01)
	gobottest.Assert(t, err.Error(), "SMBus write byte not supported")
}

func TestNewI2cDeviceReadByteData(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_READ_BYTE_DATA

	v, e := i.ReadByteData(0x01)
	gobottest.Assert(t, v, byte(0))
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceReadByteDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadByteData(0x01)
	gobottest.Assert(t, err.Error(), "SMBus read byte data not supported")
}

func TestNewI2cDeviceWriteByteData(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_WRITE_BYTE_DATA

	e := i.WriteByteData(0x01, 0x02)
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceWriteByteDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteByteData(0x01, 0x01)
	gobottest.Assert(t, err.Error(), "SMBus write byte data not supported")
}

func TestNewI2cDeviceReadWordData(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_READ_WORD_DATA

	v, e := i.ReadWordData(0x01)
	gobottest.Assert(t, v, uint16(0))
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceReadWordDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	_, err = i.ReadWordData(0x01)
	gobottest.Assert(t, err.Error(), "SMBus read word data not supported")
}

func TestNewI2cDeviceWriteWordData(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	i.funcs = I2C_FUNC_SMBUS_WRITE_WORD_DATA

	e := i.WriteWordData(0x01, 0x0102)
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceWriteWordDataNotSupported(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	err = i.WriteWordData(0x01, 0x01)
	gobottest.Assert(t, err.Error(), "SMBus write word data not supported")
}

func TestNewI2cDeviceWriteBlockData(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	e := i.WriteBlockData(0x01, []byte{0x01, 0x02, 0x03})
	gobottest.Assert(t, e, nil)
}

func TestNewI2cDeviceWriteBlockDataTooMuch(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/dev/i2c-1",
	})
	SetFilesystem(fs)

	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)

	var data []byte
	data = make([]byte, 33)
	e := i.WriteBlockData(0x01, data)
	gobottest.Assert(t, e, errors.New("Writing blocks larger than 32 bytes (33) not supported"))
}

func TestNewI2cDeviceWrite(t *testing.T) {
	SetSyscall(&MockSyscall{})
	i, err := NewI2cDevice("/dev/i2c-1")
	var _ i2c.I2cDevice = i

	gobottest.Assert(t, err, nil)

	i.SetAddress(0xff)
	buf := []byte{0x01, 0x02, 0x03}

	n, err := i.Write(buf)

	gobottest.Assert(t, n, len(buf))
	gobottest.Assert(t, err, nil)
}
