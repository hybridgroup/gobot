//nolint:nonamedreturns // ok for tests
package adaptors

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/system"
)

// make sure that this adaptor fulfills all the required interfaces
var (
	_ gobot.DigitalPinnerProvider = (*DigitalPinsAdaptor)(nil)
	_ gpio.DigitalReader          = (*DigitalPinsAdaptor)(nil)
	_ gpio.DigitalWriter          = (*DigitalPinsAdaptor)(nil)
)

func initTestDigitalPinsAdaptorWithMockedFilesystem(mockPaths []string) (*DigitalPinsAdaptor, *system.MockFilesystem) {
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(mockPaths)
	a := NewDigitalPinsAdaptor(sys, testDigitalPinTranslator)
	if err := a.Connect(); err != nil {
		panic(err)
	}
	return a, fs
}

func testDigitalPinTranslator(pin string) (string, int, error) {
	line, err := strconv.Atoi(pin)
	if err != nil {
		return "", 0, fmt.Errorf("not a valid pin")
	}
	line = line + 11 // just for tests
	return "", line, err
}

func TestDigitalPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()

	a := NewDigitalPinsAdaptor(sys, translate)
	assert.Equal(t, (map[string]gobot.DigitalPinner)(nil), a.pins)

	_, err := a.DigitalRead("13")
	require.ErrorContains(t, err, "not connected for pin 13")

	err = a.DigitalWrite("7", 1)
	require.ErrorContains(t, err, "not connected for pin 7")

	err = a.Connect()
	require.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.DigitalPinner)(nil), a.pins)
	assert.Empty(t, a.pins)
}

func TestDigitalPinsFinalize(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio14/direction",
		"/sys/class/gpio/gpio14/value",
	}
	sys := system.NewAccesser()
	fs := sys.UseMockFilesystem(mockedPaths)
	a := NewDigitalPinsAdaptor(sys, testDigitalPinTranslator)
	// assert that finalize before connect is working
	require.NoError(t, a.Finalize())
	// arrange
	require.NoError(t, a.Connect())
	require.NoError(t, a.DigitalWrite("3", 1))
	assert.Len(t, a.pins, 1)
	// act
	err := a.Finalize()
	// assert
	require.NoError(t, err)
	assert.Empty(t, a.pins)
	// assert that finalize after finalize is working
	require.NoError(t, a.Finalize())
	// arrange missing sysfs file
	require.NoError(t, a.Connect())
	require.NoError(t, a.DigitalWrite("3", 2))
	delete(fs.Files, "/sys/class/gpio/unexport")
	err = a.Finalize()
	require.ErrorContains(t, err, "/sys/class/gpio/unexport: no such file")
}

func TestDigitalPinsReConnect(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio15/direction",
		"/sys/class/gpio/gpio15/value",
	}
	a, _ := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)
	require.NoError(t, a.DigitalWrite("4", 1))
	assert.Len(t, a.pins, 1)
	require.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	require.NoError(t, err)
	assert.NotNil(t, a.pins)
	assert.Empty(t, a.pins)
}

func TestDigitalIO(t *testing.T) {
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio25/value",
		"/sys/class/gpio/gpio25/direction",
	}
	a, _ := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)

	err := a.DigitalWrite("14", 1)
	require.NoError(t, err)

	i, err := a.DigitalRead("14")
	require.NoError(t, err)
	assert.Equal(t, 1, i)
}

func TestDigitalRead(t *testing.T) {
	// arrange
	mockedPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio24/value",
		"/sys/class/gpio/gpio24/direction",
	}
	a, fs := initTestDigitalPinsAdaptorWithMockedFilesystem(mockedPaths)
	fs.Files["/sys/class/gpio/gpio24/value"].Contents = "1"

	// assert read correct value without error
	i, err := a.DigitalRead("13")
	require.NoError(t, err)
	assert.Equal(t, 1, i)

	// assert error bubbling for read errors
	fs.WithReadError = true
	_, err = a.DigitalRead("13")
	require.ErrorContains(t, err, "read error")

	// assert error bubbling for write errors
	fs.WithWriteError = true
	_, err = a.DigitalRead("7")
	require.ErrorContains(t, err, "write error")
}

func TestDigitalPinConcurrency(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(oldProcs)

	translate := func(pin string) (string, int, error) { line, err := strconv.Atoi(pin); return "", line, err }
	sys := system.NewAccesser()

	for retry := 0; retry < 20; retry++ {

		a := NewDigitalPinsAdaptor(sys, translate)
		_ = a.Connect()
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
