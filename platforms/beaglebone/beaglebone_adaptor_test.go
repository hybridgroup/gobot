package beaglebone

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/sysfs"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ aio.AnalogReader = (*Adaptor)(nil)
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ sysfs.DigitalPinnerProvider = (*Adaptor)(nil)
var _ sysfs.PWMPinnerProvider = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)
var _ spi.Connector = (*Adaptor)(nil)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *sysfs.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sysfs.UseMockFilesystem(mockPaths)

	a.findPin = func(pinPath string) (string, error) {
		switch pinPath {
		case "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip*":
			return "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4", nil
		case "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip*":
			return "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2", nil
		case "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip*":
			return "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0", nil
		case "/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip*":
			return "/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip0", nil
		default:
			return pinPath, nil
		}
	}

	return a, fs
}

func TestGeneral(t *testing.T) {
	mockPaths := []string{
		"/dev/i2c-2",
		"/sys/devices/platform/bone_capemgr",
		"/sys/devices/platform/ocp/ocp:P8_07_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_11_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_22_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_21_pinmux/state",
		"/sys/class/leds/beaglebone:green:usr1/brightness",
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/export",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/unexport",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/enable",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/period",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm0/polarity",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/enable",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/polarity",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/class/gpio/gpio66/value",
		"/sys/class/gpio/gpio66/direction",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
		"/sys/class/gpio/gpio30/value",
		"/sys/class/gpio/gpio30/direction",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	// PWM
	gobottest.Assert(t, a.PwmWrite("P9_99", 175), errors.New("Not a valid PWM pin"))
	a.PwmWrite("P9_21", 175)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period"].Contents,
		"500000",
	)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle"].Contents,
		"343137",
	)

	a.ServoWrite("P9_21", 100)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/period"].Contents,
		"500000",
	)
	gobottest.Assert(
		t,
		fs.Files["/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/pwm1/duty_cycle"].Contents,
		"66666",
	)
	gobottest.Assert(t, a.ServoWrite("P9_99", 175), errors.New("Not a valid PWM pin"))

	fs.WithReadError = true
	gobottest.Assert(t, strings.Contains(a.PwmWrite("P9_21", 175).Error(), "read error"), true)
	fs.WithReadError = false

	fs.WithWriteError = true
	gobottest.Assert(t, strings.Contains(a.PwmWrite("P9_22", 175).Error(), "write error"), true)
	fs.WithWriteError = false

	// Analog
	fs.Files["/sys/bus/iio/devices/iio:device0/in_voltage1_raw"].Contents = "567\n"
	i, err := a.AnalogRead("P9_40")
	gobottest.Assert(t, i, 567)
	gobottest.Assert(t, err, nil)

	_, err = a.AnalogRead("P9_99")
	gobottest.Assert(t, err, errors.New("Not a valid analog pin"))

	fs.WithReadError = true
	_, err = a.AnalogRead("P9_40")
	gobottest.Assert(t, err, errors.New("read error"))
	fs.WithReadError = false

	// DigitalIO
	a.DigitalWrite("usr1", 1)
	gobottest.Assert(t,
		fs.Files["/sys/class/leds/beaglebone:green:usr1/brightness"].Contents,
		"1",
	)

	// no such LED
	err = a.DigitalWrite("usr10101", 1)
	gobottest.Assert(t, err.Error(), " : /sys/class/leds/beaglebone:green:usr10101/brightness: No such file.")

	a.DigitalWrite("P9_12", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio60/value"].Contents, "1")

	gobottest.Assert(t, a.DigitalWrite("P9_99", 1), errors.New("Not a valid pin"))

	_, err = a.DigitalRead("P9_99")
	gobottest.Assert(t, err, errors.New("Not a valid pin"))

	fs.Files["/sys/class/gpio/gpio66/value"].Contents = "1"
	i, err = a.DigitalRead("P8_07")
	gobottest.Assert(t, i, 1)
	gobottest.Assert(t, err, nil)

	fs.WithReadError = true
	_, err = a.DigitalRead("P8_07")
	gobottest.Assert(t, err, errors.New("read error"))
	fs.WithReadError = false

	fs.WithWriteError = true
	_, err = a.DigitalRead("P9_11")
	gobottest.Assert(t, err, errors.New("write error"))
	fs.WithWriteError = false

	// I2c
	a.sysfs.UseMockSyscall()

	con, err := a.GetConnection(0xff, 2)
	gobottest.Assert(t, err, nil)

	_, err = con.Write([]byte{0x00, 0x01})
	gobottest.Assert(t, err, nil)

	data := []byte{42, 42}
	_, err = con.Read(data)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, data, []byte{0x00, 0x01})

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestSPI(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultChip(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultBits(), 8)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(500000))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 10000000)
	gobottest.Assert(t, err.Error(), "Bus number 10 out of range")

	// TODO: tests for real connection currently not possible, because not using sysfs.Accessor using
	// TODO: test tx/rx here...
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Beaglebone"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 2)
}

func TestGetConnectionInvalidBus(t *testing.T) {
	a := NewAdaptor()
	_, err := a.GetConnection(0x01, 99)
	gobottest.Assert(t, err, errors.New("Bus number 99 out of range"))
}

func TestAnalogReadFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/whatever",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	_, err := a.AnalogRead("P9_40")
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/bus/iio/devices/iio:device0/in_voltage1_raw: No such file."), true)
}

func TestDigitalPinDirectionFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio60/value",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("P9_12", 1)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/gpio60/direction: No such file."), true)

	err = a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport: No such file."), true)
}

func TestDigitalPinFinalizeFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio60/value",
		"/sys/class/gpio/gpio60/direction",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("P9_12", 1)
	gobottest.Assert(t, err, nil)

	err = a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport: No such file."), true)
}

func TestPocketName(t *testing.T) {
	a := NewPocketBeagleAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "PocketBeagle"), true)
}
