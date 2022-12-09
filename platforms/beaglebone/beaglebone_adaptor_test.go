package beaglebone

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/drivers/spi"
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
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)
var _ spi.Connector = (*Adaptor)(nil)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestPWM(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/ocp/ocp:P9_22_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_21_pinmux/state",
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
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	gobottest.Assert(t, a.PwmWrite("P9_99", 175), errors.New("'P9_99' is not a valid id for a PWM pin"))
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

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestAnalog(t *testing.T) {
	mockPaths := []string{
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

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

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestDigitalIO(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/ocp/ocp:P8_07_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_11_pinmux/state",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
		"/sys/class/leds/beaglebone:green:usr1/brightness",
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

	// DigitalIO
	a.DigitalWrite("usr1", 1)
	gobottest.Assert(t,
		fs.Files["/sys/class/leds/beaglebone:green:usr1/brightness"].Contents,
		"1",
	)

	// no such LED
	err := a.DigitalWrite("usr10101", 1)
	gobottest.Assert(t, err.Error(), " : /sys/class/leds/beaglebone:green:usr10101/brightness: No such file.")

	a.DigitalWrite("P9_12", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio60/value"].Contents, "1")

	gobottest.Assert(t, a.DigitalWrite("P9_99", 1), errors.New("'P9_99' is not a valid id for a digital pin"))

	_, err = a.DigitalRead("P9_99")
	gobottest.Assert(t, err, errors.New("'P9_99' is not a valid id for a digital pin"))

	fs.Files["/sys/class/gpio/gpio66/value"].Contents = "1"
	i, err := a.DigitalRead("P8_07")
	gobottest.Assert(t, i, 1)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, a.Finalize(), nil)
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Beaglebone"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
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

	// no pin added after previous problem, so no pin to unexport in finalize
	err = a.Finalize()
	gobottest.Assert(t, nil, err)
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

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultChip(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultBits(), 8)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(500000))
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.DefaultI2cBus(), 2)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-2"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetI2cConnection(0xff, 2)
	gobottest.Assert(t, err, nil)
	_, err = con.Write([]byte{0xbf})
	gobottest.Assert(t, err, nil)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	gobottest.Assert(t, strings.Contains(err.Error(), "close error"), true)
}

func Test_validateSpiBusNumber(t *testing.T) {
	var tests = map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_ok": {
			busNr: 1,
		},
		"number_2_error": {
			busNr:   2,
			wantErr: fmt.Errorf("Bus number 2 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateSpiBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}

func Test_validateI2cBusNumber(t *testing.T) {
	var tests = map[string]struct {
		busNr   int
		wantErr error
	}{
		"number_negative_error": {
			busNr:   -1,
			wantErr: fmt.Errorf("Bus number -1 out of range"),
		},
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_error": {
			busNr:   1,
			wantErr: fmt.Errorf("Bus number 1 out of range"),
		},
		"number_2_ok": {
			busNr: 2,
		},
		"number_3_error": {
			busNr:   3,
			wantErr: fmt.Errorf("Bus number 3 out of range"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			err := a.validateI2cBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}

func Test_translateAndMuxPWMPin(t *testing.T) {
	// arrange
	mockPaths := []string{
		"/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4/",
		"/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2/",
		"/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0/",
		"/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip0/",
	}
	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	var tests = map[string]struct {
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"P8_13": {
			wantDir:     "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4",
			wantChannel: 1,
		},
		"P8_19": {
			wantDir:     "/sys/devices/platform/ocp/48304000.epwmss/48304200.pwm/pwm/pwmchip4",
			wantChannel: 0,
		},
		"P9_14": {
			wantDir:     "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2",
			wantChannel: 0,
		},
		"P9_16": {
			wantDir:     "/sys/devices/platform/ocp/48302000.epwmss/48302200.pwm/pwm/pwmchip2",
			wantChannel: 1,
		},
		"P9_21": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0",
			wantChannel: 1,
		},
		"P9_22": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/pwmchip0",
			wantChannel: 0,
		},
		"P9_42": {
			wantDir:     "/sys/devices/platform/ocp/48300000.epwmss/48300100.ecap/pwm/pwmchip0",
			wantChannel: 0,
		},
		"P9_99": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'P9_99' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			muxPath := fmt.Sprintf("/sys/devices/platform/ocp/ocp:%s_pinmux/state", name)
			fs.Add(muxPath)
			// act
			path, channel, err := a.translateAndMuxPWMPin(name)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, path, tc.wantDir)
			gobottest.Assert(t, channel, tc.wantChannel)
			if tc.wantErr == nil {
				gobottest.Assert(t, fs.Files[muxPath].Contents, "pwm")
			}
		})
	}
}
