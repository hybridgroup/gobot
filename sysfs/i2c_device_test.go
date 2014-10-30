package sysfs

import (
	"io"
	"os"
	"testing"
)

func TestNewI2cDevice(t *testing.T) {
	i, _ := NewI2cDevice(os.DevNull, 0xff)
	var _ io.ReadWriteCloser = i
}
