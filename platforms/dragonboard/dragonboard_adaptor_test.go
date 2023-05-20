package dragonboard

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/gobottest"
)

// make sure that this Adaptor fulfills all the required interfaces
var _ gobot.Adaptor = (*Adaptor)(nil)
var _ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
var _ gpio.DigitalReader = (*Adaptor)(nil)
var _ gpio.DigitalWriter = (*Adaptor)(nil)
var _ i2c.Connector = (*Adaptor)(nil)

func initTestAdaptor(t *testing.T) *Adaptor {
	a := NewAdaptor()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestName(t *testing.T) {
	a := initTestAdaptor(t)
	gobottest.Assert(t, strings.HasPrefix(a.Name(), "DragonBoard"), true)
	a.SetName("NewName")
	gobottest.Assert(t, a.Name(), "NewName")
}

func TestDigitalIO(t *testing.T) {
	a := initTestAdaptor(t)
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	_ = a.DigitalWrite("GPIO_B", 1)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio12/value"].Contents, "1")

	fs.Files["/sys/class/gpio/gpio36/value"].Contents = "1"
	i, _ := a.DigitalRead("GPIO_A")
	gobottest.Assert(t, i, 1)

	gobottest.Assert(t, a.DigitalWrite("GPIO_M", 1), errors.New("'GPIO_M' is not a valid id for a digital pin"))
	gobottest.Assert(t, a.Finalize(), nil)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a := initTestAdaptor(t)
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio36/value",
		"/sys/class/gpio/gpio36/direction",
		"/sys/class/gpio/gpio12/value",
		"/sys/class/gpio/gpio12/direction",
	}
	fs := a.sys.UseMockFilesystem(mockPaths)

	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("GPIO_B", 1), nil)

	fs.WithWriteError = true

	err := a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "write error"), true)
}

func TestI2cDefaultBus(t *testing.T) {
	a := initTestAdaptor(t)
	gobottest.Assert(t, a.DefaultI2cBus(), 0)
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	gobottest.Assert(t, a.Connect(), nil)
	con, err := a.GetI2cConnection(0xff, 1)
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
			err := a.validateI2cBusNumber(tc.busNr)
			// assert
			gobottest.Assert(t, err, tc.wantErr)
		})
	}
}
