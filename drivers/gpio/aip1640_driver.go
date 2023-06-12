package gpio

import (
	"time"

	"gobot.io/x/gobot/v2"
)

// Commands of the driver
const (
	AIP1640DataCmd  = 0x40
	AIP1640DispCtrl = 0x88
	AIP1640AddrCmd  = 0xC0

	AIP1640FixedAddr = 0x04
)

// AIP1640Driver is the gobot driver for the AIP1640 LED driver used in the WEMOS D1 mini Matrix LED Shield.
// It has some similarities with the TM16xx LED drivers
//
// Datasheet CN: https://datasheet.lcsc.com/szlcsc/AiP1640_C82650.pdf
//
// Library ported from: https://github.com/wemos/WEMOS_Matrix_LED_Shield_Arduino_Library
type AIP1640Driver struct {
	pinClock   *DirectPinDriver
	pinData    *DirectPinDriver
	name       string
	intensity  byte
	buffer     [8]byte
	connection gobot.Connection
	gobot.Commander
}

// NewAIP1640Driver return a new AIP1640Driver given a gobot.Connection and the clock, data and strobe pins
func NewAIP1640Driver(a gobot.Connection, clockPin string, dataPin string) *AIP1640Driver {
	t := &AIP1640Driver{
		name:       gobot.DefaultName("AIP1640Driver"),
		pinClock:   NewDirectPinDriver(a, clockPin),
		pinData:    NewDirectPinDriver(a, dataPin),
		intensity:  7,
		connection: a,
		Commander:  gobot.NewCommander(),
	}

	/* TODO : Add commands */

	return t
}

// Start initializes the tm1638, it uses a SPI-like communication protocol
func (a *AIP1640Driver) Start() error {
	if err := a.pinData.On(); err != nil {
		return err
	}
	return a.pinClock.On()
}

// Halt implements the Driver interface
func (a *AIP1640Driver) Halt() error { return nil }

// Name returns the AIP1640Drivers name
func (a *AIP1640Driver) Name() string { return a.name }

// SetName sets the AIP1640Drivers name
func (a *AIP1640Driver) SetName(n string) { a.name = n }

// Connection returns the AIP1640Driver Connection
func (a *AIP1640Driver) Connection() gobot.Connection {
	return a.connection
}

// SetIntensity changes the intensity (from 1 to 7) of the display
func (a *AIP1640Driver) SetIntensity(level byte) {
	if level >= 7 {
		level = 7
	}
	a.intensity = level
}

// Display sends the buffer to the display (ie. turns on/off the corresponding LEDs)
func (a *AIP1640Driver) Display() error {
	for i := 0; i < 8; i++ {
		if err := a.sendData(byte(i), a.buffer[i]); err != nil {
			return err
		}

		if err := a.pinData.Off(); err != nil {
			return err
		}
		if err := a.pinClock.Off(); err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
		if err := a.pinClock.On(); err != nil {
			return err
		}
		if err := a.pinData.On(); err != nil {
			return err
		}
	}

	return a.sendCommand(AIP1640DispCtrl | a.intensity)
}

// Clear empties the buffer (turns off all the LEDs)
func (a *AIP1640Driver) Clear() {
	for i := 0; i < 8; i++ {
		a.buffer[i] = 0x00
	}
}

// DrawPixel turns on or off a specific in the buffer
func (a *AIP1640Driver) DrawPixel(x, y byte, enabled bool) {
	if x >= 8 || y >= 8 {
		return
	}
	y = 7 - y
	if enabled {
		a.buffer[y] |= 1 << x
	} else {
		a.buffer[y] &^= 1 << x
	}
}

// DrawRow sets any given row of LEDs in the buffer
func (a *AIP1640Driver) DrawRow(row, data byte) {
	if row >= 8 {
		return
	}
	a.buffer[7-row] = data
}

// DrawMatrix sets the whole buffer
func (a *AIP1640Driver) DrawMatrix(data [8]byte) {
	for i := 0; i < 8; i++ {
		a.buffer[7-i] = data[i]
	}
}

// sendCommand is an auxiliary function to send commands to the AIP1640Driver module
func (a *AIP1640Driver) sendCommand(cmd byte) error {
	if err := a.pinData.Off(); err != nil {
		return err
	}
	if err := a.send(cmd); err != nil {
		return err
	}
	return a.pinData.On()
}

// sendData is an auxiliary function to send data to the AIP1640Driver module
func (a *AIP1640Driver) sendData(address byte, data byte) error {
	if err := a.sendCommand(AIP1640DataCmd | AIP1640FixedAddr); err != nil {
		return err
	}
	if err := a.pinData.Off(); err != nil {
		return err
	}
	if err := a.send(AIP1640AddrCmd | address); err != nil {
		return err
	}
	if err := a.send(data); err != nil {
		return err
	}
	return a.pinData.On()
}

// send writes data on the module
func (a *AIP1640Driver) send(data byte) error {
	for i := 0; i < 8; i++ {
		if err := a.pinClock.Off(); err != nil {
			return err
		}

		if (data & 1) > 0 {
			if err := a.pinData.On(); err != nil {
				return err
			}
		} else {
			if err := a.pinData.Off(); err != nil {
				return err
			}
		}
		data >>= 1

		if err := a.pinClock.On(); err != nil {
			return err
		}
	}

	return nil
}
