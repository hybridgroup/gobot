package edison

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fullfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ sysfs.PWMPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

var testPinFiles = []string{
	"/sys/bus/iio/devices/iio:device1/in_voltage0_raw",
	"/sys/kernel/debug/gpio_debug/gpio111/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio115/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio114/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio109/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio131/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio129/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio28/current_pinmux",
	"/sys/kernel/debug/gpio_debug/gpio27/current_pinmux",
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm1/period",
	"/sys/class/pwm/pwmchip0/pwm1/enable",
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio13/value",
	"/sys/class/gpio/gpio13/direction",
	"/sys/class/gpio/gpio40/value",
	"/sys/class/gpio/gpio40/direction",
	"/sys/class/gpio/gpio128/value",
	"/sys/class/gpio/gpio128/direction",
	"/sys/class/gpio/gpio221/value",
	"/sys/class/gpio/gpio221/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio229/value",
	"/sys/class/gpio/gpio229/direction",
	"/sys/class/gpio/gpio253/value",
	"/sys/class/gpio/gpio253/direction",
	"/sys/class/gpio/gpio261/value",
	"/sys/class/gpio/gpio261/direction",
	"/sys/class/gpio/gpio214/value",
	"/sys/class/gpio/gpio214/direction",
	"/sys/class/gpio/gpio14/direction",
	"/sys/class/gpio/gpio14/value",
	"/sys/class/gpio/gpio165/direction",
	"/sys/class/gpio/gpio165/value",
	"/sys/class/gpio/gpio212/direction",
	"/sys/class/gpio/gpio212/value",
	"/sys/class/gpio/gpio213/direction",
	"/sys/class/gpio/gpio213/value",
	"/sys/class/gpio/gpio236/direction",
	"/sys/class/gpio/gpio236/value",
	"/sys/class/gpio/gpio237/direction",
	"/sys/class/gpio/gpio237/value",
	"/sys/class/gpio/gpio204/direction",
	"/sys/class/gpio/gpio204/value",
	"/sys/class/gpio/gpio205/direction",
	"/sys/class/gpio/gpio205/value",
	"/sys/class/gpio/gpio263/direction",
	"/sys/class/gpio/gpio263/value",
	"/sys/class/gpio/gpio262/direction",
	"/sys/class/gpio/gpio262/value",
	"/sys/class/gpio/gpio240/direction",
	"/sys/class/gpio/gpio240/value",
	"/sys/class/gpio/gpio241/direction",
	"/sys/class/gpio/gpio241/value",
	"/sys/class/gpio/gpio242/direction",
	"/sys/class/gpio/gpio242/value",
	"/sys/class/gpio/gpio218/direction",
	"/sys/class/gpio/gpio218/value",
	"/sys/class/gpio/gpio250/direction",
	"/sys/class/gpio/gpio250/value",
	"/dev/i2c-6",
}

func initTestAdaptor() (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem(testPinFiles)
	sysfs.SetFilesystem(fs)
	fs.Files["/sys/class/pwm/pwmchip0/pwm1/period"].Contents = "5000"
	a.Connect()
	return a, fs
}

func TestEdisonAdaptorName(t *testing.T) {
	a, _ := initTestAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Edison"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestAdaptorConnect(t *testing.T) {
	a, _ := initTestAdaptor()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 6)
	gobottest.Assert(t, a.Board(), "arduino")

	gobottest.Assert(t, a.Connect(), nil)
}

func TestAdaptorArduinoSetupFail263(t *testing.T) {
	a, fs := initTestAdaptor()
	delete(fs.Files, "/sys/class/gpio/gpio263/direction")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio263/direction: No such file"), true)
}

func TestAdaptorArduinoSetupFail240(t *testing.T) {
	a, fs := initTestAdaptor()
	delete(fs.Files, "/sys/class/gpio/gpio240/direction")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio240/direction: No such file"), true)
}

func TestAdaptorArduinoSetupFail111(t *testing.T) {
	a, fs := initTestAdaptor()
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux: No such file"), true)
}

func TestAdaptorArduinoSetupFail131(t *testing.T) {
	a, fs := initTestAdaptor()
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux: No such file"), true)
}

func TestAdaptorArduinoI2CSetupFailTristate(t *testing.T) {
	a, fs := initTestAdaptor()

	gobottest.Assert(t, a.arduinoSetup(), nil)

	fs.WithWriteError = true
	err := a.arduinoI2CSetup()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorArduinoI2CSetupFail14(t *testing.T) {
	a, fs := initTestAdaptor()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/gpio14/direction")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio14/direction: No such file"), true)
}

func TestAdaptorArduinoI2CSetupUnexportFail(t *testing.T) {
	a, fs := initTestAdaptor()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/unexport")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport: No such file"), true)
}

func TestAdaptorArduinoI2CSetupFail236(t *testing.T) {
	a, fs := initTestAdaptor()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/gpio236/direction")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio236/direction: No such file"), true)
}

func TestAdaptorArduinoI2CSetupFail28(t *testing.T) {
	a, fs := initTestAdaptor()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux: No such file"), true)
}

func TestAdaptorConnectArduinoError(t *testing.T) {
	a, _ := initTestAdaptor()
	a.writeFile = func(string, []byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestAdaptorConnectArduinoWriteError(t *testing.T) {
	a, fs := initTestAdaptor()
	fs.WithWriteError = true
	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestAdaptorConnectSparkfun(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("sparkfun")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
	gobottest.Assert(t, a.Board(), "sparkfun")
}

func TestAdaptorConnectMiniboard(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("miniboard")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
	gobottest.Assert(t, a.Board(), "miniboard")
}

func TestAdaptorConnectUnknown(t *testing.T) {
	a, _ := initTestAdaptor()
	a.SetBoard("wha")
	gobottest.Refute(t, a.Connect(), nil)
}

func TestAdaptorFinalize(t *testing.T) {
	a, _ := initTestAdaptor()
	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	a.GetConnection(0xff, 6)

	gobottest.Assert(t, a.Finalize(), nil)

	sysfs.SetFilesystem(sysfs.NewMockFilesystem([]string{}))
	gobottest.Refute(t, a.Finalize(), nil)
}

func TestAdaptorFinalizeError(t *testing.T) {
	a, fs := initTestAdaptor()
	a.PwmWrite("5", 100)

	fs.WithWriteError = true
	gobottest.Refute(t, a.Finalize(), nil)
}

func TestAdaptorDigitalIO(t *testing.T) {
	a, fs := initTestAdaptor()

	a.DigitalWrite("13", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio40/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestAdaptorDigitalPinInFileError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		// "/sys/class/gpio/gpio40/value",
		// "/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio229/value", // resistor
		"/sys/class/gpio/gpio229/direction",
		"/sys/class/gpio/gpio243/value",
		"/sys/class/gpio/gpio243/direction",
		"/sys/class/gpio/gpio261/value", // level shifter
		"/sys/class/gpio/gpio261/direction",
	})
	sysfs.SetFilesystem(fs)

	a.Connect()

	_, err := a.DigitalPin("13", "in")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)

}

func TestAdaptorDigitalPinInResistorFileError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		// "/sys/class/gpio/gpio229/value", // resistor
		// "/sys/class/gpio/gpio229/direction",
		"/sys/class/gpio/gpio243/value",
		"/sys/class/gpio/gpio243/direction",
		"/sys/class/gpio/gpio261/value", // level shifter
		"/sys/class/gpio/gpio261/direction",
	})
	sysfs.SetFilesystem(fs)

	a.Connect()

	_, err := a.DigitalPin("13", "in")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestAdaptorDigitalPinInLevelShifterFileError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio229/value", // resistor
		"/sys/class/gpio/gpio229/direction",
		"/sys/class/gpio/gpio243/value",
		"/sys/class/gpio/gpio243/direction",
		// "/sys/class/gpio/gpio261/value", // level shifter
		// "/sys/class/gpio/gpio261/direction",
	})
	sysfs.SetFilesystem(fs)

	a.Connect()

	_, err := a.DigitalPin("13", "in")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestAdaptorDigitalPinInMuxFileError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio229/value", // resistor
		"/sys/class/gpio/gpio229/direction",
		// "/sys/class/gpio/gpio243/value",
		// "/sys/class/gpio/gpio243/direction",
		"/sys/class/gpio/gpio261/value", // level shifter
		"/sys/class/gpio/gpio261/direction",
	})
	sysfs.SetFilesystem(fs)

	a.Connect()

	_, err := a.DigitalPin("13", "in")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestAdaptorDigitalWriteError(t *testing.T) {
	a, fs := initTestAdaptor()
	fs.WithWriteError = true

	err := a.DigitalWrite("13", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorDigitalReadWriteError(t *testing.T) {
	a, fs := initTestAdaptor()
	fs.WithWriteError = true

	_, err := a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorI2c(t *testing.T) {
	a, _ := initTestAdaptor()

	sysfs.SetSyscall(&sysfs.MockSyscall{})
	con, err := a.GetConnection(0xff, 6)
	gobottest.Assert(t, err, nil)

	con.Write([]byte{0x00, 0x01})
	data := []byte{42, 42}
	con.Read(data)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAdaptorI2cInvalidBus(t *testing.T) {
	a, _ := initTestAdaptor()
	_, err := a.GetConnection(0xff, 3)
	gobottest.Assert(t, err, errors.New("Unsupported I2C bus"))
}

func TestAdaptorPwm(t *testing.T) {
	a, fs := initTestAdaptor()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("7", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}

func TestAdaptorPwmExportError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio13/direction",
		"/sys/class/gpio/gpio13/value",
		"/sys/class/gpio/gpio221/direction",
		"/sys/class/gpio/gpio221/value",
		"/sys/class/gpio/gpio253/direction",
		"/sys/class/gpio/gpio253/value",
		//"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm1/period",
		"/sys/class/pwm/pwmchip0/pwm1/enable",
	})
	sysfs.SetFilesystem(fs)
	a.Connect()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/export: No such file"), true)
}

func TestAdaptorPwmEnableError(t *testing.T) {
	a := NewAdaptor()
	fs := sysfs.NewMockFilesystem([]string{
		"/sys/kernel/debug/gpio_debug/gpio13/current_pinmux",
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio13/direction",
		"/sys/class/gpio/gpio13/value",
		"/sys/class/gpio/gpio221/direction",
		"/sys/class/gpio/gpio221/value",
		"/sys/class/gpio/gpio253/direction",
		"/sys/class/gpio/gpio253/value",
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm1/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm1/period",
		//"/sys/class/pwm/pwmchip0/pwm1/enable",
	})
	sysfs.SetFilesystem(fs)
	a.Connect()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/pwm1/enable: No such file"), true)
}

func TestAdaptorPwmWritePinError(t *testing.T) {
	a, _ := initTestAdaptor()

	a.writeFile = func(string, []byte) (int, error) {
		return 0, errors.New("write error")
	}

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorPwmWriteError(t *testing.T) {
	a, fs := initTestAdaptor()

	fs.WithWriteError = true

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestAdaptorPwmReadError(t *testing.T) {
	a, fs := initTestAdaptor()

	fs.WithReadError = true

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, errors.New("read error"))
}

func TestAdaptorAnalog(t *testing.T) {
	a, fs := initTestAdaptor()

	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"
	i, _ := a.AnalogRead("0")
	gobottest.Assert(t, i, 250)
}

func TestAdaptorAnalogError(t *testing.T) {
	a, _ := initTestAdaptor()

	a.readFile = func(string) ([]byte, error) {
		return nil, errors.New("read error")
	}
	_, err := a.AnalogRead("0")
	gobottest.Assert(t, err, errors.New("read error"))
}
