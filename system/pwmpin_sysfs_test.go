package system

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.PWMPinner = (*pwmPinSysFs)(nil)

const (
	normal   = "normal"
	inverted = "inversed"
)

func initTestPWMPinSysFsWithMockedFilesystem(mockPaths []string) (*pwmPinSysFs, *MockFilesystem) {
	fs := newMockFilesystem(mockPaths)
	sfa := &sysfsFileAccess{fs: fs, readBufLen: 200}
	pin := newPWMPinSysfs(sfa, "/sys/class/pwm/pwmchip0", 10, normal, inverted)
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
	require.NoError(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/pwm/pwmchip0/unexport"].Contents)

	err = pin.Export()
	require.NoError(t, err)
	assert.Equal(t, "10", fs.Files["/sys/class/pwm/pwmchip0/export"].Contents)

	require.NoError(t, pin.SetPolarity(false))
	assert.Equal(t, inverted, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents)
	pol, _ := pin.Polarity()
	assert.False(t, pol)
	require.NoError(t, pin.SetPolarity(true))
	assert.Equal(t, normal, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents)
	pol, _ = pin.Polarity()
	assert.True(t, pol)

	assert.NotEqual(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents)
	err = pin.SetEnabled(true)
	require.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents)
	err = pin.SetPolarity(true)
	require.ErrorContains(t, err, "Cannot set PWM polarity when enabled")

	fs.Files["/sys/class/pwm/pwmchip0/pwm10/period"].Contents = "6"
	data, _ := pin.Period()
	assert.Equal(t, uint32(6), data)
	require.NoError(t, pin.SetPeriod(100000))
	data, _ = pin.Period()
	assert.Equal(t, uint32(100000), data)

	assert.NotEqual(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents)
	err = pin.SetDutyCycle(100)
	require.NoError(t, err)
	assert.Equal(t, "100", fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents)
	data, _ = pin.DutyCycle()
	assert.Equal(t, uint32(100), data)
}

func TestPwmPinAlreadyExported(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/export"].simulateWriteError = &os.PathError{Err: Syscall_EBUSY}
	// act & assert: no error indicates that the pin was already exported
	require.NoError(t, pin.Export())
}

func TestPwmPinExportError(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/export"].simulateWriteError = &os.PathError{Err: Syscall_EFAULT}
	// act
	err := pin.Export()
	// assert
	require.ErrorContains(t, err, "Export() failed for id 10 with  : bad address")
}

func TestPwmPinUnxportError(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/unexport"].simulateWriteError = &os.PathError{Err: Syscall_EBUSY}
	// act
	err := pin.Unexport()
	// assert
	require.ErrorContains(t, err, "Unexport() failed for id 10 with  : device or resource busy")
}

func TestPwmPinPeriodError(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm10/period"].simulateReadError = &os.PathError{Err: Syscall_EBUSY}
	// act
	_, err := pin.Period()
	// assert
	require.ErrorContains(t, err, "Period() failed for id 10 with  : device or resource busy")
}

func TestPwmPinPolarityError(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm10/polarity",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].simulateReadError = &os.PathError{Err: Syscall_EBUSY}
	// act
	_, err := pin.Polarity()
	// assert
	require.ErrorContains(t, err, "Polarity() failed for id 10 with  : device or resource busy")
}

func TestPwmPinDutyCycleError(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	pin, fs := initTestPWMPinSysFsWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].simulateReadError = &os.PathError{Err: Syscall_EBUSY}
	// act
	_, err := pin.DutyCycle()
	// assert
	require.ErrorContains(t, err, "DutyCycle() failed for id 10 with  : device or resource busy")
}
