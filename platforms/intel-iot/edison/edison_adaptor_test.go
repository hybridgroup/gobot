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
	"gobot.io/x/gobot/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gobot.PWMPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
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

var pwmMockPathsMux13 = []string{
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
	"/sys/class/pwm/pwmchip0/pwm1/enable",
}

var pwmMockPathsMux40 = []string{
	"/sys/kernel/debug/gpio_debug/gpio40/current_pinmux",
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio40/value",
	"/sys/class/gpio/gpio40/direction",
	"/sys/class/gpio/gpio229/value", // resistor
	"/sys/class/gpio/gpio229/direction",
	"/sys/class/gpio/gpio243/value",
	"/sys/class/gpio/gpio243/direction",
	"/sys/class/gpio/gpio261/value", // level shifter
	"/sys/class/gpio/gpio261/direction",
}

func initTestAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(testPinFiles)
	fs.Files["/sys/class/pwm/pwmchip0/pwm1/period"].Contents = "5000"
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Edison"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestConnect(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.GetDefaultBus(), 6)
	gobottest.Assert(t, a.Board(), "arduino")
	gobottest.Assert(t, a.Connect(), nil)
}

func TestArduinoSetupFail263(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/gpio/gpio263/direction")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio263/direction: No such file"), true)
}

func TestArduinoSetupFail240(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/gpio/gpio240/direction")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio240/direction: No such file"), true)
}

func TestArduinoSetupFail111(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio111/current_pinmux: No such file"), true)
}

func TestArduinoSetupFail131(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux")

	err := a.arduinoSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio131/current_pinmux: No such file"), true)
}

func TestArduinoI2CSetupFailTristate(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.arduinoSetup(), nil)

	fs.WithWriteError = true
	err := a.arduinoI2CSetup()
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestArduinoI2CSetupFail14(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/gpio14/direction")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio14/direction: No such file"), true)
}

func TestArduinoI2CSetupUnexportFail(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/unexport")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport: No such file"), true)
}

func TestArduinoI2CSetupFail236(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/class/gpio/gpio236/direction")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio236/direction: No such file"), true)
}

func TestArduinoI2CSetupFail28(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, a.arduinoSetup(), nil)
	delete(fs.Files, "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux")

	err := a.arduinoI2CSetup()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/kernel/debug/gpio_debug/gpio28/current_pinmux: No such file"), true)
}

func TestConnectArduinoError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestConnectArduinoWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestConnectSparkfun(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	a.SetBoard("sparkfun")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
	gobottest.Assert(t, a.Board(), "sparkfun")
}

func TestConnectMiniboard(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	a.SetBoard("miniboard")
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.GetDefaultBus(), 1)
	gobottest.Assert(t, a.Board(), "miniboard")
}

func TestConnectUnknown(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	a.SetBoard("wha")

	err := a.Connect()
	gobottest.Assert(t, strings.Contains(err.Error(), "Unknown board type: wha"), true)
}

func TestFinalize(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("3", 1)
	a.PwmWrite("5", 100)

	a.GetConnection(0xff, 6)
	gobottest.Assert(t, a.Finalize(), nil)

	// remove one file to force Finalize error
	delete(fs.Files, "/sys/class/gpio/unexport")
	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "4 errors occurred"), true)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport"), true)
}

func TestFinalizeError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	a.PwmWrite("5", 100)

	fs.WithWriteError = true
	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "6 errors occurred"), true)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
	gobottest.Assert(t, strings.Contains(err.Error(), "Enable(false) failed for pin 1 with write error"), true)
	gobottest.Assert(t, strings.Contains(err.Error(), "Unexport() failed for pin 1 with write error"), true)
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("13", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio40/value"].Contents, "1")

	a.DigitalWrite("2", 0)
	i, err := a.DigitalRead("2")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)
}

func TestDigitalPinInFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio40/value")
	delete(fs.Files, "/sys/class/gpio/gpio40/direction")
	a.Connect()

	_, err := a.DigitalPin("13")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)

}

func TestDigitalPinInResistorFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio229/value")
	delete(fs.Files, "/sys/class/gpio/gpio229/direction")
	a.Connect()

	_, err := a.DigitalPin("13")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestDigitalPinInLevelShifterFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio261/value")
	delete(fs.Files, "/sys/class/gpio/gpio261/direction")
	a.Connect()

	_, err := a.DigitalPin("13")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestDigitalPinInMuxFileError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux40)
	delete(fs.Files, "/sys/class/gpio/gpio243/value")
	delete(fs.Files, "/sys/class/gpio/gpio243/direction")
	a.Connect()

	_, err := a.DigitalPin("13")
	gobottest.Assert(t, strings.Contains(err.Error(), "No such file"), true)
}

func TestDigitalWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.DigitalWrite("13", 1)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestDigitalReadWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	_, err := a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestI2c(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()
	a.sys.UseMockSyscall()

	con, err := a.GetConnection(0xff, 6)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestI2cInvalidBus(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	_, err := a.GetConnection(0xff, 3)
	gobottest.Assert(t, err, errors.New("Unsupported I2C bus"))
}

func TestPwm(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm1/duty_cycle"].Contents, "1960")

	err = a.PwmWrite("7", 100)
	gobottest.Assert(t, err, errors.New("Not a PWM pin"))
}

func TestPwmExportError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux13)
	delete(fs.Files, "/sys/class/pwm/pwmchip0/export")
	a.Connect()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/export: No such file"), true)
}

func TestPwmEnableError(t *testing.T) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(pwmMockPathsMux13)
	delete(fs.Files, "/sys/class/pwm/pwmchip0/pwm1/enable")
	a.Connect()

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/pwm1/enable: No such file"), true)
}

func TestPwmWritePinError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, err, errors.New("write error"))
}

func TestPwmWriteError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithWriteError = true

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestPwmReadError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithReadError = true

	err := a.PwmWrite("5", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "read error"), true)
}

func TestAnalog(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.Files["/sys/bus/iio/devices/iio:device1/in_voltage0_raw"].Contents = "1000\n"

	i, _ := a.AnalogRead("0")
	gobottest.Assert(t, i, 250)
}

func TestAnalogError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.WithReadError = true

	_, err := a.AnalogRead("0")
	gobottest.Assert(t, err, errors.New("read error"))
}
