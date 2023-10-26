package up2

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/drivers/spi"
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
	_ spi.Connector               = (*Adaptor)(nil)
)

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
	assert.True(t, strings.HasPrefix(a.Name(), "UP2"))
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestDigitalIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	_ = a.DigitalWrite("7", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio462/value"].Contents)

	fs.Files["/sys/class/gpio/gpio432/value"].Contents = "1"
	i, _ := a.DigitalRead("13")
	assert.Equal(t, 1, i)

	_ = a.DigitalWrite("green", 1)
	assert.Equal(t,
		"1",
		fs.Files["/sys/class/leds/upboard:green:/brightness"].Contents,
	)

	assert.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	assert.NoError(t, a.Finalize())
}

func TestPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents = "0"
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "0"

	err := a.PwmWrite("32", 100)
	assert.NoError(t, err)

	assert.Equal(t, "0", fs.Files["/sys/class/pwm/pwmchip0/export"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/pwm/pwmchip0/pwm0/enable"].Contents)
	assert.Equal(t, "3921568", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents) // pwmPeriodDefault
	assert.Equal(t, "normal", fs.Files["/sys/class/pwm/pwmchip0/pwm0/polarity"].Contents)

	err = a.ServoWrite("32", 0)
	assert.NoError(t, err)

	assert.Equal(t, "500000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents)

	err = a.ServoWrite("32", 180)
	assert.NoError(t, err)

	assert.Equal(t, "2000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents)
	assert.Equal(t, "10000000", fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents)

	assert.NoError(t, a.Finalize())
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	assert.NoError(t, a.DigitalWrite("7", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/duty_cycle"].Contents = "0"
	fs.Files["/sys/class/pwm/pwmchip0/pwm0/period"].Contents = "0"

	assert.NoError(t, a.PwmWrite("32", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	assert.Contains(t, err.Error(), "write error")
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func Test_validateSpiBusNumber(t *testing.T) {
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 5, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-5"})
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 5)
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func Test_translatePWMPin(t *testing.T) {
	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantChannel, channel)
		})
	}
}
