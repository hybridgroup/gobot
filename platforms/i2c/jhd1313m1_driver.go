package i2c

import (
	"time"

	"github.com/hybridgroup/gobot"
)

const (
	REG_RED   = 0x04
	REG_GREEN = 0x03
	REG_BLUE  = 0x02

	LCD_CLEARDISPLAY        = 0x01
	LCD_RETURNHOME          = 0x02
	LCD_ENTRYMODESET        = 0x04
	LCD_DISPLAYCONTROL      = 0x08
	LCD_CURSORSHIFT         = 0x10
	LCD_FUNCTIONSET         = 0x20
	LCD_SETCGRAMADDR        = 0x40
	LCD_SETDDRAMADDR        = 0x80
	LCD_ENTRYRIGHT          = 0x00
	LCD_ENTRYLEFT           = 0x02
	LCD_ENTRYSHIFTINCREMENT = 0x01
	LCD_ENTRYSHIFTDECREMENT = 0x00
	LCD_DISPLAYON           = 0x04
	LCD_DISPLAYOFF          = 0x00
	LCD_CURSORON            = 0x02
	LCD_CURSOROFF           = 0x00
	LCD_BLINKON             = 0x01
	LCD_BLINKOFF            = 0x00
	LCD_DISPLAYMOVE         = 0x08
	LCD_CURSORMOVE          = 0x00
	LCD_MOVERIGHT           = 0x04
	LCD_MOVELEFT            = 0x00
	LCD_2LINE               = 0x08
)

var _ gobot.Driver = (*JHD1313M1Driver)(nil)

// JHD1313M1Driver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type JHD1313M1Driver struct {
	name       string
	connection I2c
	lcdAddress int
	rgbAddress int
}

// NewJHD1313M1Driver creates a new driver with specified name and i2c interface.
func NewJHD1313M1Driver(a I2c, name string) *JHD1313M1Driver {
	return &JHD1313M1Driver{
		name:       name,
		connection: a,
		lcdAddress: 0x3E,
		rgbAddress: 0x62,
	}
}

// Name returns the name the JHD1313M1 Driver was given when created.
func (h *JHD1313M1Driver) Name() string { return h.name }

// Connection returns the driver connection to the device.
func (h *JHD1313M1Driver) Connection() gobot.Connection {
	return h.connection.(gobot.Connection)
}

// Start starts the backlit and the screen and initializes the states.
func (h *JHD1313M1Driver) Start() []error {
	cmd := uint8(0)
	if err := h.connection.I2cStart(h.lcdAddress); err != nil {
		return []error{err}
	}

	if err := h.connection.I2cStart(h.rgbAddress); err != nil {
		return []error{err}
	}

	cmd |= LCD_2LINE

	<-time.After(30 * time.Millisecond)

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_FUNCTIONSET | cmd}); err != nil {
		return []error{err}
	}
	<-time.After(40 * time.Nanosecond)

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_FUNCTIONSET | cmd}); err != nil {
		return []error{err}
	}
	<-time.After(150 * time.Microsecond)
	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_FUNCTIONSET | cmd}); err != nil {
		return []error{err}
	}
	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_FUNCTIONSET | cmd}); err != nil {
		return []error{err}
	}
	cmd |= LCD_DISPLAYON

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_DISPLAYCONTROL | cmd}); err != nil {
		return []error{err}
	}

	h.Clear()

	cmd |= LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT

	if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x80, LCD_ENTRYMODESET | cmd}); err != nil {
		return []error{err}
	}

	if err := h.setReg(0, 1); err != nil {
		return []error{err}
	}
	if err := h.setReg(1, 0); err != nil {
		return []error{err}
	}
	if err := h.setReg(0x08, 0xAA); err != nil {
		return []error{err}
	}

	if err := h.SetRGB(255, 255, 255); err != nil {
		return []error{err}
	}

	return nil
}

// SetRGB sets the Red Green Blue value of backlit.
func (h *JHD1313M1Driver) SetRGB(r, g, b int) error {
	if err := h.setReg(REG_RED, r); err != nil {
		return err
	}
	if err := h.setReg(REG_GREEN, g); err != nil {
		return err
	}
	return h.setReg(REG_BLUE, b)
}

// Clear clears the text on the lCD display.
func (h *JHD1313M1Driver) Clear() error {
	err := h.command([]byte{LCD_CLEARDISPLAY})
	<-time.After(2 * time.Millisecond)
	return err
}

// Home sets the cursor to the origin position on the display.
func (h *JHD1313M1Driver) Home() error {
	err := h.command([]byte{LCD_RETURNHOME})
	<-time.After(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (h *JHD1313M1Driver) Write(message string) error {
	for _, val := range message {
		if err := h.connection.I2cWrite(h.lcdAddress, []byte{0x40, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// Halt is a noop function.
func (h *JHD1313M1Driver) Halt() []error { return nil }

func (h *JHD1313M1Driver) setReg(command int, data int) error {
	if err := h.connection.I2cWrite(h.rgbAddress, []byte{byte(command), byte(data)}); err != nil {
		return err
	}
	return nil
}

func (h *JHD1313M1Driver) command(buf []byte) error {
	return h.connection.I2cWrite(h.lcdAddress, append([]byte{0x80}, buf...))
}
