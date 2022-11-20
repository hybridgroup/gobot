package system

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

func TestDigitalPin(t *testing.T) {
	a := NewAccesser()
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	}
	fs := a.UseMockFilesystem(mockPaths)

	pin := a.NewDigitalPin(10)
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

	pin2 := a.NewDigitalPin(30)
	err = pin2.Write(1)
	gobottest.Assert(t, err.Error(), "pin has not been exported")

	data, err = pin2.Read()
	gobottest.Assert(t, err.Error(), "pin has not been exported")
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
	a := NewAccesser()
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio11/direction",
	}
	a.UseMockFilesystem(mockPaths)

	pin := a.NewDigitalPin(10)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Export()
	gobottest.Assert(t, err.Error(), " : /sys/class/gpio/gpio10/direction: No such file.")
}

func TestDigitalPinUnexportError(t *testing.T) {
	a := NewAccesser()
	mockPaths := []string{
		"/sys/class/gpio/unexport",
	}
	a.UseMockFilesystem(mockPaths)

	pin := a.NewDigitalPin(10)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Unexport()
	gobottest.Assert(t, err.Error(), " : device or resource busy")
}
