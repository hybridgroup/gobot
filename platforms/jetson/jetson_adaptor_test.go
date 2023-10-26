package jetson

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
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
	_ i2c.Connector               = (*Adaptor)(nil)
	_ spi.Connector               = (*Adaptor)(nil)
)

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestNewAdaptor(t *testing.T) {
	a := NewAdaptor()

	assert.True(t, strings.HasPrefix(a.Name(), "JetsonNano"))

	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestFinalize(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/dev/i2c-1",
		"/dev/i2c-0",
		"/dev/spidev0.0",
		"/dev/spidev0.1",
	}
	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	_ = a.DigitalWrite("3", 1)

	_, _ = a.GetI2cConnection(0xff, 0)
	assert.NoError(t, a.Finalize())
}

func TestPWMPinsConnect(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, (map[string]gobot.PWMPinner)(nil), a.pwmPins)

	err := a.PwmWrite("33", 1)
	assert.ErrorContains(t, err, "not connected")

	err = a.Connect()
	assert.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.PWMPinner)(nil), a.pwmPins)
	assert.Equal(t, 0, len(a.pwmPins))
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	mockPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm2/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm2/period",
		"/sys/class/pwm/pwmchip0/pwm2/enable",
	}
	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)
	assert.Equal(t, 0, len(a.pwmPins))
	assert.NoError(t, a.PwmWrite("33", 1))
	assert.Equal(t, 1, len(a.pwmPins))
	assert.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(a.pwmPins))
}

func TestDigitalIO(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio216/value",
		"/sys/class/gpio/gpio216/direction",
		"/sys/class/gpio/gpio14/value",
		"/sys/class/gpio/gpio14/direction",
	}
	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("7", 1)
	assert.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio216/value"].Contents)

	err = a.DigitalWrite("13", 1)
	assert.NoError(t, err)
	i, err := a.DigitalRead("13")
	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	assert.ErrorContains(t, a.DigitalWrite("notexist", 1), "'notexist' is not a valid id for a digital pin")
	assert.NoError(t, a.Finalize())
}

func TestDigitalPinConcurrency(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(oldProcs)

	for retry := 0; retry < 20; retry++ {

		a := NewAdaptor()
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			pinAsString := strconv.Itoa(i)
			go func(pin string) {
				defer wg.Done()
				_, _ = a.DigitalPin(pin)
			}(pinAsString)
		}

		wg.Wait()
	}
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, int64(10000000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 1, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
	assert.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 1)
	assert.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	assert.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	assert.Contains(t, err.Error(), "close error")
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
		"number_2_not_ok": {
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
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
