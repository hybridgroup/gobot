package gpio

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot"
)

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

const (
	HD44780_2NDLINEOFFSET = 0x40
)

// data bus mode
type HD44780BusMode int

const (
	HD44780_4BITMODE HD44780BusMode = iota + 1
	HD44780_8BITMODE
)

// databit pins
type HD44780DataPin struct {
	D0 string // not used if 4bit mode
	D1 string // not used if 4bit mode
	D2 string // not used if 4bit mode
	D3 string // not used if 4bit mode
	D4 string
	D5 string
	D6 string
	D7 string
}

// HD44780Driver is the gobot driver for the HD44780 LCD controller
// Datasheet: https://www.sparkfun.com/datasheets/LCD/HD44780.pdf
type HD44780Driver struct {
	name        string
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
	connection  gobot.Connection
	gobot.Commander
	mutex *sync.Mutex // mutex is needed for sequences, like CreateChar(), Write(), Start()
}

// NewHD44780Driver return a new HD44780Driver
// a: gobot.Conenction
// cols: lcd columns
// rows: lcd rows
// busMode: 4Bit or 8Bit
// pinRS: register select pin
// pinEN: clock enable pin
// pinDataBits: databit pins
func NewHD44780Driver(a gobot.Connection, cols int, rows int, busMode HD44780BusMode, pinRS string, pinEN string, pinDataBits HD44780DataPin) *HD44780Driver {
	h := &HD44780Driver{
		name:       "HD44780Driver",
		cols:       cols,
		rows:       rows,
		busMode:    busMode,
		pinRS:      NewDirectPinDriver(a, pinRS),
		pinEN:      NewDirectPinDriver(a, pinEN),
		connection: a,
		Commander:  gobot.NewCommander(),
		mutex:      &sync.Mutex{},
	}

	if h.busMode == HD44780_4BITMODE {
		h.pinDataBits = make([]*DirectPinDriver, 4)
		h.pinDataBits[0] = NewDirectPinDriver(a, pinDataBits.D4)
		h.pinDataBits[1] = NewDirectPinDriver(a, pinDataBits.D5)
		h.pinDataBits[2] = NewDirectPinDriver(a, pinDataBits.D6)
		h.pinDataBits[3] = NewDirectPinDriver(a, pinDataBits.D7)
	} else {
		h.pinDataBits = make([]*DirectPinDriver, 8)
		h.pinDataBits[0] = NewDirectPinDriver(a, pinDataBits.D0)
		h.pinDataBits[1] = NewDirectPinDriver(a, pinDataBits.D1)
		h.pinDataBits[2] = NewDirectPinDriver(a, pinDataBits.D2)
		h.pinDataBits[3] = NewDirectPinDriver(a, pinDataBits.D3)
		h.pinDataBits[4] = NewDirectPinDriver(a, pinDataBits.D4)
		h.pinDataBits[5] = NewDirectPinDriver(a, pinDataBits.D5)
		h.pinDataBits[6] = NewDirectPinDriver(a, pinDataBits.D6)
		h.pinDataBits[7] = NewDirectPinDriver(a, pinDataBits.D7)
	}

	h.rowOffsets[0] = 0x00
	h.rowOffsets[1] = HD44780_2NDLINEOFFSET
	h.rowOffsets[2] = 0x00 + cols
	h.rowOffsets[3] = HD44780_2NDLINEOFFSET + cols

	/* TODO : Add commands */

	return h
}

// Halt implements the Driver interface
func (h *HD44780Driver) Halt() error { return nil }

// Name returns the HD44780Driver name
func (h *HD44780Driver) Name() string { return h.name }

// SetName sets the HD44780Driver name
func (h *HD44780Driver) SetName(n string) { h.name = n }

// Connecton returns the HD44780Driver Connection
func (h *HD44780Driver) Connection() gobot.Connection {
	return h.connection
}

// Start initializes the HD44780 LCD controller
// refer to page 45/46 of hitachi HD44780 datasheet
func (h *HD44780Driver) Start() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, bitPin := range h.pinDataBits {
		if bitPin.Pin() == "" {
			return errors.New("Initialization error")
		}
	}

	time.Sleep(50 * time.Millisecond)

	if err := h.activateWriteMode(); err != nil {
		return err
	}

	if h.busMode == HD44780_4BITMODE {
		if err := h.writeDataPins(0x03); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)

		if err := h.writeDataPins(0x03); err != nil {
			return err
		}
		time.Sleep(100 * time.Microsecond)

		if err := h.writeDataPins(0x03); err != nil {
			return err
		}
		time.Sleep(100 * time.Microsecond)

		if err := h.writeDataPins(0x02); err != nil {
			return err
		}
	} else {
		if err := h.sendCommand(0x30); err != nil {
			return err
		}
		time.Sleep(5 * time.Millisecond)

		if err := h.sendCommand(0x30); err != nil {
			return err
		}
		time.Sleep(100 * time.Microsecond)

		if err := h.sendCommand(0x30); err != nil {
			return err
		}
	}
	time.Sleep(100 * time.Microsecond)

	if h.busMode == HD44780_4BITMODE {
		h.displayFunc |= HD44780_4BITBUS
	} else {
		h.displayFunc |= HD44780_8BITBUS
	}

	if h.rows > 1 {
		h.displayFunc |= HD44780_2LINE
	} else {
		h.displayFunc |= HD44780_1LINE
	}

	h.displayFunc |= HD44780_5x8DOTS
	h.displayCtrl = HD44780_DISPLAYON | HD44780_BLINKOFF | HD44780_CURSOROFF
	h.displayMode = HD44780_ENTRYLEFT | HD44780_ENTRYSHIFTDECREMENT

	if err := h.sendCommand(HD44780_DISPLAYCONTROL | h.displayCtrl); err != nil {
		return err
	}
	if err := h.sendCommand(HD44780_FUNCTIONSET | h.displayFunc); err != nil {
		return err
	}
	if err := h.sendCommand(HD44780_ENTRYMODESET | h.displayMode); err != nil {
		return err
	}

	return h.clear()
}

// SetRWPin initializes the RW pin
func (h *HD44780Driver) SetRWPin(pinRW string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.pinRW = NewDirectPinDriver(h.connection, pinRW)
}

// Write output text to the display
func (h *HD44780Driver) Write(message string) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	col := 0
	if (h.displayMode & HD44780_ENTRYLEFT) == 0 {
		col = h.cols - 1
	}

	row := 0
	for _, c := range message {
		if c == '\n' {
			row++
			if err := h.setCursor(col, row); err != nil {
				return err
			}
			continue
		}
		if err := h.writeChar(int(c)); err != nil {
			return err
		}
	}

	return nil
}

// Clear clear the display
func (h *HD44780Driver) Clear() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.clear()
}

// Home return cursor to home
func (h *HD44780Driver) Home() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if err := h.sendCommand(HD44780_RETURNHOME); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)

	return nil
}

// SetCursor move the cursor to the specified position
func (h *HD44780Driver) SetCursor(col int, row int) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.setCursor(col, row)
}

// Display turn the display on and off
func (h *HD44780Driver) Display(on bool) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if on {
		h.displayCtrl |= HD44780_DISPLAYON
	} else {
		h.displayCtrl &= ^HD44780_DISPLAYON
	}

	return h.sendCommand(HD44780_DISPLAYCONTROL | h.displayCtrl)
}

// Cursor turn the cursor on and off
func (h *HD44780Driver) Cursor(on bool) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if on {
		h.displayCtrl |= HD44780_CURSORON
	} else {
		h.displayCtrl &= ^HD44780_CURSORON
	}

	return h.sendCommand(HD44780_DISPLAYCONTROL | h.displayCtrl)
}

// Blink turn the blink on and off
func (h *HD44780Driver) Blink(on bool) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if on {
		h.displayCtrl |= HD44780_BLINKON
	} else {
		h.displayCtrl &= ^HD44780_BLINKON
	}

	return h.sendCommand(HD44780_DISPLAYCONTROL | h.displayCtrl)
}

// ScrollLeft scroll text left
func (h *HD44780Driver) ScrollLeft() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.sendCommand(HD44780_CURSORSHIFT | HD44780_DISPLAYMOVE | HD44780_MOVELEFT)
}

// ScrollRight scroll text right
func (h *HD44780Driver) ScrollRight() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.sendCommand(HD44780_CURSORSHIFT | HD44780_DISPLAYMOVE | HD44780_MOVERIGHT)
}

// LeftToRight display text from left to right
func (h *HD44780Driver) LeftToRight() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.displayMode |= HD44780_ENTRYLEFT
	return h.sendCommand(HD44780_ENTRYMODESET | h.displayMode)
}

// RightToLeft display text from right to left
func (h *HD44780Driver) RightToLeft() (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.displayMode &= ^HD44780_ENTRYLEFT
	return h.sendCommand(HD44780_ENTRYMODESET | h.displayMode)
}

// SendCommand send control command
func (h *HD44780Driver) SendCommand(data int) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.sendCommand(data)
}

// WriteChar output a character to the display
func (h *HD44780Driver) WriteChar(data int) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.writeChar(data)
}

// CreateChar create custom character
func (h *HD44780Driver) CreateChar(pos int, charMap [8]byte) (err error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if pos > 7 {
		return errors.New("can't set a custom character at a position greater than 7")
	}

	if err := h.sendCommand(HD44780_SETCGRAMADDR | (pos << 3)); err != nil {
		return err
	}

	for i := range charMap {
		if err := h.writeChar(int(charMap[i])); err != nil {
			return err
		}
	}

	return nil
}

func (h *HD44780Driver) sendCommand(data int) (err error) {
	if err := h.activateWriteMode(); err != nil {
		return err
	}
	if err := h.pinRS.Off(); err != nil {
		return err
	}
	if h.busMode == HD44780_4BITMODE {
		if err := h.writeDataPins(data >> 4); err != nil {
			return err
		}
	}

	return h.writeDataPins(data)
}

func (h *HD44780Driver) writeChar(data int) (err error) {
	if err := h.activateWriteMode(); err != nil {
		return err
	}

	if err := h.pinRS.On(); err != nil {
		return err
	}
	if h.busMode == HD44780_4BITMODE {
		if err := h.writeDataPins(data >> 4); err != nil {
			return err
		}
	}

	return h.writeDataPins(data)
}

func (h *HD44780Driver) clear() (err error) {
	if err := h.sendCommand(HD44780_CLEARDISPLAY); err != nil {
		return err
	}
	time.Sleep(2 * time.Millisecond)

	return nil
}

func (h *HD44780Driver) setCursor(col int, row int) (err error) {
	if col < 0 || row < 0 || col >= h.cols || row >= h.rows {
		return fmt.Errorf("Invalid position value (%d, %d), range (%d, %d)", col, row, h.cols-1, h.rows-1)
	}

	return h.sendCommand(HD44780_SETDDRAMADDR | col + h.rowOffsets[row])
}

func (h *HD44780Driver) writeDataPins(data int) (err error) {
	for i, pin := range h.pinDataBits {
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
	return h.fallingEdge()
}

// fallingEdge creates falling edge to trigger data transmission
func (h *HD44780Driver) fallingEdge() (err error) {
	if err := h.pinEN.On(); err != nil {
		return err
	}
	time.Sleep(1 * time.Microsecond)

	if err := h.pinEN.Off(); err != nil {
		return err
	}
	// fastest write operation at 190kHz mode takes 53 us
	time.Sleep(60 * time.Microsecond)

	return nil
}

func (h *HD44780Driver) activateWriteMode() (err error) {
	if h.pinRW == nil {
		return
	}
	return h.pinRW.Off()
}
