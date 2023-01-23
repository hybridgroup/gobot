package adaptors

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"runtime"
	"strconv"
	"sync"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/gobottest"
	"gobot.io/x/gobot/system"
)

// make sure that this adaptor fulfills all the required interfaces
var _ gobot.DigitalPinnerProvider = (*DigitalPinsAdaptor)(nil)
var _ gpio.DigitalReader = (*DigitalPinsAdaptor)(nil)
var _ gpio.DigitalWriter = (*DigitalPinsAdaptor)(nil)

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
	gobottest.Assert(t, len(a.pinOptions), 3)
}

func TestDigitalPinsConnect(t *testing.T) {
	translate := func(pin string) (chip string, line int, err error) { return }
	sys := system.NewAccesser()

	a := NewDigitalPinsAdaptor(sys, translate)
	gobottest.Assert(t, a.pins, (map[string]gobot.DigitalPinner)(nil))

	_, err := a.DigitalRead("13")
	gobottest.Assert(t, err.Error(), "not connected")

	err = a.DigitalWrite("7", 1)
	gobottest.Assert(t, err.Error(), "not connected")

	err = a.Connect()
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, a.pins, (map[string]gobot.DigitalPinner)(nil))
	gobottest.Assert(t, len(a.pins), 0)
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
	gobottest.Assert(t, a.Finalize(), nil)
	// arrange
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("3", 1), nil)
	gobottest.Assert(t, len(a.pins), 1)
	// act
	err := a.Finalize()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, len(a.pins), 0)
	// assert that finalize after finalize is working
	gobottest.Assert(t, a.Finalize(), nil)
	// arrange missing sysfs file
	gobottest.Assert(t, a.Connect(), nil)
	gobottest.Assert(t, a.DigitalWrite("3", 2), nil)
	delete(fs.Files, "/sys/class/gpio/unexport")
	err = a.Finalize()
	gobottest.Assert(t, strings.Contains(err.Error(), "/sys/class/gpio/unexport: No such file"), true)
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
	gobottest.Assert(t, a.DigitalWrite("4", 1), nil)
	gobottest.Assert(t, len(a.pins), 1)
	gobottest.Assert(t, a.Finalize(), nil)
	// act
	err := a.Connect()
	// assert
	gobottest.Assert(t, err, nil)
	gobottest.Refute(t, a.pins, nil)
	gobottest.Assert(t, len(a.pins), 0)
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
	gobottest.Assert(t, err, nil)

	i, err := a.DigitalRead("14")
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 1)
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, i, 1)

	// assert error bubbling for read errors
	fs.WithReadError = true
	_, err = a.DigitalRead("13")
	gobottest.Assert(t, err, errors.New("read error"))

	// assert error bubbling for write errors
	fs.WithWriteError = true
	_, err = a.DigitalRead("7")
	gobottest.Assert(t, err, errors.New("write error"))
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
	gobottest.Assert(t, err1, nil)
	gobottest.Assert(t, err2, nil)
	gobottest.Assert(t, got1, 1) // there is no mechanism to negate mocked values
	gobottest.Assert(t, got2, 0)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio25/active_low"].Contents, "1")
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio26/active_low"].Contents, "1")
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
	err := a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio18/value"].Contents, "1")

	// assert second write to same pin without error and just ignore unsupported options
	WithGpiosPullDown("7")(a)
	WithGpiosOpenSource("7")(a)
	WithGpioDebounce("7", 2*time.Second)(a)
	err = a.DigitalWrite("7", 1)
	gobottest.Assert(t, err, nil)

	// assert error on bad id
	gobottest.Assert(t, a.DigitalWrite("notexist", 1), errors.New("not a valid pin"))

	// assert error bubbling
	fs.WithWriteError = true
	err = a.DigitalWrite("7", 0)
	gobottest.Assert(t, err, errors.New("write error"))
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
	gobottest.Assert(t, err, nil)
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio19/value"].Contents, "2")
	gobottest.Assert(t, fs.Files["/sys/class/gpio/gpio19/active_low"].Contents, "1")
}

func TestDigitalPinConcurrency(t *testing.T) {
	oldProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(oldProcs)

	translate := func(pin string) (string, int, error) { line, err := strconv.Atoi(pin); return "", line, err }
	sys := system.NewAccesser()

	for retry := 0; retry < 20; retry++ {

		a := NewDigitalPinsAdaptor(sys, translate)
		a.Connect()
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			pinAsString := strconv.Itoa(i)
			go func(pin string) {
				defer wg.Done()
				a.DigitalPin(pin)
			}(pinAsString)
		}

		wg.Wait()
	}
}
