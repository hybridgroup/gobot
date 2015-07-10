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

	LCD_2NDLINEOFFSET = 0x40
)

var _ gobot.Driver = (*JHD1313M1Driver)(nil)

type JHD1313M1Driver struct {
	name       string
	connection I2c
	lcdAddress int
	rgbAddress int
}

// NewJHD1313M1Driver creates a new driver with specified name and i2c interface
func NewJHD1313M1Driver(a I2c, name string) *JHD1313M1Driver {
	return &JHD1313M1Driver{
		name:       name,
		connection: a,
		lcdAddress: 0x3E,
		rgbAddress: 0x62,
	}
}

func (h *JHD1313M1Driver) Name() string                 { return h.name }
func (h *JHD1313M1Driver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

func (h *JHD1313M1Driver) Start() (errs []error) {
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

	return
}

func (h *JHD1313M1Driver) SetRGB(r, g, b int) (err error) {
	if err = h.setReg(REG_RED, r); err != nil {
		return
	}
	if err = h.setReg(REG_GREEN, g); err != nil {
		return
	}
	return h.setReg(REG_BLUE, b)
}

func (h *JHD1313M1Driver) setReg(command int, data int) (err error) {
	if err = h.connection.I2cWrite(h.rgbAddress, []byte{byte(command), byte(data)}); err != nil {
		return
	}
	return
}

func (h *JHD1313M1Driver) Clear() (err error) {
	return h.command([]byte{LCD_CLEARDISPLAY})
}

func (h *JHD1313M1Driver) Home() (err error) {
	return h.command([]byte{LCD_RETURNHOME})
}

func (h *JHD1313M1Driver) Write(message string) (err error) {
	for _, val := range message {
		if val == '\n' {
			if err = h.SetPosition(16); err != nil {
				return
			}
			continue
		}
		if err = h.connection.I2cWrite(h.lcdAddress, []byte{0x40, byte(val)}); err != nil {
			break
		}
	}
	return
}

// SetPosition sets the cursor and the data display to pos.
// 0..15 are the positions in the first display line.
// 16..32 are the positions in the second display line.
func (h *JHD1313M1Driver) SetPosition(pos int) (err error) {
	if pos < 0 || pos > 31 {
		err = ErrInvalidPosition
		return
	}
	offset := byte(pos)
	if pos >= 16 {
		offset -= 16
		offset |= LCD_2NDLINEOFFSET
	}
	err = h.command([]byte{LCD_SETDDRAMADDR | offset})
	return
}

func (h *JHD1313M1Driver) Halt() (errs []error) { return }

func (h *JHD1313M1Driver) command(buf []byte) (err error) {
	return h.connection.I2cWrite(h.lcdAddress, append([]byte{0x80}, buf...))
}
