package sysfs

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestDigitalPin(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	SetFilesystem(fs)

	pin := NewDigitalPin(10, "custom")
	gobottest.Assert(t, pin.pin, "10")
	gobottest.Assert(t, pin.label, "custom")

	pin = NewDigitalPin(10)
	gobottest.Assert(t, pin.pin, "10")
	gobottest.Assert(t, pin.label, "gpio10")
	gobottest.Assert(t, pin.value, nil)

	err := pin.Unexport()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/unexport"].Contents, "10")

	err = pin.Export()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/export"].Contents, "10")
	gobottest.Refute(t, pin.value, nil)

	err = pin.Write(1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio10/value"].Contents, "1")

	err = pin.Direction(IN)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio10/direction"].Contents, "in")

	data, _ := pin.Read()
	gobottest.Assert(t, 1, data)

	pin2 := NewDigitalPin(30, "custom")
	err = pin2.Write(1)
	gobottest.Refute(t, err, nil)

	data, err = pin2.Read()
	gobottest.Refute(t, err, nil)
	gobottest.Assert(t, data, 0)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EINVAL}
	}

	err = pin.Unexport()
	gobottest.Assert(t, err, nil)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Unexport()
	gobottest.Assert(t, err.(*os.PathError).Err, errors.New("write error"))

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err = pin.Export()
	gobottest.Assert(t, err, nil)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Export()
	gobottest.Assert(t, err.(*os.PathError).Err, errors.New("write error"))
}

func TestDigitalPinExportError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	SetFilesystem(fs)

	pin := NewDigitalPin(10, "custom")
	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Export()
	gobottest.Refute(t, err, nil)
}

func TestDigitalPinUnexportError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	SetFilesystem(fs)

	pin := NewDigitalPin(10, "custom")
	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Unexport()
	gobottest.Refute(t, err, nil)
}
