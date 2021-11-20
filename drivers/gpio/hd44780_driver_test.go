package gpio

import (
	"errors"
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*HD44780Driver)(nil)

// --------- HELPERS
func initTestHD44780Driver() (driver *HD44780Driver) {
	driver, _ = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	return
}

func initTestHD44780Driver4BitModeWithStubbedAdaptor() (*HD44780Driver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	dataPins := HD44780DataPin{
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "12",
	}

	return NewHD44780Driver(adaptor, 2, 16, HD44780_4BITMODE, "13", "15", dataPins), adaptor
}

func initTestHD44780Driver8BitModeWithStubbedAdaptor() (*HD44780Driver, *gpioTestAdaptor) {
	adaptor := newGpioTestAdaptor()
	dataPins := HD44780DataPin{
		D0: "31",
		D1: "33",
		D2: "35",
		D3: "37",
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "12",
	}

	return NewHD44780Driver(adaptor, 2, 16, HD44780_8BITMODE, "13", "15", dataPins), adaptor
}

// --------- TESTS
func TestHD44780Driver(t *testing.T) {
	var a interface{} = initTestHD44780Driver()
	_, ok := a.(*HD44780Driver)
	if !ok {
		t.Errorf("NewHD44780Driver() should have returned a *HD44780Driver")
	}
}

func TestHD44780DriverHalt(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Halt(), nil)
}

func TestHD44780DriverDefaultName(t *testing.T) {
	d := initTestHD44780Driver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "HD44780Driver"), true)
}

func TestHD44780DriverSetName(t *testing.T) {
	d := initTestHD44780Driver()
	d.SetName("my driver")
	gobottest.Assert(t, d.Name(), "my driver")
}

func TestHD44780DriverStart(t *testing.T) {
	d := initTestHD44780Driver()
	gobottest.Assert(t, d.Start(), nil)
}

func TestHD44780DriverStartError(t *testing.T) {
	a := newGpioTestAdaptor()

	var pins HD44780DataPin
	var d *HD44780Driver

	pins = HD44780DataPin{
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "",
	}
	d = NewHD44780Driver(a, 2, 16, HD44780_4BITMODE, "13", "15", pins)
	gobottest.Assert(t, d.Start(), errors.New("Initialization error"))

	pins = HD44780DataPin{
		D0: "31",
		D1: "33",
		D2: "35",
		D3: "37",
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "",
	}
	d = NewHD44780Driver(a, 2, 16, HD44780_8BITMODE, "13", "15", pins)
	gobottest.Assert(t, d.Start(), errors.New("Initialization error"))
}

func TestHD44780DriverWrite(t *testing.T) {
	var d *HD44780Driver

	d, _ = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("hello gobot"), nil)

	d, _ = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	d.Start()
	gobottest.Assert(t, d.Write("hello gobot"), nil)
}

func TestHD44780DriverWriteError(t *testing.T) {
	var d *HD44780Driver
	var a *gpioTestAdaptor

	d, a = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	a.testAdaptorDigitalWrite = func(string, byte) (err error) {
		return errors.New("write error")
	}
	d.Start()
	gobottest.Assert(t, d.Write("hello gobot"), errors.New("write error"))

	d, a = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	a.testAdaptorDigitalWrite = func(string, byte) (err error) {
		return errors.New("write error")
	}
	d.Start()
	gobottest.Assert(t, d.Write("hello gobot"), errors.New("write error"))
}

func TestHD44780DriverClear(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Clear(), nil)
}

func TestHD44780DriverHome(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Home(), nil)
}

func TestHD44780DriverSetCursor(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.SetCursor(0, 3), nil)
}

func TestHD44780DriverSetCursorInvalid(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.SetCursor(-1, 3), errors.New("Invalid position value"))
	gobottest.Assert(t, d.SetCursor(2, 3), errors.New("Invalid position value"))
	gobottest.Assert(t, d.SetCursor(0, -1), errors.New("Invalid position value"))
	gobottest.Assert(t, d.SetCursor(0, 16), errors.New("Invalid position value"))
}

func TestHD44780DriverDisplayOn(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Display(true), nil)
}

func TestHD44780DriverDisplayOff(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Display(false), nil)
}

func TestHD44780DriverCursorOn(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Cursor(true), nil)
}

func TestHD44780DriverCursorOff(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Cursor(false), nil)
}

func TestHD44780DriverBlinkOn(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Blink(true), nil)
}

func TestHD44780DriverBlinkOff(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.Blink(false), nil)
}

func TestHD44780DriverScrollLeft(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.ScrollLeft(), nil)
}

func TestHD44780DriverScrollRight(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.ScrollRight(), nil)
}

func TestHD44780DriverLeftToRight(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.LeftToRight(), nil)
}

func TestHD44780DriverRightToLeft(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.RightToLeft(), nil)
}

func TestHD44780DriverSendCommand(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.SendCommand(0x33), nil)
}

func TestHD44780DriverWriteChar(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	gobottest.Assert(t, d.WriteChar(0x41), nil)
}

func TestHD44780DriverCreateChar(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	gobottest.Assert(t, d.CreateChar(0, charMap), nil)
}

func TestHD44780DriverCreateCharError(t *testing.T) {
	d := initTestHD44780Driver()
	d.Start()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	gobottest.Assert(t, d.CreateChar(8, charMap), errors.New("can't set a custom character at a position greater than 7"))
}
