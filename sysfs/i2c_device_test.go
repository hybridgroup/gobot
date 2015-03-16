package sysfs

import (
	"io"
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
	gobot.Assert(t, err, nil)
	var _ io.ReadWriteCloser = i
}
