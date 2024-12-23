package adaptors

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

func TestDigitalPinsWithGpiosActiveLow(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val).apply(cfg)".
	// arrange & act, connect is mandatory to set options to the system
	a := NewDigitalPinsAdaptor(system.NewAccesser(), nil, WithGpiosActiveLow("1", "12", "33"))
	require.NoError(t, a.Connect())
	// assert
	assert.Len(t, a.pinOptions, 3)
}

func TestDigitalPinsWithDigitalPinInitializer(t *testing.T) {
	// arrange
	const (
		pinNo           = "1"
		translatedPinNo = "12"
	)
	a := NewDigitalPinsAdaptor(system.NewAccesser(), testDigitalPinTranslator)
	require.NoError(t, a.Connect())
	dpa := a.sys.UseMockDigitalPinAccess()
	pin, err := a.DigitalPin(pinNo)
	require.NoError(t, err)
	require.Equal(t, 1, dpa.Exported("", translatedPinNo)) // original initializer called on DigitalPin()
	require.NoError(t, a.digitalPinsCfg.initialize(pin))
	require.Equal(t, 2, dpa.Exported("", translatedPinNo))
	var called bool
	anotherInitializer := func(pin gobot.DigitalPinner) error {
		called = true
		return nil
	}
	WithDigitalPinInitializer(anotherInitializer).apply(a.digitalPinsCfg)
	// act
	require.NoError(t, a.digitalPinsCfg.initialize(pin))
	// assert
	assert.Equal(t, 2, dpa.Exported("", translatedPinNo))
	assert.True(t, called)
}

func TestDigitalPinsWithSysfsAccess(t *testing.T) {
	// arrange
	a := NewDigitalPinsAdaptor(system.NewAccesser(), nil)
	require.NoError(t, a.Connect())
	require.True(t, a.sys.IsGpiodDigitalPinAccess())
	require.NoError(t, a.Finalize())
	// act, connect is mandatory to set options to the system
	WithSysfsAccess().apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	// assert
	assert.True(t, a.sys.IsSysfsDigitalPinAccess())
}

func TestDigitalPinsWithGpiodAccess(t *testing.T) {
	// arrange
	a := NewDigitalPinsAdaptor(system.NewAccesser(system.WithDigitalPinSysfsAccess()), nil)
	require.NoError(t, a.Connect())
	require.True(t, a.sys.IsSysfsDigitalPinAccess())
	require.NoError(t, a.Finalize())
	// we have to mock the fs at this point to ensure the option can be applied on each test environment
	a.sys.UseMockFilesystem([]string{"/dev/gpiochip0"})
	// act, connect is mandatory to set options to the system
	WithGpiodAccess().apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	// assert
	assert.True(t, a.sys.IsGpiodDigitalPinAccess())
}

func TestDigitalReadWithGpiosActiveLow(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio25/value",
		"/sys/class/gpio/gpio25/direction",
		"/sys/class/gpio/gpio25/active_low",
		"/sys/class/gpio/gpio26/value",
		"/sys/class/gpio/gpio26/direction",
	}
	a, fs := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)
	// arrange files for for pin 14
	fs.Files["/sys/class/gpio/gpio25/value"].Contents = "1"
	fs.Files["/sys/class/gpio/gpio25/active_low"].Contents = "5"
	// arrange value file and config for pin 15
	fs.Files["/sys/class/gpio/gpio26/value"].Contents = "0"
	WithGpiosActiveLow("14").apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	// creates a new pin without inverted logic
	if _, err := a.DigitalRead("15"); err != nil {
		panic(err)
	}
	// assert for untouched content of pin 14
	assert.Equal(t, "5", fs.Files["/sys/class/gpio/gpio25/active_low"].Contents)
	// arrange direction file and config for pin 15
	fs.Add("/sys/class/gpio/gpio26/active_low")
	fs.Files["/sys/class/gpio/gpio26/active_low"].Contents = "6"
	WithGpiosActiveLow("15").apply(a.digitalPinsCfg)
	require.NoError(t, a.Finalize())
	require.NoError(t, a.Connect())
	// act
	got1, err1 := a.DigitalRead("14") // for a new pin
	got2, err2 := a.DigitalRead("15") // for an existing pin (calls ApplyOptions())
	// assert
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.Equal(t, 1, got1) // there is no mechanism to negate mocked values
	assert.Equal(t, 0, got2)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio25/active_low"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio26/active_low"].Contents)
}

func TestDigitalWriteWithOptions(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio18/value",
		"/sys/class/gpio/gpio18/direction",
	}
	a, fs := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)

	// assert write correct value without error and just ignore unsupported options
	WithGpiosPullUp("7").apply(a.digitalPinsCfg)
	WithGpiosOpenDrain("7").apply(a.digitalPinsCfg)
	WithGpioEventOnFallingEdge("7", gpioTestEventHandler).apply(a.digitalPinsCfg)
	WithGpioPollForEdgeDetection("7", 0, nil).apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	err := a.DigitalWrite("7", 1)
	require.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio18/value"].Contents)

	// assert second write to same pin without error and just ignore unsupported options
	WithGpiosPullDown("7").apply(a.digitalPinsCfg)
	WithGpiosOpenSource("7").apply(a.digitalPinsCfg)
	WithGpioDebounce("7", 2*time.Second).apply(a.digitalPinsCfg)
	WithGpioEventOnRisingEdge("7", gpioTestEventHandler).apply(a.digitalPinsCfg)
	require.NoError(t, a.Finalize())
	require.NoError(t, a.Connect())
	err = a.DigitalWrite("7", 1)
	require.NoError(t, err)

	// assert third write to same pin without error
	WithGpioEventOnBothEdges("7", gpioTestEventHandler).apply(a.digitalPinsCfg)
	require.NoError(t, a.Finalize())
	require.NoError(t, a.Connect())
	err = a.DigitalWrite("7", 1)
	require.NoError(t, err)

	// assert error on bad id
	require.ErrorContains(t, a.DigitalWrite("notexist", 1), "not a valid pin")

	// assert error bubbling
	fs.WithWriteError = true
	err = a.DigitalWrite("7", 0)
	require.ErrorContains(t, err, "write error")
}

func TestDigitalWriteWithGpiosActiveLow(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio19/value",
		"/sys/class/gpio/gpio19/direction",
		"/sys/class/gpio/gpio19/active_low",
	}
	a, fs := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/gpio/gpio19/active_low"].Contents = "5"
	WithGpiosActiveLow("8").apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	// act
	err := a.DigitalWrite("8", 2)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2", fs.Files["/sys/class/gpio/gpio19/value"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio19/active_low"].Contents)
}

func TestDigitalPinsWithSpiGpioAccess(t *testing.T) {
	// arrange
	const (
		sclkPin           = "1"
		ncsPin            = "2"
		sdoPin            = "3"
		sdiPin            = "4"
		sclkPinTranslated = "12"
		ncsPinTranslated  = "13"
		sdoPinTranslated  = "14"
		sdiPinTranslated  = "15"
	)
	a := NewDigitalPinsAdaptor(system.NewAccesser(), testDigitalPinTranslator)
	dpa := a.sys.UseMockDigitalPinAccess()
	// act
	WithSpiGpioAccess(sclkPin, ncsPin, sdoPin, sdiPin).apply(a.digitalPinsCfg)
	require.NoError(t, a.Connect())
	bus, err := a.sys.NewSpiDevice(0, 0, 0, 0, 1111)
	// assert
	require.NoError(t, err)
	assert.NotNil(t, bus)
	assert.Equal(t, 1, dpa.AppliedOptions("", sclkPinTranslated))
	assert.Equal(t, 1, dpa.AppliedOptions("", ncsPinTranslated))
	assert.Equal(t, 1, dpa.AppliedOptions("", sdoPinTranslated))
	assert.Equal(t, 0, dpa.AppliedOptions("", sdiPinTranslated)) // already input, so no option applied
	assert.Equal(t, 1, dpa.Exported("", sclkPinTranslated))
	assert.Equal(t, 1, dpa.Exported("", ncsPinTranslated))
	assert.Equal(t, 1, dpa.Exported("", sdoPinTranslated))
	assert.Equal(t, 1, dpa.Exported("", sdiPinTranslated))
}

func gpioTestEventHandler(o int, t time.Duration, et string, sn uint32, lsn uint32) {
	// the handler should never execute, because used in outputs and not supported by sysfs
	panic(fmt.Sprintf("event handler was called (%d, %d) unexpected for line %d with '%s' at %s!", sn, lsn, o, t, et))
}
