package jetson

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	require.NoError(t, a.Finalize())
}

func TestPWMPinsConnect(t *testing.T) {
	a := NewAdaptor()

	err := a.PwmWrite("33", 1)
	require.ErrorContains(t, err, "not connected")

	err = a.Connect()
	require.NoError(t, err)
}

func TestPWMPinsReConnect(t *testing.T) {
	// arrange
	mockPaths := []string{
		"/sys/class/pwm/pwmchip0/export",
		"/sys/class/pwm/pwmchip0/unexport",
		"/sys/class/pwm/pwmchip0/pwm2/duty_cycle",
		"/sys/class/pwm/pwmchip0/pwm2/period",
		"/sys/class/pwm/pwmchip0/pwm2/polarity",
		"/sys/class/pwm/pwmchip0/pwm2/enable",
	}
	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)
	require.NoError(t, a.PwmWrite("33", 1))
	require.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	require.NoError(t, err)
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
	require.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio216/value"].Contents)

	err = a.DigitalWrite("13", 1)
	require.NoError(t, err)
	i, err := a.DigitalRead("13")
	require.NoError(t, err)
	assert.Equal(t, 1, i)

	require.ErrorContains(t, a.DigitalWrite("notexist", 1), "'notexist' is not a valid id for a digital pin")
	require.NoError(t, a.Finalize())
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
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 1)
	require.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.ErrorContains(t, err, "close error")
}

func Test_translatePWMPin(t *testing.T) {
	tests := map[string]struct {
		pin         string
		wantDir     string
		wantChannel int
		wantErr     error
	}{
		"32_pwm0": {
			pin:         "32",
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 0,
		},
		"33_pwm2": {
			pin:         "33",
			wantDir:     "/sys/class/pwm/pwmchip0",
			wantChannel: 2,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// arrange
			a := NewAdaptor()
			// act
			dir, channel, err := a.translatePWMPin(tc.pin)
			// assert
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, dir)
			assert.Equal(t, tc.wantChannel, channel)
		})
	}
}
