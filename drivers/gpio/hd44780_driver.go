package gpio

import (
	"errors"
	"fmt"
	"time"

	"gobot.io/x/gobot/v2"
)

// Commands for the driver
const (
	HD44780_CLEARDISPLAY        = 0x01
	HD44780_RETURNHOME          = 0x02
	HD44780_ENTRYMODESET        = 0x04
	HD44780_DISPLAYCONTROL      = 0x08
	HD44780_CURSORSHIFT         = 0x10
	HD44780_FUNCTIONSET         = 0x20
	HD44780_SETCGRAMADDR        = 0x40
	HD44780_SETDDRAMADDR        = 0x80
	HD44780_ENTRYRIGHT          = 0x00
	HD44780_ENTRYLEFT           = 0x02
	HD44780_ENTRYSHIFTINCREMENT = 0x01
	HD44780_ENTRYSHIFTDECREMENT = 0x00
	HD44780_DISPLAYON           = 0x04
	HD44780_DISPLAYOFF          = 0x00
	HD44780_CURSORON            = 0x02
	HD44780_CURSOROFF           = 0x00
	HD44780_BLINKON             = 0x01
	HD44780_BLINKOFF            = 0x00
	HD44780_DISPLAYMOVE         = 0x08
	HD44780_CURSORMOVE          = 0x00
	HD44780_MOVERIGHT           = 0x04
	HD44780_MOVELEFT            = 0x00
	HD44780_1LINE               = 0x00
	HD44780_2LINE               = 0x08
	HD44780_5x8DOTS             = 0x00
	HD44780_5x10DOTS            = 0x04
	HD44780_4BITBUS             = 0x00
	HD44780_8BITBUS             = 0x10
)

// Some useful constants for the driver
const (
	HD44780_2NDLINEOFFSET = 0x40
)

// HD44780BusMode is the data bus mode
type HD44780BusMode int

// Bus modes of the driver
const (
	HD44780_4BITMODE HD44780BusMode = iota + 1
	HD44780_8BITMODE
)

// HD44780DataPin are the data bit pins
type HD44780DataPin struct {
	D0 string // not used if 4Bit mode
	D1 string // not used if 4Bit mode
	D2 string // not used if 4Bit mode
	D3 string // not used if 4Bit mode
	D4 string
	D5 string
	D6 string
	D7 string
}

// hd44780OptionApplier needs to be implemented by each configurable option type
type hd44780OptionApplier interface {
	apply(cfg *hd44780Configuration)
}

// hd44780Configuration contains all changeable attributes of the driver
type hd44780Configuration struct {
	pinRW string
}

// hd44780PinRWOption is the type for applying a R/W pin to the configuration
type hd44780PinRWOption string

// HD44780Driver is the gobot driver for the HD44780 LCD controller
// Datasheet: https://www.sparkfun.com/datasheets/LCD/HD44780.pdf
type HD44780Driver struct {
	*driver
	hd44780Cfg  *hd44780Configuration
	cols        int
	rows        int
	rowOffsets  [4]int
	busMode     HD44780BusMode
	pinRS       *DirectPinDriver
	pinEN       *DirectPinDriver
	pinRW       *DirectPinDriver
	pinDataBits []*DirectPinDriver
	displayCtrl int
	displayFunc int
	displayMode int
}

// NewHD44780Driver return a new HD44780Driver
// a: gobot.Connection
// cols: lcd columns
// rows: lcd rows
// busMode: 4Bit or 8Bit
// pinRS: register select pin
// pinEN: clock enable pin
// pinDataBits: databit pins
//
// Supported options:
//
//	"WithName"
func NewHD44780Driver(
	a gobot.Connection,
	cols int,
	rows int,
	busMode HD44780BusMode,
	pinRS string,
	pinEN string,
	pinDataBits HD44780DataPin,
	opts ...interface{},
) *HD44780Driver {
	d := &HD44780Driver{
		driver:     newDriver(a, "HD44780"),
		hd44780Cfg: &hd44780Configuration{},
		cols:       cols,
		rows:       rows,
		busMode:    busMode,
		pinRS:      NewDirectPinDriver(a, pinRS),
		pinEN:      NewDirectPinDriver(a, pinEN),
	}
	d.afterStart = d.initialize

	for _, opt := range opts {
		switch o := opt.(type) {
		case optionApplier:
			o.apply(d.driverCfg)
		case hd44780OptionApplier:
			o.apply(d.hd44780Cfg)
		default:
			panic(fmt.Sprintf("'%s' can not be applied on '%s'", opt, d.driverCfg.name))
		}
	}

	if d.busMode == HD44780_4BITMODE {
		d.pinDataBits = make([]*DirectPinDriver, 4)
		d.pinDataBits[0] = NewDirectPinDriver(a, pinDataBits.D4)
		d.pinDataBits[1] = NewDirectPinDriver(a, pinDataBits.D5)
		d.pinDataBits[2] = NewDirectPinDriver(a, pinDataBits.D6)
		d.pinDataBits[3] = NewDirectPinDriver(a, pinDataBits.D7)
	} else {
		d.pinDataBits = make([]*DirectPinDriver, 8)
		d.pinDataBits[0] = NewDirectPinDriver(a, pinDataBits.D0)
		d.pinDataBits[1] = NewDirectPinDriver(a, pinDataBits.D1)
		d.pinDataBits[2] = NewDirectPinDriver(a, pinDataBits.D2)
		d.pinDataBits[3] = NewDirectPinDriver(a, pinDataBits.D3)
		d.pinDataBits[4] = NewDirectPinDriver(a, pinDataBits.D4)
		d.pinDataBits[5] = NewDirectPinDriver(a, pinDataBits.D5)
		d.pinDataBits[6] = NewDirectPinDriver(a, pinDataBits.D6)
		d.pinDataBits[7] = NewDirectPinDriver(a, pinDataBits.D7)
	}

	if d.hd44780Cfg.pinRW != "" {
		d.pinRW = NewDirectPinDriver(d.connection, d.hd44780Cfg.pinRW)
	}

	d.rowOffsets[0] = 0x00
	d.rowOffsets[1] = HD44780_2NDLINEOFFSET
	d.rowOffsets[2] = 0x00 + cols
	d.rowOffsets[3] = HD44780_2NDLINEOFFSET + cols

	/* TODO : Add commands */

	return d
}

// WithHD44780RWPin sets the RW pin for next initializing.
func WithHD44780RWPin(pin string) hd44780OptionApplier {
	return hd44780PinRWOption(pin)
}

// initialize initializes the HD44780 LCD controller
// refer to page 45/46 of Hitachi HD44780 datasheet
func (d *HD44780Driver) initialize() error {
	for _, bitPin := range d.pinDataBits {
		if bitPin.Pin() == "" {
			return errors.New("Initialization error")
		}
	}

	time.Sleep(50 * time.Millisecond)

	if err := d.activateWriteMode(); err != nil {
		return err
	}

	// for initialization refer to documentation, page 45 and 46
	if d.busMode == HD44780_4BITMODE {
		if err := d.writeDataPins(0x03); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)
		if err := d.writeDataPins(0x03); err != nil {
			return err
		}
		time.Sleep(100 * time.Microsecond)
		if err := d.writeDataPins(0x03); err != nil {
			return err
		}
		// no additional delay is necessary now
		if err := d.writeDataPins(0x02); err != nil {
			return err
		}
	} else {
		if err := d.sendCommand(0x30); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)
		if err := d.sendCommand(0x30); err != nil {
			return err
		}
		time.Sleep(100 * time.Microsecond)
		if err := d.sendCommand(0x30); err != nil {
			return err
		}
		// no additional delay is necessary now
	}

	if d.busMode == HD44780_4BITMODE {
		d.displayFunc |= HD44780_4BITBUS
	} else {
		d.displayFunc |= HD44780_8BITBUS
	}

	if d.rows > 1 {
		d.displayFunc |= HD44780_2LINE
	} else {
		d.displayFunc |= HD44780_1LINE
	}

	d.displayFunc |= HD44780_5x8DOTS
	d.displayCtrl = HD44780_DISPLAYON | HD44780_BLINKOFF | HD44780_CURSOROFF
	d.displayMode = HD44780_ENTRYLEFT | HD44780_ENTRYSHIFTDECREMENT

	if err := d.sendCommand(HD44780_FUNCTIONSET | d.displayFunc); err != nil {
		return err
	}

	if err := d.sendCommand(HD44780_DISPLAYCONTROL | d.displayCtrl); err != nil {
		return err
	}

	if err := d.clear(); err != nil {
		return err
	}

	if err := d.sendCommand(HD44780_ENTRYMODESET | d.displayMode); err != nil {
		return err
	}

	// see documentation, page 45, 46: the busy flag can't be checked before
	return nil
}

// Write output text to the display
func (d *HD44780Driver) Write(message string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	col := 0
	if (d.displayMode & HD44780_ENTRYLEFT) == 0 {
		col = d.cols - 1
	}

	row := 0
	for _, c := range message {
		if c == '\n' {
			row++
			if err := d.setCursor(col, row); err != nil {
				return err
			}
			continue
		}
		if err := d.writeChar(int(c)); err != nil {
			return err
		}
	}

	return nil
}

// Clear clear the display
func (d *HD44780Driver) Clear() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.clear()
}

// Home return cursor to home
func (d *HD44780Driver) Home() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if err := d.sendCommand(HD44780_RETURNHOME); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)

	return nil
}

// SetCursor move the cursor to the specified position
func (d *HD44780Driver) SetCursor(col int, row int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.setCursor(col, row)
}

// Display turn the display on and off
func (d *HD44780Driver) Display(on bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if on {
		d.displayCtrl |= HD44780_DISPLAYON
	} else {
		d.displayCtrl &= ^HD44780_DISPLAYON
	}

	return d.sendCommand(HD44780_DISPLAYCONTROL | d.displayCtrl)
}

// Cursor turn the cursor on and off
func (d *HD44780Driver) Cursor(on bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if on {
		d.displayCtrl |= HD44780_CURSORON
	} else {
		d.displayCtrl &= ^HD44780_CURSORON
	}

	return d.sendCommand(HD44780_DISPLAYCONTROL | d.displayCtrl)
}

// Blink turn the blink on and off
func (d *HD44780Driver) Blink(on bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if on {
		d.displayCtrl |= HD44780_BLINKON
	} else {
		d.displayCtrl &= ^HD44780_BLINKON
	}

	return d.sendCommand(HD44780_DISPLAYCONTROL | d.displayCtrl)
}

// ScrollLeft scroll text left
func (d *HD44780Driver) ScrollLeft() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.sendCommand(HD44780_CURSORSHIFT | HD44780_DISPLAYMOVE | HD44780_MOVELEFT)
}

// ScrollRight scroll text right
func (d *HD44780Driver) ScrollRight() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.sendCommand(HD44780_CURSORSHIFT | HD44780_DISPLAYMOVE | HD44780_MOVERIGHT)
}

// LeftToRight display text from left to right
func (d *HD44780Driver) LeftToRight() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.displayMode |= HD44780_ENTRYLEFT
	return d.sendCommand(HD44780_ENTRYMODESET | d.displayMode)
}

// RightToLeft display text from right to left
func (d *HD44780Driver) RightToLeft() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.displayMode &= ^HD44780_ENTRYLEFT
	return d.sendCommand(HD44780_ENTRYMODESET | d.displayMode)
}

// SendCommand send control command
func (d *HD44780Driver) SendCommand(data int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.sendCommand(data)
}

// WriteChar output a character to the display
func (d *HD44780Driver) WriteChar(data int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.writeChar(data)
}

// CreateChar create custom character
func (d *HD44780Driver) CreateChar(pos int, charMap [8]byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if pos > 7 {
		return errors.New("can't set a custom character at a position greater than 7")
	}

	if err := d.sendCommand(HD44780_SETCGRAMADDR | (pos << 3)); err != nil {
		return err
	}

	for i := range charMap {
		if err := d.writeChar(int(charMap[i])); err != nil {
			return err
		}
	}

	return nil
}

func (d *HD44780Driver) sendCommand(data int) error {
	if err := d.activateWriteMode(); err != nil {
		return err
	}
	if err := d.pinRS.Off(); err != nil {
		return err
	}
	if d.busMode == HD44780_4BITMODE {
		if err := d.writeDataPins(data >> 4); err != nil {
			return err
		}
	}

	return d.writeDataPins(data)
}

func (d *HD44780Driver) writeChar(data int) error {
	if err := d.activateWriteMode(); err != nil {
		return err
	}

	if err := d.pinRS.On(); err != nil {
		return err
	}
	if d.busMode == HD44780_4BITMODE {
		if err := d.writeDataPins(data >> 4); err != nil {
			return err
		}
	}

	return d.writeDataPins(data)
}

func (d *HD44780Driver) clear() error {
	if err := d.sendCommand(HD44780_CLEARDISPLAY); err != nil {
		return err
	}

	// clear is time consuming, see documentation for JHD1313
	// for lower clock speed it takes more time
	time.Sleep(4 * time.Millisecond)

	return nil
}

func (d *HD44780Driver) setCursor(col int, row int) error {
	if col < 0 || row < 0 || col >= d.cols || row >= d.rows {
		return fmt.Errorf("Invalid position value (%d, %d), range (%d, %d)", col, row, d.cols-1, d.rows-1)
	}

	return d.sendCommand(HD44780_SETDDRAMADDR | col + d.rowOffsets[row])
}

func (d *HD44780Driver) writeDataPins(data int) error {
	for i, pin := range d.pinDataBits {
		if ((data >> i) & 0x01) == 0x01 {
			if err := pin.On(); err != nil {
				return err
			}
		} else {
			if err := pin.Off(); err != nil {
				return err
			}
		}
	}
	return d.fallingEdge()
}

// fallingEdge creates falling edge to trigger data transmission
func (d *HD44780Driver) fallingEdge() error {
	if err := d.pinEN.On(); err != nil {
		return err
	}
	time.Sleep(1 * time.Microsecond)

	if err := d.pinEN.Off(); err != nil {
		return err
	}
	// fastest write operation at 190kHz mode takes 53 us
	time.Sleep(60 * time.Microsecond)

	return nil
}

func (d *HD44780Driver) activateWriteMode() error {
	if d.pinRW == nil {
		return nil
	}
	return d.pinRW.Off()
}

func (o hd44780PinRWOption) String() string {
	return "hd44780 RW pin option"
}

func (o hd44780PinRWOption) apply(cfg *hd44780Configuration) {
	cfg.pinRW = string(o)
}
