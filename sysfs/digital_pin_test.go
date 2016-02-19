package sysfs

import (
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/hybridgroup/gobot"
)

func TestDigitalPin(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	})

	SetFilesystem(fs)

	pin := NewDigitalPin(10, "custom").(*digitalPin)
	gobot.Assert(t, pin.pin, "10")
	gobot.Assert(t, pin.label, "custom")

	pin = NewDigitalPin(10).(*digitalPin)
	gobot.Assert(t, pin.pin, "10")
	gobot.Assert(t, pin.label, "gpio10")
	gobot.Assert(t, pin.value, nil)

	err := pin.Unexport()
	gobot.Assert(t, err, nil)
	gobot.Assert(t, fs.Files["/sys/class/gpio/unexport"].Contents, "10")

	err = pin.Export()
	gobot.Assert(t, err, nil)
	gobot.Assert(t, fs.Files["/sys/class/gpio/export"].Contents, "10")
	gobot.Refute(t, pin.value, nil)

	err = pin.Write(1)
	gobot.Assert(t, err, nil)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio10/value"].Contents, "1")

	err = pin.Direction(IN)
	gobot.Assert(t, err, nil)
	gobot.Assert(t, fs.Files["/sys/class/gpio/gpio10/direction"].Contents, "in")

	data, _ := pin.Read()
	gobot.Assert(t, 1, data)

	pin2 := NewDigitalPin(30, "custom")
	err = pin2.Write(1)
	gobot.Refute(t, err, nil)

	data, err = pin2.Read()
	gobot.Refute(t, err, nil)
	gobot.Assert(t, data, 0)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EINVAL}
	}

	err = pin.Unexport()
	gobot.Assert(t, err, nil)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Unexport()
	gobot.Assert(t, err.(*os.PathError).Err, errors.New("write error"))

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err = pin.Export()
	gobot.Assert(t, err, nil)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Export()
	gobot.Assert(t, err.(*os.PathError).Err, errors.New("write error"))
}
