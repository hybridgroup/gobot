package i2c

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2"
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
	LCD_CMD                 = 0x80
	LCD_DATA                = 0x40

	LCD_2NDLINEOFFSET = 0x40
)

// CustomLCDChars is a map of CGRAM characters that can be loaded
// into a LCD screen to display custom characters. Some LCD screens such
// as the Grove screen (jhd1313m1) isn't loaded with latin 1 characters.
// It's up to the developer to load the set up to 8 custom characters and
// update the input text so the character is swapped by a byte reflecting
// the position of the custom character to use.
// See SetCustomChar
var CustomLCDChars = map[string][8]byte{
	"é":       {130, 132, 142, 145, 159, 144, 142, 128},
	"è":       {136, 132, 142, 145, 159, 144, 142, 128},
	"ê":       {132, 138, 142, 145, 159, 144, 142, 128},
	"à":       {136, 134, 128, 142, 145, 147, 141, 128},
	"â":       {132, 138, 128, 142, 145, 147, 141, 128},
	"á":       {2, 4, 14, 1, 15, 17, 15, 0},
	"î":       {132, 138, 128, 140, 132, 132, 142, 128},
	"í":       {2, 4, 12, 4, 4, 4, 14, 0},
	"û":       {132, 138, 128, 145, 145, 147, 141, 128},
	"ù":       {136, 134, 128, 145, 145, 147, 141, 128},
	"ñ":       {14, 0, 22, 25, 17, 17, 17, 0},
	"ó":       {2, 4, 14, 17, 17, 17, 14, 0},
	"heart":   {0, 10, 31, 31, 31, 14, 4, 0},
	"smiley":  {0, 0, 10, 0, 0, 17, 14, 0},
	"frowney": {0, 0, 10, 0, 0, 0, 14, 17},
}

var jhd1313m1ErrInvalidPosition = fmt.Errorf("Invalid position value")

// JHD1313M1Driver is a driver for the Jhd1313m1 LCD display which has two i2c addreses,
// one belongs to a controller and the other controls solely the backlight.
// This module was tested with the Seed Grove LCD RGB Backlight v2.0 display which requires 5V to operate.
// http://www.seeedstudio.com/wiki/Grove_-_LCD_RGB_Backlight
type JHD1313M1Driver struct {
	name      string
	connector Connector
	Config
	gobot.Commander
	lcdAddress    int
	lcdConnection Connection
	rgbAddress    int
	rgbConnection Connection
}

// NewJHD1313M1Driver creates a new driver with specified i2c interface.
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
func NewJHD1313M1Driver(a Connector, options ...func(Config)) *JHD1313M1Driver {
	d := &JHD1313M1Driver{
		name:       gobot.DefaultName("JHD1313M1"),
		connector:  a,
		Config:     NewConfig(),
		Commander:  gobot.NewCommander(),
		lcdAddress: 0x3E,
		rgbAddress: 0x62,
	}

	for _, option := range options {
		option(d)
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r, _ := strconv.Atoi(params["r"].(string))
		g, _ := strconv.Atoi(params["g"].(string))
		b, _ := strconv.Atoi(params["b"].(string))
		return d.SetRGB(r, g, b)
	})
	d.AddCommand("Clear", func(_ map[string]interface{}) interface{} {
		return d.Clear()
	})
	d.AddCommand("Home", func(_ map[string]interface{}) interface{} {
		return d.Home()
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("Write", func(params map[string]interface{}) interface{} {
		msg := params["msg"].(string)
		return d.Write(msg)
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("SetPosition", func(params map[string]interface{}) interface{} {
		pos, _ := strconv.Atoi(params["pos"].(string))
		return d.SetPosition(pos)
	})
	//nolint:forcetypeassert // ok here
	d.AddCommand("Scroll", func(params map[string]interface{}) interface{} {
		lr, _ := strconv.ParseBool(params["lr"].(string))
		return d.Scroll(lr)
	})

	return d
}

// Name returns the name the JHD1313M1 Driver was given when created.
func (d *JHD1313M1Driver) Name() string { return d.name }

// SetName sets the name for the JHD1313M1 Driver.
func (d *JHD1313M1Driver) SetName(n string) { d.name = n }

// Connection returns the driver connection to the device.
func (d *JHD1313M1Driver) Connection() gobot.Connection {
	if conn, ok := d.connector.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

// Start starts the backlit and the screen and initializes the states.
func (d *JHD1313M1Driver) Start() error {
	bus := d.GetBusOrDefault(d.connector.DefaultI2cBus())

	var err error
	if d.lcdConnection, err = d.connector.GetI2cConnection(d.lcdAddress, bus); err != nil {
		return err
	}

	if d.rgbConnection, err = d.connector.GetI2cConnection(d.rgbAddress, bus); err != nil {
		return err
	}

	// SEE PAGE 45/46 FOR INITIALIZATION SPECIFICATION!
	// according to datasheet, we need at least 40ms after power rises above 2.7V
	// before sending commands. Arduino can turn on way befer 4.5V so we'll wait 50
	time.Sleep(50 * time.Millisecond)

	// this is according to the hitachi HD44780 datasheet
	// page 45 figure 23
	// Send function set command sequence
	init_payload := []byte{LCD_CMD, LCD_FUNCTIONSET | LCD_2LINE}
	if _, err := d.lcdConnection.Write(init_payload); err != nil {
		return err
	}

	// wait more than 4.1ms
	time.Sleep(4500 * time.Microsecond)
	// second try
	if _, err := d.lcdConnection.Write(init_payload); err != nil {
		return err
	}

	time.Sleep(150 * time.Microsecond)
	// third go
	if _, err := d.lcdConnection.Write(init_payload); err != nil {
		return err
	}

	if _, err := d.lcdConnection.Write([]byte{LCD_CMD, LCD_DISPLAYCONTROL | LCD_DISPLAYON}); err != nil {
		return err
	}

	time.Sleep(100 * time.Microsecond)
	if err := d.Clear(); err != nil {
		return err
	}

	if _, err := d.lcdConnection.Write([]byte{
		LCD_CMD, LCD_ENTRYMODESET | LCD_ENTRYLEFT | LCD_ENTRYSHIFTDECREMENT,
	}); err != nil {
		return err
	}

	if err := d.setReg(0, 0); err != nil {
		return err
	}
	if err := d.setReg(1, 0); err != nil {
		return err
	}
	if err := d.setReg(0x08, 0xAA); err != nil {
		return err
	}

	return d.SetRGB(255, 255, 255)
}

// SetRGB sets the Red Green Blue value of backlit.
func (d *JHD1313M1Driver) SetRGB(r, g, b int) error {
	if err := d.setReg(REG_RED, r); err != nil {
		return err
	}
	if err := d.setReg(REG_GREEN, g); err != nil {
		return err
	}
	return d.setReg(REG_BLUE, b)
}

// Clear clears the text on the lCD display.
func (d *JHD1313M1Driver) Clear() error {
	return d.command([]byte{LCD_CLEARDISPLAY})
}

// Home sets the cursor to the origin position on the display.
func (d *JHD1313M1Driver) Home() error {
	err := d.command([]byte{LCD_RETURNHOME})
	// This wait fixes a race condition when calling home and clear back to back.
	time.Sleep(2 * time.Millisecond)
	return err
}

// Write displays the passed message on the screen.
func (d *JHD1313M1Driver) Write(message string) error {
	// This wait fixes an odd bug where the clear function doesn't always work properly.
	time.Sleep(1 * time.Millisecond)
	for _, val := range message {
		if val == '\n' {
			if err := d.SetPosition(16); err != nil {
				return err
			}
			continue
		}
		if _, err := d.lcdConnection.Write([]byte{LCD_DATA, byte(val)}); err != nil {
			return err
		}
	}
	return nil
}

// SetPosition sets the cursor and the data display to pos.
// 0..15 are the positions in the first display line.
// 16..32 are the positions in the second display line.
func (d *JHD1313M1Driver) SetPosition(pos int) error {
	if pos < 0 || pos > 31 {
		return jhd1313m1ErrInvalidPosition
	}
	offset := byte(pos)
	if pos >= 16 {
		offset -= 16
		offset |= LCD_2NDLINEOFFSET
	}

	return d.command([]byte{LCD_SETDDRAMADDR | offset})
}

// Scroll sets the scrolling direction for the display, either left to right, or
// right to left.
func (d *JHD1313M1Driver) Scroll(lr bool) error {
	if lr {
		_, err := d.lcdConnection.Write([]byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVELEFT})
		return err
	}

	_, err := d.lcdConnection.Write([]byte{LCD_CMD, LCD_CURSORSHIFT | LCD_DISPLAYMOVE | LCD_MOVERIGHT})
	return err
}

// Halt is a noop function.
func (d *JHD1313M1Driver) Halt() error { return nil }

// SetCustomChar sets one of the 8 CGRAM locations with a custom character.
// The custom character can be used by writing a byte of value 0 to 7.
// When you are using LCD as 5x8 dots in function set then you can define a total of 8 user defined patterns
// (1 Byte for each row and 8 rows for each pattern).
// Use http://www.8051projects.net/lcd-interfacing/lcd-custom-character.php to create your own
// characters.
// To use a custom character, write byte value of the custom character position as a string after
// having setup the custom character.
func (d *JHD1313M1Driver) SetCustomChar(pos int, charMap [8]byte) error {
	if pos > 7 {
		return fmt.Errorf("can't set a custom character at a position greater than 7")
	}
	location := uint8(pos) //nolint:gosec // checked before
	if err := d.command([]byte{LCD_SETCGRAMADDR | (location << 3)}); err != nil {
		return err
	}
	_, err := d.lcdConnection.Write(append([]byte{LCD_DATA}, charMap[:]...))
	return err
}

func (d *JHD1313M1Driver) setReg(command int, data int) error {
	_, err := d.rgbConnection.Write([]byte{byte(command), byte(data)})
	return err
}

func (d *JHD1313M1Driver) command(buf []byte) error {
	_, err := d.lcdConnection.Write(append([]byte{LCD_CMD}, buf...))
	return err
}
