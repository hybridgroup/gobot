package sysfs

import (
	"os"
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestNewI2cDevice(t *testing.T) {
	fs := NewMockFilesystem([]string{})
	SetFilesystem(fs)

	i, err := NewI2cDevice(os.DevNull, 0xff)
	gobot.Refute(t, err, nil)

	fs = NewMockFilesystem([]string{
		"/dev/i2c-1",
	})

	SetFilesystem(fs)

	i, err = NewI2cDevice("/dev/i2c-1", 0xff)
	gobot.Refute(t, err, nil)

	SetSyscall(&MockSyscall{})

	i, err = NewI2cDevice("/dev/i2c-1", 0xff)
	var _ I2cDevice = i

	gobot.Assert(t, err, nil)

	gobot.Assert(t, i.SetAddress(0xff), nil)

	buf := []byte{0x01, 0x02, 0x03}

	n, err := i.Write(buf)

	gobot.Assert(t, n, len(buf))
	gobot.Assert(t, err, nil)

	buf = make([]byte, 4)

	n, err = i.Read(buf)

	gobot.Assert(t, n, 4)
	gobot.Assert(t, err, nil)

}
