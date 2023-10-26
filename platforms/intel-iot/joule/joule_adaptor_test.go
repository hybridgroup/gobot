package joule

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
)

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

	assert.True(t, strings.HasPrefix(a.Name(), "Joule"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestFinalize(t *testing.T) {
	a, _ := initTestAdaptorWithMockedFilesystem()

	_ = a.DigitalWrite("J12_1", 1)
	_ = a.PwmWrite("J12_26", 100)

	assert.NoError(t, a.Finalize())

	// assert finalize after finalize is working
	assert.NoError(t, a.Finalize())

	// assert re-connect is working
	assert.NoError(t, a.Connect())
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	_ = a.DigitalWrite("J12_1", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio451/value"].Contents)

	_ = a.DigitalWrite("J12_1", 0)

	i, err := a.DigitalRead("J12_1")
	assert.NoError(t, err)
	assert.Equal(t, 0, i)

	_, err = a.DigitalRead("P9_99")
	assert.ErrorContains(t, err, "'P9_99' is not a valid id for a digital pin")
}

func TestPwm(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()

	err := a.PwmWrite("J12_26", 100)
	assert.NoError(t, err)
	assert.Equal(t, "3921568", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)

	err = a.PwmWrite("4", 100)
	assert.ErrorContains(t, err, "'4' is not a valid id for a pin")

	err = a.PwmWrite("J12_1", 100)
	assert.ErrorContains(t, err, "'J12_1' is not a valid id for a PWM pin")
}

func TestPwmPinExportError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/export")

	err := a.PwmWrite("J12_26", 100)
	assert.Contains(t, err.Error(), "/sys/class/pwm/pwmchip0/export: no such file")
}

func TestPwmPinEnableError(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	delete(fs.Files, "/sys/class/pwm/pwmchip0/pwm0/enable")

	err := a.PwmWrite("J12_26", 100)
	assert.Contains(t, err.Error(), "/sys/class/pwm/pwmchip0/pwm0/enable: no such file")
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 0, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-2"})
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 2)
	assert.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	assert.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	assert.Contains(t, err.Error(), "close error")
}

func Test_validateI2cBusNumber(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
