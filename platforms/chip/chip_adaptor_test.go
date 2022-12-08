package chip

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot"
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
var _ gpio.PwmWriter = (*Adaptor)(nil)
var _ gpio.ServoWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

var mockPaths = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	"/sys/class/gpio/gpio50/value",
	"/sys/class/gpio/gpio50/direction",
	"/sys/class/gpio/gpio139/value",
	"/sys/class/gpio/gpio139/direction",
	"/sys/class/pwm/pwmchip0/export",
	"/sys/class/pwm/pwmchip0/unexport",
	"/sys/class/pwm/pwmchip0/pwm0/enable",
	"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
	"/sys/class/pwm/pwmchip0/pwm0/polarity",
	"/sys/class/pwm/pwmchip0/pwm0/period",
}

func initTestAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func initTestProAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewProAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "CHIP"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestNewProAdaptor(t *testing.T) {
	a := NewProAdaptor()
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "CHIP Pro"), true)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("CSID7", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.PwmWrite("PWM0", 100), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P10", 1), errors.New("'XIO-P10' is not a valid id for a digital pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestProDigitalIO(t *testing.T) {
	a, fs := initTestProAdaptorWithMockedFilesystem()
	a.Connect()

	a.DigitalWrite("CSID7", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio139/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("XIO-P0", 1), errors.New("'XIO-P0' is not a valid id for a digital pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	a.Connect()

	err := a.PwmWrite("PWM0", 100)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/export"].Contents, "0")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents, "normal")

	err = a.ServoWrite("PWM0", 0)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "500000")

	err = a.ServoWrite("PWM0", 180)
	gobottest.Assert(t, err, nil)

	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "2000000")
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.GetDefaultBus(), 1)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-2"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetConnection(0xff, 2)
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
		"number_0_ok": {
			busNr: 0,
		},
		"number_1_ok": {
			busNr: 1,
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

func Test_translatePWMPin(t *testing.T) {
	var tests = map[string]struct {
		usePro      bool
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"PWM0": {
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 0,
		},
		"PWM1": {
			usePro:      true,
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 1,
		},
		"33_1": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'33_1' is not a valid id for a pin"),
		},
		"AP-EINT3": {
			wantDir:     "",
			wantChannel: -1,
			wantErr:     fmt.Errorf("'AP-EINT3' is not a valid id for a PWM pin"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			var a *Adaptor
			if tc.usePro {
				a = NewProAdaptor()
			} else {
				a = NewAdaptor()
			}
			// act
			dir, channel, err := a.translatePWMPin(name)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
			gobottest.Assert(t, dir, tc.wantDir)
			gobottest.Assert(t, channel, tc.wantChannel)
		})
	}
}
