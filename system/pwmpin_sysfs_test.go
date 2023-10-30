package system

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.PWMPinner = (*pwmPinSysFs)(nil)

const (
	normal   = "normal"
	inverted = "inverted"
)

func initTestPWMPinSysFsWithMockedFilesystem(mockPaths []string) (*pwmPinSysFs, *MockFilesystem) {
	fs := newMockFilesystem(mockPaths)
	pin := newPWMPinSysfs(fs, "/sys/class/pwm/pwmchip0", 10, normal, inverted)
	return pin, fs
}

func TestPwmPin(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm10/polarity",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	assert.Equal(t, "10", pin.pin)

	err := pin.Unexport()
	assert.NoError(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/pwm/pwmchip0/unexport"].Contents)

	err = pin.Export()
	assert.NoError(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/pwm/pwmchip0/export"].Contents)

	assert.NoError(t, pin.SetPolarity(false))
	assert.Equal(t, inverted, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents)
	pol, _ := pin.Polarity()
	assert.False(t, pol)
	assert.NoError(t, pin.SetPolarity(true))
	assert.Equal(t, normal, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents)
	pol, _ = pin.Polarity()
	assert.True(t, pol)

	assert.NotEqual(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents)
	err = pin.SetEnabled(true)
	assert.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents)
	err = pin.SetPolarity(true)
	assert.ErrorContains(t, err, "Cannot set PWM polarity when enabled")

	fs.Files["/sys/class/pwm/pwmchip0/pwm10/period"].Contents = "6"
	data, _ := pin.Period()
	assert.Equal(t, uint32(6), data)
	assert.NoError(t, pin.SetPeriod(100000))
	data, _ = pin.Period()
	assert.Equal(t, uint32(100000), data)

	assert.NotEqual(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents)
	err = pin.SetDutyCycle(100)
	assert.NoError(t, err)
	assert.Equal(t, "100", fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents)
	data, _ = pin.DutyCycle()
	assert.Equal(t, uint32(100), data)
}

func TestPwmPinAlreadyExported(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EBUSY}
	}

	// no error indicates that the pin was already exported
	assert.NoError(t, pin.Export())
}

func TestPwmPinExportError(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EFAULT}
	}

	// no error indicates that the pin was already exported
	err := pin.Export()
	assert.ErrorContains(t, err, "Export() failed for id 10 with  : bad address")
}

func TestPwmPinUnxportError(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: Syscall_EBUSY}
	}

	err := pin.Unexport()
	assert.ErrorContains(t, err, "Unexport() failed for id 10 with  : device or resource busy")
}

func TestPwmPinPeriodError(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: Syscall_EBUSY}
	}

	_, err := pin.Period()
	assert.ErrorContains(t, err, "Period() failed for id 10 with  : device or resource busy")
}

func TestPwmPinPolarityError(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: Syscall_EBUSY}
	}

	_, err := pin.Polarity()
	assert.ErrorContains(t, err, "Polarity() failed for id 10 with  : device or resource busy")
}

func TestPwmPinDutyCycleError(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, _ := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: Syscall_EBUSY}
	}

	_, err := pin.DutyCycle()
	assert.ErrorContains(t, err, "DutyCycle() failed for id 10 with  : device or resource busy")
}
