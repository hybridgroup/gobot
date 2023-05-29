package system

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/gobottest"
)

var _ gobot.DigitalPinner = (*digitalPinSysfs)(nil)
var _ gobot.DigitalPinValuer = (*digitalPinSysfs)(nil)
var _ gobot.DigitalPinOptioner = (*digitalPinSysfs)(nil)
var _ gobot.DigitalPinOptionApplier = (*digitalPinSysfs)(nil)

func initTestDigitalPinSysFsWithMockedFilesystem(mockPaths []string) (*digitalPinSysfs, *MockFilesystem) {
	fs := newMockFilesystem(mockPaths)
	pin := newDigitalPinSysfs(fs, "10")
	return pin, fs
}

func TestDigitalPin(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	}
	pin, fs := initTestDigitalPinSysFsWithMockedFilesystem(mockPaths)

	gobottest.Assert(t, pin.pin, "10")
	gobottest.Assert(t, pin.label, "gpio10")
	gobottest.Assert(t, pin.valFile, nil)

	err := pin.Unexport()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/unexport"].Contents, "10")

	err = pin.Export()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/export"].Contents, "10")
	gobottest.Refute(t, pin.valFile, nil)

	err = pin.Write(1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio10/value"].Contents, "1")

	err = pin.ApplyOptions(WithPinDirectionInput())
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio10/direction"].Contents, "in")

	data, _ := pin.Read()
	gobottest.Assert(t, 1, data)

	pin2 := newDigitalPinSysfs(fs, "30")
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

	// assert a busy error is dropped (just means "already exported")
	cnt := 0
	writeFile = func(File, []byte) (int, error) {
		cnt++
		if cnt == 1 {
			return 0, &os.PathError{Err: syscall.EBUSY}
		}
		return 0, nil
	}
	err = pin.Export()
	gobottest.Assert(t, err, nil)

	// assert write error on export
	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}
	err = pin.Export()
	gobottest.Assert(t, err.(*os.PathError).Err, errors.New("write error"))
}

func TestDigitalPinExportError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio11/direction",
	}
	pin, _ := initTestDigitalPinSysFsWithMockedFilesystem(mockPaths)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Export()
	gobottest.Assert(t, err.Error(), " : /sys/class/gpio/gpio10/direction: no such file")
}

func TestDigitalPinUnexportError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/unexport",
	}
	pin, _ := initTestDigitalPinSysFsWithMockedFilesystem(mockPaths)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Unexport()
	gobottest.Assert(t, err.Error(), " : device or resource busy")
}
