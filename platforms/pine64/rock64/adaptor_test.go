package rock64

import (
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

// make sure that this Adaptor fulfills all the required interfaces
var (
	_ gobot.Adaptor               = (*Adaptor)(nil)
	_ gobot.DigitalPinnerProvider = (*Adaptor)(nil)
	_ gpio.DigitalReader          = (*Adaptor)(nil)
	_ gpio.DigitalWriter          = (*Adaptor)(nil)
	_ aio.AnalogReader            = (*Adaptor)(nil)
	_ i2c.Connector               = (*Adaptor)(nil)
)

func initConnectedTestAdaptorWithMockedFilesystem(mockPaths []string) (*Adaptor, *system.MockFilesystem) {
	a := initConnectedTestAdaptor()
	fs := a.sys.UseMockFilesystem(mockPaths)
	return a, fs
}

func initConnectedTestAdaptor() *Adaptor {
	a := NewAdaptor()
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a
}

func TestNewAdaptor(t *testing.T) {
	// arrange & act
	a := NewAdaptor()
	// assert
	assert.IsType(t, &Adaptor{}, a)
	assert.True(t, strings.HasPrefix(a.Name(), "ROCK64"))
	assert.NotNil(t, a.sys)
	assert.NotNil(t, a.mutex)
	assert.NotNil(t, a.AnalogPinsAdaptor)
	assert.NotNil(t, a.DigitalPinsAdaptor)
	assert.NotNil(t, a.I2cBusAdaptor)
	assert.NotNil(t, a.SpiBusAdaptor)
	assert.True(t, a.sys.HasDigitalPinCdevAccess())
	// act & assert
	a.SetName("NewName")
	assert.Equal(t, "NewName", a.Name())
}

func TestNewAdaptorWithOption(t *testing.T) {
	// arrange & act
	a := NewAdaptor(adaptors.WithGpiosActiveLow("1"), adaptors.WithGpioSysfsAccess())
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
	assert.Equal(t, []int{1}, dpa.Written("gpiochip1", "28"))
	// arrange, act & assert read
	dpa.UseValues("gpiochip2", "1", []int{3})
	i, err := a.DigitalRead("10")
	require.NoError(t, err)
	assert.Equal(t, 3, i)
	// act and assert unknown pin
	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	// act and assert finalize
	require.NoError(t, a.Finalize())
	assert.Equal(t, 0, dpa.Exported("gpiochip1", "28"))
	assert.Equal(t, 0, dpa.Exported("gpiochip2", "1"))
}

func TestDigitalIOSysfs(t *testing.T) {
	// some basic tests, further tests are done in "digitalpinsadaptor.go"
	// arrange
	a := NewAdaptor(adaptors.WithGpioSysfsAccess())
	require.NoError(t, a.Connect())
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.HasDigitalPinSysfsAccess())
	// act & assert write
	err := a.DigitalWrite("7", 1)
	require.NoError(t, err)
	assert.Equal(t, []int{1}, dpa.Written("", "60"))
	// arrange, act & assert read
	dpa.UseValues("", "65", []int{4})
	i, err := a.DigitalRead("10")
	require.NoError(t, err)
	assert.Equal(t, 4, i)
	// act and assert unknown pin
	require.ErrorContains(t, a.DigitalWrite("99", 1), "'99' is not a valid id for a digital pin")
	// act and assert finalize
	require.NoError(t, a.Finalize())
	assert.Equal(t, 0, dpa.Exported("", "60"))
	assert.Equal(t, 0, dpa.Exported("", "65"))
}

func TestAnalogRead(t *testing.T) {
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

func TestFinalizeErrorAfterGPIO(t *testing.T) {
	// arrange
	a := initConnectedTestAdaptor()
	dpa := a.sys.UseMockDigitalPinAccess()
	require.True(t, a.sys.HasDigitalPinCdevAccess())
	require.NoError(t, a.DigitalWrite("7", 1))
	dpa.UseUnexportError("gpiochip1", "28")
	// act
	err := a.Finalize()
	// assert
	require.ErrorContains(t, err, "unexport error")
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
	a := initConnectedTestAdaptor()
	a.sys.UseMockSyscall()
	fs := a.sys.UseMockFilesystem([]string{"/dev/i2c-1"})
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
