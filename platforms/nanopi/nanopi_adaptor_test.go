package nanopi

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
	pwmDir           = "/sys/devices/platform/soc/1c21400.pwm/pwm/pwmchip0/" //nolint:gosec // false positive
	pwmExportPath    = pwmDir + "export"
	pwmUnexportPath  = pwmDir + "unexport"
	pwmPwmDir        = pwmDir + "pwm0/"
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

func initConnectedTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := initConnectedTestAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	return a, fs
}

func initConnectedTestAdaptor() *Adaptor {
	a := NewNeoAdaptor()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewNeoAdaptor()
	// assert
	assert.IsType(t, &Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "NanoPi NEO Board"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.PWMPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.True(t, a.sys.HasDigitalPinCdevAccess())
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewNeoAdaptor(adaptors.WithGpiosActiveLow("1"), adaptors.WithGpioSysfsAccess())
	// assert
	require.NoError(t, a.Connect())
	assert.True(t, a.sys.HasDigitalPinSysfsAccess())
}

func TestDigitalIO(t *testing.T) {
	// some basic tests, further tests are done in "digitalpinsadaptor.go"
	// arrange
	a := initConnectedTestAdaptor()
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.HasDigitalPinCdevAccess())
	// act & assert write
	err := a.DigitalWrite("7", 1)
	require.NoError(t, err)
	assert.Equal(t, []int{1}, dpa.Written("gpiochip0", "203"))
	// arrange, act & assert read
	dpa.UseValues("gpiochip0", "199", []int{2})
	i, err := a.DigitalRead("10")
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	// act and assert unknown pin
	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	// act and assert finalize
	require.NoError(t, a.Finalize())
	assert.Equal(t, 0, dpa.Exported("gpiochip0", "203"))
	assert.Equal(t, 0, dpa.Exported("gpiochip0", "199"))
}

func TestDigitalIOSysfs(t *testing.T) {
	// some basic tests, further tests are done in "digitalpinsadaptor.go"
	// arrange
	a := NewNeoAdaptor(adaptors.WithGpioSysfsAccess())
	require.NoError(t, a.Connect())
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.HasDigitalPinSysfsAccess())
	// act & assert write
	err := a.DigitalWrite("7", 1)
	require.NoError(t, err)
	assert.Equal(t, []int{1}, dpa.Written("", "203"))
	// arrange, act & assert read
	dpa.UseValues("", "199", []int{2})
	i, err := a.DigitalRead("10")
	require.NoError(t, err)
	assert.Equal(t, 2, i)
	// act and assert unknown pin
	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	// act and assert finalize
	require.NoError(t, a.Finalize())
	assert.Equal(t, 0, dpa.Exported("", "203"))
	assert.Equal(t, 0, dpa.Exported("", "199"))
}

func TestAnalog(t *testing.T) {
	mockPaths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
	}

	a, fs := initConnectedTestAdaptorWithMockedFilesystem(mockPaths)

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

func TestInvalidPWMPin(t *testing.T) {
	a, fs := initConnectedTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	err := a.PwmWrite("666", 42)
	require.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	err = a.ServoWrite("666", 120)
	require.ErrorContains(t, err, "'666' is not a valid id for a PWM pin")

	err = a.PwmWrite("3", 42)
	require.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")

	err = a.ServoWrite("3", 120)
	require.ErrorContains(t, err, "'3' is not a valid id for a PWM pin")
}

func TestPwmWrite(t *testing.T) {
	// arrange
	a, fs := initConnectedTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)
	// act
	err := a.PwmWrite("PWM", 100)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, strconv.Itoa(10000000), fs.Files[pwmPeriodPath].Contents)
	assert.Equal(t, "3921568", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	require.NoError(t, a.Finalize())
}

func TestServoWrite(t *testing.T) {
	// arrange: prepare 50Hz for servos
	const (
		pin         = "PWM"
		fiftyHzNano = 20000000
	)
	a := NewNeoAdaptor(adaptors.WithPWMDefaultPeriodForPin(pin, fiftyHzNano))
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
	a, fs := initConnectedTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	newPeriod := uint32(2550000)
	// act
	err := a.SetPeriod("PWM", newPeriod)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "0", fs.Files[pwmExportPath].Contents)
	assert.Equal(t, "1", fs.Files[pwmEnablePath].Contents)
	assert.Equal(t, fmt.Sprintf("%d", newPeriod), fs.Files[pwmPeriodPath].Contents) //nolint:perfsprint // ok here
	assert.Equal(t, "0", fs.Files[pwmDutyCyclePath].Contents)
	assert.Equal(t, "normal", fs.Files[pwmPolarityPath].Contents)

	// arrange test for automatic adjustment of duty cycle to lower value
	err = a.PwmWrite("PWM", 127) // 127 is a little bit smaller than 50% of period
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(1270000), fs.Files[pwmDutyCyclePath].Contents)
	newPeriod = newPeriod / 10

	// act
	err = a.SetPeriod("PWM", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(127000), fs.Files[pwmDutyCyclePath].Contents)

	// arrange test for automatic adjustment of duty cycle to higher value
	newPeriod = newPeriod * 20

	// act
	err = a.SetPeriod("PWM", newPeriod)

	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(2540000), fs.Files[pwmDutyCyclePath].Contents)
}

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	// arrange
	a := initConnectedTestAdaptor()
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.HasDigitalPinCdevAccess())
	require.NoError(t, a.DigitalWrite("7", 1))
	dpa.UseUnexportError("gpiochip0", "203")
	// act
	err := a.Finalize()
	// assert
	require.ErrorContains(t, err, "unexport error")
}

func TestFinalizeErrorAfterPWM(t *testing.T) {
	a, fs := initConnectedTestAdaptorWithMockedFilesystem(pwmMockPaths)
	preparePwmFs(fs)

	require.NoError(t, a.PwmWrite("PWM", 1))

	fs.WithWriteError = true

	err := a.Finalize()
	require.ErrorContains(t, err, "write error")
}

func TestSpiDefaultValues(t *testing.T) {
	a := NewNeoAdaptor()

	assert.Equal(t, 0, a.SpiDefaultBusNumber())
	assert.Equal(t, 0, a.SpiDefaultChipNumber())
	assert.Equal(t, 0, a.SpiDefaultMode())
	assert.Equal(t, 8, a.SpiDefaultBitCount())
	assert.Equal(t, int64(500000), a.SpiDefaultMaxSpeed())
}

func TestI2cDefaultBus(t *testing.T) {
	a := NewNeoAdaptor()
	assert.Equal(t, 0, a.DefaultI2cBus())
}

func TestI2cFinalizeWithErrors(t *testing.T) {
	// arrange
	a := NewNeoAdaptor()
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
