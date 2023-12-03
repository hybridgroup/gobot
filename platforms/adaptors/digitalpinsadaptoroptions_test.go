//nolint:nonamedreturns // ok for tests
package adaptors

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2/system"
)

func TestDigitalPinsWithGpiosActiveLow(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val)(d)".
	// arrange
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()
	// act
	a := NewDigitalPinsAdaptor(sys, translate, WithGpiosActiveLow("1", "12", "33"))
	// assert
	assert.Len(t, a.pinOptions, 3)
}

func TestDigitalReadWithGpiosActiveLow(t *testing.T) {
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
	fs.Files["/sys/class/gpio/gpio25/value"].Contents = "1"
	fs.Files["/sys/class/gpio/gpio25/active_low"].Contents = "5"
	fs.Files["/sys/class/gpio/gpio26/value"].Contents = "0"
	WithGpiosActiveLow("14")(a)
	// creates a new pin without inverted logic
	if _, err := a.DigitalRead("15"); err != nil {
		panic(err)
	}
	fs.Add("/sys/class/gpio/gpio26/active_low")
	fs.Files["/sys/class/gpio/gpio26/active_low"].Contents = "6"
	WithGpiosActiveLow("15")(a)
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
	WithGpiosPullUp("7")(a)
	WithGpiosOpenDrain("7")(a)
	WithGpioEventOnFallingEdge("7", gpioEventHandler)(a)
	WithGpioPollForEdgeDetection("7", 0, nil)(a)
	err := a.DigitalWrite("7", 1)
	require.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio18/value"].Contents)

	// assert second write to same pin without error and just ignore unsupported options
	WithGpiosPullDown("7")(a)
	WithGpiosOpenSource("7")(a)
	WithGpioDebounce("7", 2*time.Second)(a)
	WithGpioEventOnRisingEdge("7", gpioEventHandler)(a)
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
	WithGpiosActiveLow("8")(a)
	// act
	err := a.DigitalWrite("8", 2)
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2", fs.Files["/sys/class/gpio/gpio19/value"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio19/active_low"].Contents)
}

func gpioEventHandler(o int, t time.Duration, et string, sn uint32, lsn uint32) {
	// the handler should never execute, because used in outputs and not supported by sysfs
	panic(fmt.Sprintf("event handler was called (%d, %d) unexpected for line %d with '%s' at %s!", sn, lsn, o, t, et))
}
