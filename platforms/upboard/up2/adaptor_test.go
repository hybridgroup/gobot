package up2

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
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
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)
var _ spi.Connector = (*Adaptor)(nil)

var pwmMockPaths = []string{
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm0/enable",
	"/sys/class/pwm/pwmchip0/pwm0/period",
	"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm0/polarity",
}

var gpioMockPaths = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio462/value",
	"/sys/class/gpio/gpio462/direction",
	"/sys/class/gpio/gpio432/value",
	"/sys/class/gpio/gpio432/direction",
	"/sys/class/leds/upboard:green:/brightness",
}

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "UP2"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	a.DigitalWrite("7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio462/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio432/value"].Contents = "1"
	i, _ := a.DigitalRead("13")
	gobottest.Assert(t, i, 1)

	a.DigitalWrite("green", 1)
	gobottest.Assert(t,
		fs.Files["/sys/class/leds/upboard:green:/brightness"].Contents,
		"1",
	)

	gobottest.Assert(t, a.DigitalWrite("99", 1), errors.New("'99' is not a valid id for a digital pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestSPI(t *testing.T) {
	a := NewAdaptor()

	gobottest.Assert(t, a.GetSpiDefaultBus(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMode(), 0)
	gobottest.Assert(t, a.GetSpiDefaultMaxSpeed(), int64(500000))

	_, err := a.GetSpiConnection(10, 0, 0, 8, 500000)
	gobottest.Assert(t, err.Error(), "Bus number 10 out of range")

	// TODO: tests for real connection currently not possible, because not using system.Accessor using
	// TODO: test tx/rx here...
}

func TestPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)

	err := a.PwmWrite("32", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("32", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("32", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	gobottest.Assert(t, a.DigitalWrite("7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)

	gobottest.Assert(t, a.PwmWrite("32", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.DefaultBus(), 5)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-5"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetConnection(0xff, 5)
	gobottest.Assert(t, err, nil)
	_, err = con.Write([]byte{0xbf})
	gobottest.Assert(t, err, nil)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	gobottest.Assert(t, strings.Contains(err.Error(), "close error"), true)
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
		"number_4_error": {
			busNr:   4,
			wantErr: fmt.Errorf("Bus number 4 out of range"),
		},
		"number_5_ok": {
			busNr: 5,
		},
		"number_6_ok": {
			busNr: 6,
		},
		"number_7_error": {
			busNr:   7,
			wantErr: fmt.Errorf("Bus number 7 out of range"),
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

func Test_translatePWMPin(t *testing.T) {
	var tests = map[string]struct {
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"16": {
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 3,
		},
		"32": {
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 0,
		},
		"33": {
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 1,
		},
		"PWM0": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'PWM0' is not a valid id for a pin"),
		},
		"7": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'7' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			dir, channel, err := a.translatePWMPin(name)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, dir, tc.wantDir)
			gobottest.Assert(t, channel, tc.wantChannel)
		})
	}
}
