package tinkerboard

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
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/system"
)

const (
	gpio17Path  = "/sys/class/gpio/gpio17/"
	gpio160Path = "/sys/class/gpio/gpio160/"
)

const (
	pwmDir           = "/sys/devices/platform/ff680020.pwm/pwm/pwmchip2/" //nolint:gosec // false positive
	pwmPwmDir        = pwmDir + "pwm0/"
	pwmExportPath    = pwmDir + "export"
	pwmUnexportPath  = pwmDir + "unexport"
	pwmEnablePath    = pwmPwmDir + "enable"
	pwmPeriodPath    = pwmPwmDir + "period"
	pwmDutyCyclePath = pwmPwmDir + "duty_cycle"
	pwmPolarityPath  = pwmPwmDir + "polarity"

	pwmInvertedIdentifier = "inversed"
)

var pwmMockPaths = []string{
	pwmExportPath,
	pwmUnexportPath,
	pwmEnablePath,
	pwmPeriodPath,
	pwmDutyCyclePath,
	pwmPolarityPath,
}

var gpioMockPaths = []string{
	"/sys/class/gpio/export",
	"/sys/class/gpio/unexport",
	gpio17Path + "value",
	gpio17Path + "direction",
	gpio160Path + "value",
	gpio160Path + "direction",
}

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gobot.PWMPinnerProvider     = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ gpio.PwmWriter              = (*Adaptor)(nil)
	_ gpio.ServoWriter            = (*Adaptor)(nil)
	_ aio.AnalogReader            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
)

func preparePwmFs(fs *system.MockFilesystem) {
	fs.Files[pwmEnablePath].Contents = "0"
	fs.Files[pwmPeriodPath].Contents = "0"
	fs.Files[pwmDutyCyclePath].Contents = "0"
	fs.Files[pwmPolarityPath].Contents = pwmInvertedIdentifier
}

func initTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := NewAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewAdaptor()
	// assert
	assert.IsType(t, &Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "Tinker Board"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.mutex)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.NotNil(t, a.OneWireBusAdaptor)
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewAdaptor(adaptors.WithGpiosActiveLow("1"))
	// assert
	require.NoError(t, a.Connect())
}

func TestDigitalIO(t *testing.T) {
	// only basic tests needed, further tests are done in "digitalpinsadaptor.go"
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	_ = a.DigitalWrite("7", 1)
	assert.Equal(t, "1", fs.Files[gpio17Path+"value"].Contents)

	fs.Files[gpio160Path+"value"].Contents = "1"
	i, _ := a.DigitalRead("10")
	assert.Equal(t, 1, i)

	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	require.NoError(t, a.Finalize())
}

func TestAnalogRead(t *testing.T) {
	mockPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
	}

	a, fs := initTestAdaptorWithMockedFilesystem(mockPaths)

	fs.Files["/sys/class/thermal/thermal_zone0/temp"].Contents = "567\n"
	got, err := a.AnalogRead("thermal_zone0")
	require.NoError(t, err)
	assert.Equal(t, 567, got)

	_, err = a.AnalogRead("thermal_zone10")
	require.ErrorContains(t, err, "'thermal_zone10' is not a valid id for an analog pin")

	fs.WithReadError = true
	_, err = a.AnalogRead("thermal_zone0")
	require.ErrorContains(t, err, "read error")
	fs.WithReadError = false

	require.NoError(t, a.Finalize())
}

func TestPwmWrite(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)
	// act
	err := a.PwmWrite("33", 100)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, "10000000", fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)
	// act & assert invalid pin
	err = a.PwmWrite("666", 42)
	require.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	require.NoError(t, a.Finalize())
}

func TestServoWrite(t *testing.T) {
	// arrange: prepare 50Hz for servos
	const (
		pin         = "33"
		fiftyHzNano = 20000000
	)
	a := NewAdaptor(adaptors.WithPWMDefaultPeriodForPin(pin, fiftyHzNano))
	fs := a.sys.UseMockFilesystem(pwmMockPaths)
	preparePwmFs(fs)
	require.NoError(t, a.Connect())
	// act & assert for 0° (min default value)
	err := a.ServoWrite(pin, 0)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "500000", fs.Files[pwmDutyCyclePath].Contents)
	// act & assert for 180° (max default value)
	err = a.ServoWrite(pin, 180)
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(fiftyHzNano), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "2500000", fs.Files[pwmDutyCyclePath].Contents)
	// act & assert invalid pins
	err = a.ServoWrite("3", 120)
	require.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")

	require.NoError(t, a.Finalize())
}

func TestSetPeriod(t *testing.T) {
	// arrange
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("33", newPeriod)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", newPeriod), fs.Files[pwmPeriodPath].Contents) //nolint:perfsprint // ok here
	assert.Equal(t, "0", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("33", 127) // 127 is a little bit smaller than 50% of period
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(1270000), fs.Files[pwmDutyCyclePath].Contents)
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(127000), fs.Files[pwmDutyCyclePath].Contents)

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("33", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(2540000), fs.Files[pwmDutyCyclePath].Contents)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(gpioMockPaths)

	require.NoError(t, a.DigitalWrite("7", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	require.NoError(t, a.PwmWrite("33", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
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
	assert.Equal(t, 1, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-4"})
	require.NoError(t, a.Connect())
	con, err := a.GetI2cConnection(0xff, 4)
	require.NoError(t, err)
	_, err = con.Write([]byte{0xbf})
	require.NoError(t, err)
	fs.WithCloseError = true
	// act
	err = a.Finalize()
	// assert
	require.ErrorContains(t, err, "close error")
}
