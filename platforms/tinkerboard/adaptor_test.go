package tinkerboard

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

const (
	gpio17Path  = "/sys/class/gpio/gpio17/"
	gpio160Path = "/sys/class/gpio/gpio160/"
)

const (
	pwm2Dir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2/"
	pwm2pwm0Dir       = pwm2Dir + "pwm0/"
	pwm2ExportPath    = pwm2Dir + "export"
	pwm2UnexportPath  = pwm2Dir + "unexport"
	pwm2EnablePath    = pwm2pwm0Dir + "enable"
	pwm2PeriodPath    = pwm2pwm0Dir + "period"
	pwm2DutyCyclePath = pwm2pwm0Dir + "duty_cycle"
	pwm2PolarityPath  = pwm2pwm0Dir + "polarity"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ sysfs.PWMPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func gpioFs() *sysfs.MockFilesystem {
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		gpio17Path + "value",
		gpio17Path + "direction",
		gpio160Path + "value",
		gpio160Path + "direction",
	})

	return fs
}

func pwmFs(t *testing.T) *sysfs.MockFilesystem {
	fs := sysfs.NewMockFilesystem([]string{
		pwm2ExportPath,
		pwm2UnexportPath,
		pwm2EnablePath,
		pwm2PeriodPath,
		pwm2DutyCyclePath,
		pwm2PolarityPath,
	})
	gobottest.Assert(t, writePwmPath(fs, pwm2EnablePath, "0"), nil)
	gobottest.Assert(t, writePwmPath(fs, pwm2PeriodPath, "0"), nil)
	gobottest.Assert(t, writePwmPath(fs, pwm2DutyCyclePath, "0"), nil)
	gobottest.Assert(t, writePwmPath(fs, pwm2PolarityPath, pwmInverted), nil)
	return fs
}

func writePwmPath(fs *sysfs.MockFilesystem, filePath string, value string) error {
	file, err := fs.OpenFile(filePath, 0, 0)
	if err != nil {
		return err
	}
	_, err = file.WriteString(value)
	return err
}

func i2cFs() *sysfs.MockFilesystem {
	return sysfs.NewMockFilesystem([]string{"/dev/i2c-1"})
}

func initTestTinkerboard(fs *sysfs.MockFilesystem) *Adaptor {
	a := NewAdaptor()
	sysfs.SetFilesystem(fs)
	return a
}

func TestTinkerboardName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Tinker Board"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestTinkerboardDigitalIO(t *testing.T) {
	fs := gpioFs()
	a := initTestTinkerboard(fs)
	a.Connect()

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files[gpio17Path+"value"].Contents, "1")

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("Not a valid pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestTinkerboardDigitalWriteError(t *testing.T) {
	fs := gpioFs()
	fs.WithWriteError = true
	a := initTestTinkerboard(fs)

	err := a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTinkerboardDigitalReadWriteError(t *testing.T) {
	fs := gpioFs()
	fs.WithWriteError = true
	a := initTestTinkerboard(fs)

	_, err := a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestTinkerboardInvalidPWMPin(t *testing.T) {
	fs := pwmFs(t)
	a := initTestTinkerboard(fs)

	err := a.PwmWrite("666", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("666", 120)
	gobottest.Refute(t, err, nil)

	err = a.PwmWrite("3", 42)
	gobottest.Refute(t, err, nil)

	err = a.ServoWrite("3", 120)
	gobottest.Refute(t, err, nil)
}

func TestTinkerboardPwmWrite(t *testing.T) {
	fs := pwmFs(t)
	a := initTestTinkerboard(fs)

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwm2ExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2EnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwm2PeriodPath].Contents, fmt.Sprintf("%d", pwmPeriodDefault))
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "3921568")
	gobottest.Assert(t, fs.Files[pwm2PolarityPath].Contents, "normal")

	err = a.ServoWrite("33", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "500000")

	err = a.ServoWrite("33", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestTinkerboardPwmWriteError(t *testing.T) {
	fs := pwmFs(t)
	fs.WithWriteError = true
	a := initTestTinkerboard(fs)

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestTinkerboardPwmWriteReadError(t *testing.T) {
	fs := pwmFs(t)
	fs.WithReadError = true
	a := initTestTinkerboard(fs)

	err := a.PwmWrite("33", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "read error"), true)
}

func TestTinkerboardSetPeriod(t *testing.T) {
	// arrange
	fs := pwmFs(t)
	a := initTestTinkerboard(fs)
	newPeriod := uint32(2550000)

	// act
	err := a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2ExportPath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2EnablePath].Contents, "1")
	gobottest.Assert(t, fs.Files[pwm2PeriodPath].Contents, fmt.Sprintf("%d", newPeriod))
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, "0")
	gobottest.Assert(t, fs.Files[pwm2PolarityPath].Contents, "normal")

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 1270000))
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 127000))

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files[pwm2DutyCyclePath].Contents, fmt.Sprintf("%d", 2540000))
}

func TestTinkerboardI2c(t *testing.T) {
	fs := i2cFs()
	a := initTestTinkerboard(fs)
	sysfs.SetSyscall(&sysfs.MockSyscall{})

	con, err := a.GetConnection(0xff, 1)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestTinkerboardDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestTinkerboardGetConnectionInvalidBus(t *testing.T) {
	a := NewAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestTinkerboardFinalizeErrorAfterGPIO(t *testing.T) {
	fs := gpioFs()
	a := initTestTinkerboard(fs)

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestTinkerboardFinalizeErrorAfterPWM(t *testing.T) {
	fs := pwmFs(t)
	a := initTestTinkerboard(fs)

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("33", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}
