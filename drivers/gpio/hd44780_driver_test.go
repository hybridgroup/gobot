package gpio

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*HD44780Driver)(nil)

func initTestHD44780Driver() *HD44780Driver {
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	if err := d.Start(); err != nil {
		panic(err)
	}
	return d
}

func initTestHD44780Driver4BitModeWithStubbedAdaptor() (*HD44780Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	dataPins := HD44780DataPin{
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "12",
	}

	return NewHD44780Driver(a, 2, 16, HD44780_4BITMODE, "13", "15", dataPins), a
}

func initTestHD44780Driver8BitModeWithStubbedAdaptor() (*HD44780Driver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
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

	return NewHD44780Driver(a, 2, 16, HD44780_8BITMODE, "13", "15", dataPins), a
}

func TestNewHD44780Driver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	dataPins := HD44780DataPin{
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "12",
	}
	// act
	d := NewHD44780Driver(a, 16, 2, HD44780_4BITMODE, "13", "15", dataPins)
	// assert
	assert.IsType(t, &HD44780Driver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "HD44780"))
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.Equal(t, "", d.hd44780Cfg.pinRW)
	assert.Equal(t, 16, d.cols)
	assert.Equal(t, 2, d.rows)
	assert.Equal(t, HD44780_4BITMODE, d.busMode)
	assert.NotNil(t, d.pinRS)
	assert.NotNil(t, d.pinEN)
	assert.Nil(t, d.pinRW) // will be set optionally
	assert.NotNil(t, d.pinRS)
	assert.Equal(t, [4]int{0, 64, 16, 80}, d.rowOffsets)
	assert.Len(t, d.pinDataBits, 4)
	for _, b := range d.pinDataBits {
		assert.NotNil(t, b)
	}
	assert.Equal(t, 0, d.displayCtrl)
	assert.Equal(t, 0, d.displayFunc)
	assert.Equal(t, 0, d.displayMode)
}

func TestNewHD44780Driver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "LCD output"
		pinRW  = "3"
	)
	dataPins := HD44780DataPin{
		D4: "22",
		D5: "18",
		D6: "16",
		D7: "12",
	}
	panicFunc := func() {
		NewHD44780Driver(newGpioTestAdaptor(), 16, 2, HD44780_4BITMODE, "1", "2", dataPins, WithName("crazy"),
			aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewHD44780Driver(newGpioTestAdaptor(), 16, 2, HD44780_4BITMODE, "1", "2", dataPins, WithName(myName),
		WithHD44780RWPin(pinRW))
	// assert
	assert.Equal(t, pinRW, d.hd44780Cfg.pinRW)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestHD44780Start(t *testing.T) {
	// arrange
	d, _ := initTestHD44780Driver4BitModeWithStubbedAdaptor()
	// act & assert: tests also initialize()
	require.NoError(t, d.Start())
}

func TestHD44780StartError(t *testing.T) {
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
	require.EqualError(t, d.Start(), "Initialization error")

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
	require.EqualError(t, d.Start(), "Initialization error")
}

func TestHD44780Write(t *testing.T) {
	var d *HD44780Driver

	d, _ = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Write("hello gobot"))

	d, _ = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	_ = d.Start()
	require.NoError(t, d.Write("hello gobot"))
}

func TestHD44780WriteError(t *testing.T) {
	var d *HD44780Driver
	var a *gpioTestAdaptor

	d, a = initTestHD44780Driver4BitModeWithStubbedAdaptor()
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	_ = d.Start()
	require.EqualError(t, d.Write("hello gobot"), "write error")

	d, a = initTestHD44780Driver8BitModeWithStubbedAdaptor()
	a.digitalWriteFunc = func(string, byte) error {
		return errors.New("write error")
	}
	_ = d.Start()
	require.EqualError(t, d.Write("hello gobot"), "write error")
}

func TestHD44780Clear(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Clear())
}

func TestHD44780Home(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Home())
}

func TestHD44780SetCursor(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.SetCursor(0, 3))
}

func TestHD44780SetCursorInvalid(t *testing.T) {
	d := initTestHD44780Driver()

	require.EqualError(t, d.SetCursor(-1, 3), "Invalid position value (-1, 3), range (1, 15)")
	require.EqualError(t, d.SetCursor(2, 3), "Invalid position value (2, 3), range (1, 15)")
	require.EqualError(t, d.SetCursor(0, -1), "Invalid position value (0, -1), range (1, 15)")
	require.EqualError(t, d.SetCursor(0, 16), "Invalid position value (0, 16), range (1, 15)")
}

func TestHD44780DisplayOn(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Display(true))
}

func TestHD44780DisplayOff(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Display(false))
}

func TestHD44780CursorOn(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Cursor(true))
}

func TestHD44780CursorOff(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Cursor(false))
}

func TestHD44780BlinkOn(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Blink(true))
}

func TestHD44780BlinkOff(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.Blink(false))
}

func TestHD44780ScrollLeft(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.ScrollLeft())
}

func TestHD44780ScrollRight(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.ScrollRight())
}

func TestHD44780LeftToRight(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.LeftToRight())
}

func TestHD44780RightToLeft(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.RightToLeft())
}

func TestHD44780SendCommand(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.SendCommand(0x33))
}

func TestHD44780WriteChar(t *testing.T) {
	d := initTestHD44780Driver()
	require.NoError(t, d.WriteChar(0x41))
}

func TestHD44780CreateChar(t *testing.T) {
	d := initTestHD44780Driver()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	require.NoError(t, d.CreateChar(0, charMap))
}

func TestHD44780CreateCharError(t *testing.T) {
	d := initTestHD44780Driver()
	charMap := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	require.EqualError(t, d.CreateChar(8, charMap), "can't set a custom character at a position greater than 7")
}
