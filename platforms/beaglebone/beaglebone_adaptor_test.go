package beaglebone

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/drivers/i2c"
	"gobot.io/x/gobot/v2/drivers/spi"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ aio.AnalogReader            = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ gpio.ServoWriter            = (*Adaptor)(nil)
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

const (
	pwmDir               = "/sys/devices/platform/ocp/48300000.epwmss/48300200.pwm/pwm/" //nolint:gosec // false positive
	pwmChip0Dir          = pwmDir + "pwmchip0/"
	pwmChip0ExportPath   = pwmChip0Dir + "export"
	pwmChip0UnexportPath = pwmChip0Dir + "unexport"
	pwmChip0Pwm0Dir      = pwmChip0Dir + "pwm0/"
	pwmChip0Pwm1Dir      = pwmChip0Dir + "pwm1/"

	pwm0EnablePath    = pwmChip0Pwm0Dir + "enable"
	pwm0PeriodPath    = pwmChip0Pwm0Dir + "period"
	pwm0DutyCyclePath = pwmChip0Pwm0Dir + "duty_cycle"
	pwm0PolarityPath  = pwmChip0Pwm0Dir + "polarity"

	pwm1EnablePath    = pwmChip0Pwm1Dir + "enable"
	pwm1PeriodPath    = pwmChip0Pwm1Dir + "period"
	pwm1DutyCyclePath = pwmChip0Pwm1Dir + "duty_cycle"
	pwm1PolarityPath  = pwmChip0Pwm1Dir + "polarity"
)

var pwmMockPaths = []string{
	"/sys/devices/platform/ocp/ocp:P9_22_pinmux/state",
	"/sys/devices/platform/ocp/ocp:P9_21_pinmux/state",
	"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
	pwmChip0ExportPath,
	pwmChip0UnexportPath,
	pwm0EnablePath,
	pwm0PeriodPath,
	pwm0DutyCyclePath,
	pwm0PolarityPath,
	pwm1EnablePath,
	pwm1PeriodPath,
	pwm1DutyCyclePath,
	pwm1PolarityPath,
}

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewAdaptor()
	// assert
	assert.IsType(t, &Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "Beaglebone"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.mutex)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.Equal(t, bbbPinMap, a.pinMap)
	assert.NotNil(t, a.pwmPinTranslate)
	assert.Equal(t, "/sys/class/leds/beaglebone:green:", a.usrLed)
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestPWMWrite(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	fs.Files[pwm1DutyCyclePath].Contents = "0"
	fs.Files[pwm1PeriodPath].Contents = "0"
	// act & assert wrong pin
	require.ErrorContains(t, a.PwmWrite("P9_99", 175), "'P9_99' is not a valid id for a PWM pin")

	// act & assert values
	_ = a.PwmWrite("P9_21", 175)
	assert.Equal(t, "500000", fs.Files[pwm1PeriodPath].Contents)
	assert.Equal(t, "343137", fs.Files[pwm1DutyCyclePath].Contents)

	require.NoError(t, a.Finalize())
}

func TestServoWrite(t *testing.T) {
	// arrange: prepare 50Hz for servos
	const (
		pin         = "P9_21"
		fiftyHzNano = 20000000
	)
	a := NewAdaptor(adaptors.WithPWMDefaultPeriodForPin(pin, fiftyHzNano))
	fs := a.sys.UseMockFilesystem(pwmMockPaths)
	require.NoError(t, a.Connect())
	// act & assert for 0° (min default value)
	err := a.ServoWrite(pin, 0)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwm1PeriodPath].Contents)
	assert.Equal(t, "500000", fs.Files[pwm1DutyCyclePath].Contents)
	// act & assert for 180° (max default value)
	err = a.ServoWrite(pin, 180)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwm1PeriodPath].Contents)
	assert.Equal(t, "2500000", fs.Files[pwm1DutyCyclePath].Contents)
	// act & assert invalid pins
	err = a.ServoWrite("3", 120)
	require.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")

	require.NoError(t, a.Finalize())
}

func TestAnalog(t *testing.T) {
	mockPaths := []string{
		"/sys/bus/iio/devices/iio:device0/in_voltage1_raw",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	fs.Files["/sys/bus/iio/devices/iio:device0/in_voltage1_raw"].Contents = "567\n"
	i, err := a.AnalogRead("P9_40")
	assert.Equal(t, 567, i)
	require.NoError(t, err)

	_, err = a.AnalogRead("P9_99")
	require.ErrorContains(t, err, "not a valid id for an analog pin")

	fs.WithReadError = true
	_, err = a.AnalogRead("P9_40")
	require.ErrorContains(t, err, "read error")
	fs.WithReadError = false

	require.NoError(t, a.Finalize())
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
	_ = a.DigitalWrite("usr1", 1)
	assert.Equal(t,
		"1",
		fs.Files["/sys/class/leds/beaglebone:green:usr1/brightness"].Contents,
	)

	// no such LED
	err := a.DigitalWrite("usr10101", 1)
	require.ErrorContains(t, err, " : /sys/class/leds/beaglebone:green:usr10101/brightness: no such file")

	_ = a.DigitalWrite("P9_12", 1)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio60/value"].Contents)

	require.ErrorContains(t, a.DigitalWrite("P9_99", 1), "'P9_99' is not a valid id for a digital pin")

	_, err = a.DigitalRead("P9_99")
	require.ErrorContains(t, err, "'P9_99' is not a valid id for a digital pin")

	fs.Files["/sys/class/gpio/gpio66/value"].Contents = "1"
	i, err := a.DigitalRead("P8_07")
	assert.Equal(t, 1, i)
	require.NoError(t, err)

	require.NoError(t, a.Finalize())
}

func TestAnalogReadFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/devices/platform/whatever",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	_, err := a.AnalogRead("P9_40")
	require.ErrorContains(t, err, "/sys/bus/iio/devices/iio:device0/in_voltage1_raw: no such file")
}

func TestDigitalPinDirectionFileError(t *testing.T) {
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/gpio60/value",
		"/sys/devices/platform/ocp/ocp:P9_12_pinmux/state",
	}

	a, _ := initTestAdaptorWithMockedFilesystem(mockPaths)

	err := a.DigitalWrite("P9_12", 1)
	require.ErrorContains(t, err, "/sys/class/gpio/gpio60/direction: no such file")

	// no pin added after previous problem, so no pin to unexport in finalize
	err = a.Finalize()
	require.NoError(t, err)
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
	require.NoError(t, err)

	err = a.Finalize()
	require.ErrorContains(t, err, "/sys/class/gpio/unexport: no such file")
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, 8, a.SpiDefaultBitCount())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewAdaptor()
	assert.Equal(t, 2, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-2"})
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 2)
	require.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.ErrorContains(t, err, "close error")
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

	tests := map[string]struct {
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
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantDir, path)
			assert.Equal(t, tc.wantChannel, channel)
			if tc.wantErr == nil {
				assert.Equal(t, "pwm", fs.Files[muxPath].Contents)
			}
		})
	}
}
