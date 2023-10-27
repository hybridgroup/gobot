package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*HD44780Driver)(nil)

// --------- HELPERS
func initTestHD44780Driver() *HD44780Driver {
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d
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
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	var a interface{} = d
	_, ok := a.(*HD44780Driver)
	if !ok {
		t.Errorf("NewHD44780Driver() should have returned a *HD44780Driver")
	}
}

func TestHD44780DriverHalt(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Halt())
}

func TestHD44780DriverDefaultName(t *testing.T) {
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	assert.True(t, strings.HasPrefix(d.Name(), "HD44780Driver"))
}

func TestHD44780DriverSetName(t *testing.T) {
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	d.SetName("my driver")
	assert.Equal(t, "my driver", d.Name())
}

func TestHD44780DriverStart(t *testing.T) {
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	assert.NoError(t, d.Start())
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
	assert.ErrorContains(t, d.Start(), "Initialization error")

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
	assert.ErrorContains(t, d.Start(), "Initialization error")
}

func TestHD44780DriverWrite(t *testing.T) {
	var d *HD44780Driver

	d, _ = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	_ = d.Start()
	assert.NoError(t, d.Write("hello gobot"))

	d, _ = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	_ = d.Start()
	assert.NoError(t, d.Write("hello gobot"))
}

func TestHD44780DriverWriteError(t *testing.T) {
	var d *HD44780Driver
	var a *gpioTestAdaptor

	d, a = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}
	_ = d.Start()
	assert.ErrorContains(t, d.Write("hello gobot"), "write error")

	d, a = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	a.digitalWriteFunc = func(string, byte) (err error) {
		return errors.New("write error")
	}
	_ = d.Start()
	assert.ErrorContains(t, d.Write("hello gobot"), "write error")
}

func TestHD44780DriverClear(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Clear())
}

func TestHD44780DriverHome(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Home())
}

func TestHD44780DriverSetCursor(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.SetCursor(0, 3))
}

func TestHD44780DriverSetCursorInvalid(t *testing.T) {
	d := initTestHD44780Driver()

	assert.ErrorContains(t, d.SetCursor(-1, 3), "Invalid position value (-1, 3), range (1, 15)")
	assert.ErrorContains(t, d.SetCursor(2, 3), "Invalid position value (2, 3), range (1, 15)")
	assert.ErrorContains(t, d.SetCursor(0, -1), "Invalid position value (0, -1), range (1, 15)")
	assert.ErrorContains(t, d.SetCursor(0, 16), "Invalid position value (0, 16), range (1, 15)")
}

func TestHD44780DriverDisplayOn(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Display(true))
}

func TestHD44780DriverDisplayOff(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Display(false))
}

func TestHD44780DriverCursorOn(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Cursor(true))
}

func TestHD44780DriverCursorOff(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Cursor(false))
}

func TestHD44780DriverBlinkOn(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Blink(true))
}

func TestHD44780DriverBlinkOff(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.Blink(false))
}

func TestHD44780DriverScrollLeft(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.ScrollLeft())
}

func TestHD44780DriverScrollRight(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.ScrollRight())
}

func TestHD44780DriverLeftToRight(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.LeftToRight())
}

func TestHD44780DriverRightToLeft(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.RightToLeft())
}

func TestHD44780DriverSendCommand(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.SendCommand(0x33))
}

func TestHD44780DriverWriteChar(t *testing.T) {
	d := initTestHD44780Driver()
	assert.NoError(t, d.WriteChar(0x41))
}

func TestHD44780DriverCreateChar(t *testing.T) {
	d := initTestHD44780Driver()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	assert.NoError(t, d.CreateChar(0, charMap))
}

func TestHD44780DriverCreateCharError(t *testing.T) {
	d := initTestHD44780Driver()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	assert.ErrorContains(t, d.CreateChar(8, charMap), "can't set a custom character at a position greater than 7")
}
