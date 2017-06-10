package sysfs

import (
	"os"
	"syscall"
	"testing"

	"gobot.io/x/gobot/gobottest"
)

var _ PWMPinner = (*PWMPin)(nil)

func TestPwmPin(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm10/polarity",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
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
	gobottest.Refute(t, pin.InvertPolarity(false), nil)

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
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.write = func(string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	// no error indicates that the pin was already exported
	gobottest.Assert(t, pin.Export(), nil)
}

func TestPwmPinExportError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.write = func(string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EFAULT}
	}

	// no error indicates that the pin was already exported
	gobottest.Refute(t, pin.Export(), nil)
}

func TestPwmPinUnxportError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.write = func(string, []byte) (int, error) {
		return 0, &os.PathError{Err: syscall.EBUSY}
	}

	gobottest.Refute(t, pin.Unexport(), nil)
}

func TestPwmPinPeriodError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.read = func(string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.Period()
	gobottest.Refute(t, err, nil)
}

func TestPwmPinPolarityError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.read = func(string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.Polarity()
	gobottest.Refute(t, err, nil)
}

func TestPwmPinDutyCycleError(t *testing.T) {
	fs := NewMockFilesystem([]string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm10/enable",
		"/sys/class/pwm/pwmchip0/pwm10/period",
		"/sys/class/pwm/pwmchip0/pwm10/duty_cycle",
	})

	SetFilesystem(fs)

	pin := NewPWMPin(10)
	pin.read = func(string) ([]byte, error) {
		return nil, &os.PathError{Err: syscall.EBUSY}
	}

	_, err := pin.DutyCycle()
	gobottest.Refute(t, err, nil)
}
