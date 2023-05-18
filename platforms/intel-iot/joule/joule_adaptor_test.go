package joule

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
var _ i2c.Connector = (*Adaptor)(nil)

func initTestAdaptorWithMockedFilesystem() (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	mockPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm0/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm0/period",
		"/sys/class/pwm/pwmchip0/pwm0/enable",
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio13/value",
		"/sys/class/gpio/gpio13/direction",
		"/sys/class/gpio/gpio40/value",
		"/sys/class/gpio/gpio40/direction",
		"/sys/class/gpio/gpio446/value",
		"/sys/class/gpio/gpio446/direction",
		"/sys/class/gpio/gpio463/value",
		"/sys/class/gpio/gpio463/direction",
		"/sys/class/gpio/gpio421/value",
		"/sys/class/gpio/gpio421/direction",
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
		"/sys/class/gpio/gpio451/direction",
		"/sys/class/gpio/gpio451/value",
		"/dev/i2c-0",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "5000"
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestName(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	gobottest.Assert(t, strings.HasPrefix(a.Name(), "Joule"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestFinalize(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("J12_1", 1)
	a.PwmWrite("J12_26", 100)

	gobottest.Assert(t, a.Finalize(), nil)

	// assert finalize after finalize is working
	gobottest.Assert(t, a.Finalize(), nil)

	// assert re-connect is working
	gobottest.Assert(t, a.Connect(), nil)
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	a.DigitalWrite("J12_1", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio451/value"].Contents, "1")

	a.DigitalWrite("J12_1", 0)

	i, err := a.DigitalRead("J12_1")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 0)

	_, err = a.DigitalRead("P9_99")
	gobottest.Assert(t, err, errors.New("'P9_99' is not a valid id for a digital pin"))
}

func TestPwm(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents, "3921568")

	err = a.PwmWrite("4", 100)
	gobottest.Assert(t, err, errors.New("'4' is not a valid id for a pin"))

	err = a.PwmWrite("J12_1", 100)
	gobottest.Assert(t, err, errors.New("'J12_1' is not a valid id for a PWM pin"))
}

func TestPwmPinExportError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/export")

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/export: no such file"), true)
}

func TestPwmPinEnableError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/pwm0/enable")

	err := a.PwmWrite("J12_26", 100)
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/pwm/pwmchip0/pwm0/enable: no such file"), true)
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	gobottest.Assert(t, a.DefaultI2cBus(), 0)
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
