package system

import (
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

var _ PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm10/polarity",
	}
	fs := a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)
	gobottest.Assert(t, pin.pin, "10")

	err := pin.Unexport()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/unexport"].Contents, "10")

	err = pin.Export()
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "10")

	gobottest.Assert(t, pin.InvertPolarity(true), nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents, "inverted")
	pol, _ := pin.Polarity()
	gobottest.Assert(t, pol, "inverted")
	gobottest.Assert(t, pin.InvertPolarity(false), nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/polarity"].Contents, "normal")
	pol, _ = pin.Polarity()
	gobottest.Assert(t, pol, "normal")

	gobottest.Refute(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents, "1")
	err = pin.Enable(true)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/enable"].Contents, "1")
	err = pin.InvertPolarity(false)
	gobottest.Assert(t, err.Error(), "Cannot set PWM polarity when enabled")

	fs.Files["/sys/class/pwm/pwmchip0/pwm10/period"].Contents = "6"
	data, _ := pin.Period()
	gobottest.Assert(t, data, uint32(6))
	gobottest.Assert(t, pin.SetPeriod(100000), nil)
	data, _ = pin.Period()
	gobottest.Assert(t, data, uint32(100000))

	gobottest.Refute(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents, "1")
	err = pin.SetDutyCycle(100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm10/duty_cycle"].Contents, "100")
	data, _ = pin.DutyCycle()
	gobottest.Assert(t, data, uint32(100))
}

func TestPwmPinAlreadyExported(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)
	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	// no error indicates that the pin was already exported
	gobottest.Assert(t, pin.Export(), nil)
}

func TestPwmPinExportError(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)

	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EFAULT}
	}

	// no error indicates that the pin was already exported
	err := pin.Export()
	gobottest.Assert(t, err.Error(), "Export() failed for pin 10 with  : bad address")
}

func TestPwmPinUnxportError(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)

	pin.write = func(filesystem, string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	err := pin.Unexport()
	gobottest.Assert(t, err.Error(), "Unexport() failed for pin 10 with  : device or resource busy")
}

func TestPwmPinPeriodError(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.Period()
	gobottest.Assert(t, err.Error(), "Period() failed for pin 10 with  : device or resource busy")
}

func TestPwmPinPolarityError(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.Polarity()
	gobottest.Assert(t, err.Error(), "Polarity() failed for pin 10 with  : device or resource busy")
}

func TestPwmPinDutyCycleError(t *testing.T) {
	a := NewAccesser()
	mockedPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	}
	a.UseMockFilesystem(mockedPaths)

	pin := a.NewPWMPin("/sys/class/pwm/pwmchip0", 10)

	pin.read = func(filesystem, string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.DutyCycle()
	gobottest.Assert(t, err.Error(), "DutyCycle() failed for pin 10 with  : device or resource busy")
}
