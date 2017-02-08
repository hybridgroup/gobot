package sysfs

import (
	"os"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestNewI2cDevice(t *testing.T) {
	fs := NewMockFilesystem([]string{})
	SetFilesystem(fs)

	i, err := NewI2cDevice(os.DevNull)
	gobottest.Refute(t, err, nil)

	fs = NewMockFilesystem([]string{
		"/dev/i2c-1",
	})

	SetFilesystem(fs)

	i, err = NewI2cDevice("/dev/i2c-1")
	gobottest.Refute(t, err, nil)

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
