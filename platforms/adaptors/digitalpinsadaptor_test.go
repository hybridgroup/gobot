package adaptors

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestDigitalPinsWithGpiosActiveLow(t *testing.T) {
	// This is a general test, that options are applied in constructor. Further tests for options
	// can also be done by call of "WithOption(val)(d)".
	// arrange
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()
	// act
	a := NewDigitalPinsAdaptor(sys, translate, WithGpiosActiveLow("1", "12", "33"))
	// assert
	assert.Equal(t, 3, len(a.pinOptions))
}

func TestDigitalPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()

	a := NewDigitalPinsAdaptor(sys, translate)
	assert.Equal(t, (map[string]gobot.DigitalPinner)(nil), a.pins)

	_, err := a.DigitalRead("13")
	assert.ErrorContains(t, err, "not connected for pin 13")

	err = a.DigitalWrite("7", 1)
	assert.ErrorContains(t, err, "not connected for pin 7")

	err = a.Connect()
	assert.NoError(t, err)
	assert.NotEqual(t, (map[string]gobot.DigitalPinner)(nil), a.pins)
	assert.Equal(t, 0, len(a.pins))
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
	assert.NoError(t, a.Finalize())
	// arrange
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.DigitalWrite("3", 1))
	assert.Equal(t, 1, len(a.pins))
	// act
	err := a.Finalize()
	// assert
	assert.NoError(t, err)
	assert.Equal(t, 0, len(a.pins))
	// assert that finalize after finalize is working
	assert.NoError(t, a.Finalize())
	// arrange missing sysfs file
	assert.NoError(t, a.Connect())
	assert.NoError(t, a.DigitalWrite("3", 2))
	delete(fs.Files, "/sys/class/gpio/unexport")
	err = a.Finalize()
	assert.Contains(t, err.Error(), "/sys/class/gpio/unexport: no such file")
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
	assert.NoError(t, a.DigitalWrite("4", 1))
	assert.Equal(t, 1, len(a.pins))
	assert.NoError(t, a.Finalize())
	// act
	err := a.Connect()
	// assert
	assert.NoError(t, err)
	assert.NotNil(t, a.pins)
	assert.Equal(t, 0, len(a.pins))
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
	assert.NoError(t, err)

	i, err := a.DigitalRead("14")
	assert.NoError(t, err)
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
	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	// assert error bubbling for read errors
	fs.WithReadError = true
	_, err = a.DigitalRead("13")
	assert.ErrorContains(t, err, "read error")

	// assert error bubbling for write errors
	fs.WithWriteError = true
	_, err = a.DigitalRead("7")
	assert.ErrorContains(t, err, "write error")
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
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, 1, got1) // there is no mechanism to negate mocked values
	assert.Equal(t, 0, got2)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio25/active_low"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio26/active_low"].Contents)
}

func TestDigitalWrite(t *testing.T) {
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
	assert.NoError(t, err)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio18/value"].Contents)

	// assert second write to same pin without error and just ignore unsupported options
	WithGpiosPullDown("7")(a)
	WithGpiosOpenSource("7")(a)
	WithGpioDebounce("7", 2*time.Second)(a)
	WithGpioEventOnRisingEdge("7", gpioEventHandler)(a)
	err = a.DigitalWrite("7", 1)
	assert.NoError(t, err)

	// assert error on bad id
	assert.ErrorContains(t, a.DigitalWrite("notexist", 1), "not a valid pin")

	// assert error bubbling
	fs.WithWriteError = true
	err = a.DigitalWrite("7", 0)
	assert.ErrorContains(t, err, "write error")
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
	assert.NoError(t, err)
	assert.Equal(t, "2", fs.Files["/sys/class/gpio/gpio19/value"].Contents)
	assert.Equal(t, "1", fs.Files["/sys/class/gpio/gpio19/active_low"].Contents)
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

func gpioEventHandler(o int, t time.Duration, et string, sn uint32, lsn uint32) {
	// the handler should never execute, because used in outputs and not supported by sysfs
	panic(fmt.Sprintf("event handler was called (%d, %d) unexpected for line %d with '%s' at %s!", sn, lsn, o, t, et))
}
