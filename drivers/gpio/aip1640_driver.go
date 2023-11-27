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
	*driver
	pinClock  *DirectPinDriver
	pinData   *DirectPinDriver
	intensity byte
	buffer    [8]byte
}

// NewAIP1640Driver return a new driver for AIP1640 LED driver given a gobot.Connection and the clock,
// data and strobe pins.
//
// Supported options:
//
//	"WithName"
func NewAIP1640Driver(a gobot.Connection, clockPin string, dataPin string, opts ...interface{}) *AIP1640Driver {
	d := &AIP1640Driver{
		driver:    newDriver(a, "AIP1640", opts...),
		pinClock:  NewDirectPinDriver(a, clockPin),
		pinData:   NewDirectPinDriver(a, dataPin),
		intensity: 7,
	}
	d.afterStart = d.initialize

	/* TODO : Add commands */

	return d
}

// SetIntensity changes the intensity (from 1 to 7) of the display
func (d *AIP1640Driver) SetIntensity(level byte) {
	if level >= 7 {
		level = 7
	}
	d.intensity = level
}

// Display sends the buffer to the display (ie. turns on/off the corresponding LEDs)
func (d *AIP1640Driver) Display() error {
	for i := 0; i < 8; i++ {
		if err := d.sendData(byte(i), d.buffer[i]); err != nil {
			return err
		}

		if err := d.pinData.Off(); err != nil {
			return err
		}
		if err := d.pinClock.Off(); err != nil {
			return err
		}
		time.Sleep(1 * time.Millisecond)
		if err := d.pinClock.On(); err != nil {
			return err
		}
		if err := d.pinData.On(); err != nil {
			return err
		}
	}

	return d.sendCommand(AIP1640DispCtrl | d.intensity)
}

// Clear empties the buffer (turns off all the LEDs)
func (d *AIP1640Driver) Clear() {
	for i := 0; i < 8; i++ {
		d.buffer[i] = 0x00
	}
}

// DrawPixel turns on or off a specific in the buffer
func (d *AIP1640Driver) DrawPixel(x, y byte, enabled bool) {
	if x >= 8 || y >= 8 {
		return
	}
	y = 7 - y
	if enabled {
		d.buffer[y] |= 1 << x
	} else {
		d.buffer[y] &^= 1 << x
	}
}

// DrawRow sets any given row of LEDs in the buffer
func (d *AIP1640Driver) DrawRow(row, data byte) {
	if row >= 8 {
		return
	}
	d.buffer[7-row] = data
}

// DrawMatrix sets the whole buffer
func (d *AIP1640Driver) DrawMatrix(data [8]byte) {
	for i := 0; i < 8; i++ {
		d.buffer[7-i] = data[i]
	}
}

// initialize initializes the tm1638, it uses a SPI-like communication protocol
func (d *AIP1640Driver) initialize() error {
	if err := d.pinData.On(); err != nil {
		return err
	}
	return d.pinClock.On()
}

// sendCommand is an auxiliary function to send commands to the AIP1640Driver module
func (d *AIP1640Driver) sendCommand(cmd byte) error {
	if err := d.pinData.Off(); err != nil {
		return err
	}
	if err := d.send(cmd); err != nil {
		return err
	}
	return d.pinData.On()
}

// sendData is an auxiliary function to send data to the AIP1640Driver module
func (d *AIP1640Driver) sendData(address byte, data byte) error {
	if err := d.sendCommand(AIP1640DataCmd | AIP1640FixedAddr); err != nil {
		return err
	}
	if err := d.pinData.Off(); err != nil {
		return err
	}
	if err := d.send(AIP1640AddrCmd | address); err != nil {
		return err
	}
	if err := d.send(data); err != nil {
		return err
	}
	return d.pinData.On()
}

// send writes data on the module
func (d *AIP1640Driver) send(data byte) error {
	for i := 0; i < 8; i++ {
		if err := d.pinClock.Off(); err != nil {
			return err
		}

		if (data & 1) > 0 {
			if err := d.pinData.On(); err != nil {
				return err
			}
		} else {
			if err := d.pinData.Off(); err != nil {
				return err
			}
		}
		data >>= 1

		if err := d.pinClock.On(); err != nil {
			return err
		}
	}

	return nil
}
