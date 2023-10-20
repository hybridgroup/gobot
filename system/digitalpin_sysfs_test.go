package system

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var (
	_ gobot.DigitalPinner           = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinValuer        = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinOptioner      = (*digitalPinSysfs)(nil)
	_ gobot.DigitalPinOptionApplier = (*digitalPinSysfs)(nil)
)

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

	assert.Equal(t, "10", pin.pin)
	assert.Equal(t, "gpio10", pin.label)
	assert.Nil(t, pin.valFile)

	err := pin.Unexport()
	assert.Nil(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/gpio/unexport"].Contents)

	err = pin.Export()
	assert.Nil(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/gpio/export"].Contents)
	assert.NotNil(t, pin.valFile)

	err = pin.Write(1)
	assert.Nil(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio10/value"].Contents)

	err = pin.ApplyOptions(WithPinDirectionInput())
	assert.Nil(t, err)
	assert.Equal(t, "in", fs.Files["/sys/class/gpio/gpio10/direction"].Contents)

	data, _ := pin.Read()
	assert.Equal(t, data, 1)

	pin2 := newDigitalPinSysfs(fs, "30")
	err = pin2.Write(1)
	assert.Errorf(t, err, "pin has not been exported")

	data, err = pin2.Read()
	assert.Errorf(t, err, "pin has not been exported")
	assert.Equal(t, 0, data)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EINVAL}
	}

	err = pin.Unexport()
	assert.Nil(t, err)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}

	err = pin.Unexport()
	assert.Errorf(t, err.(*os.PathError).Err, "write error")

	// assert a busy error is dropped (just means "already exported")
	cnt := 0
	writeFile = func(File, []byte) (int, error) {
		cnt++
		if cnt == 1 {
			return 0, &os.PathError{Err: Syscall_EBUSY}
		}
		return 0, nil
	}
	err = pin.Export()
	assert.Nil(t, err)

	// assert write error on export
	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: errors.New("write error")}
	}
	err = pin.Export()
	assert.Errorf(t, err.(*os.PathError).Err, "write error")
}

func TestDigitalPinExportError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio11/direction",
	}
	pin, _ := initTestDigitalPinSysFsWithMockedFilesystem(mockPaths)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EBUSY}
	}

	err := pin.Export()
	assert.Errorf(t, err, " : /sys/class/gpio/gpio10/direction: no such file")
}

func TestDigitalPinUnexportError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/unexport",
	}
	pin, _ := initTestDigitalPinSysFsWithMockedFilesystem(mockPaths)

	writeFile = func(File, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EBUSY}
	}

	err := pin.Unexport()
	assert.Errorf(t, err, " : device or resource busy")
}
