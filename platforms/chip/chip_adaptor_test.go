package chip

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
	_ gpio.ServoWriter            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
)

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
	assert.True(t, strings.HasPrefix(a.Name(), "CHIP"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewProAdaptor(t *testing.T) {
	a := NewProAdaptor()
	assert.True(t, strings.HasPrefix(a.Name(), "CHIP Pro"))
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.DigitalWrite("CSID7", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents = "0"
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "0"

	assert.NoError(t, a.Connect())
	assert.NoError(t, a.PwmWrite("PWM0", 100))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	_ = a.Connect()

	_ = a.DigitalWrite("CSID7", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio139/value"].Contents)

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	assert.Equal(t, 1, i)

	assert.ErrorContains(t, a.DigitalWrite("XIO-P10", 1), "'XIO-P10' is not a valid id for a digital pin")
	assert.NoError(t, a.Finalize())
}

func TestProDigitalIO(t *testing.T) {
	a, fs := initTestProAdaptorWithMockedFilesystem()
	_ = a.Connect()

	_ = a.DigitalWrite("CSID7", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio139/value"].Contents)

	fs.Files["/sys/class/gpio/gpio50/value"].Contents = "1"
	i, _ := a.DigitalRead("TWI2-SDA")
	assert.Equal(t, 1, i)

	assert.ErrorContains(t, a.DigitalWrite("XIO-P0", 1), "'XIO-P0' is not a valid id for a digital pin")
	assert.NoError(t, a.Finalize())
}

func TestPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem()
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents = "0"
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "0"

	_ = a.Connect()

	err := a.PwmWrite("PWM0", 100)
	assert.NoError(t, err)

	assert.Equal(t, "0", fs.Files["/sys/class/pwm/pwmchip0/export"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents)
	assert.Equal(t, "3921568", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents) // pwmPeriodDefault
	assert.Equal(t, "normal", fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents)

	err = a.ServoWrite("PWM0", 0)
	assert.NoError(t, err)

	assert.Equal(t, "500000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents)

	err = a.ServoWrite("PWM0", 180)
	assert.NoError(t, err)

	assert.Equal(t, "2000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents) // pwmPeriodDefault

	assert.NoError(t, a.Finalize())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 1, a.DefaultI2cBus())
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

func Test_translatePWMPin(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantChannel, channel)
		})
	}
}
