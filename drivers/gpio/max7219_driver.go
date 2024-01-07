package gpio

import (
	"github.com/hashicorp/go-multierror"

	"gobot.io/x/gobot/v2"
)

// Access and command constants for the driver
const (
	MAX7219Digit0 = 0x01
	MAX7219Digit1 = 0x02
	MAX7219Digit2 = 0x03
	MAX7219Digit3 = 0x04
	MAX7219Digit4 = 0x05
	MAX7219Digit5 = 0x06
	MAX7219Digit6 = 0x07
	MAX7219Digit7 = 0x08

	MAX7219DecodeMode  = 0x09
	MAX7219Intensity   = 0x0a
	MAX7219ScanLimit   = 0x0b
	MAX7219Shutdown    = 0x0c
	MAX7219DisplayTest = 0x0f
)

// MAX7219Driver is the gobot driver for the MAX7219 LED driver
//
// Datasheet: https://datasheets.maximintegrated.com/en/ds/MAX7219-MAX7221.pdf
type MAX7219Driver struct {
	*driver
	pinClock *DirectPinDriver
	pinData  *DirectPinDriver
	pinCS    *DirectPinDriver
	count    uint
}

// NewMAX7219Driver return a new MAX7219Driver given a gobot.Connection, pins and how many chips are chained
//
// Supported options:
//
//	"WithName"
func NewMAX7219Driver(
	a gobot.Connection,
	clockPin, dataPin, csPin string,
	count uint,
	opts ...interface{},
) *MAX7219Driver {
	d := &MAX7219Driver{
		driver:   newDriver(a, "MAX7219", opts...),
		pinClock: NewDirectPinDriver(a, clockPin),
		pinData:  NewDirectPinDriver(a, dataPin),
		pinCS:    NewDirectPinDriver(a, csPin),
		count:    count,
	}
	d.afterStart = d.initialize

	/* TODO : Add commands */

	return d
}

// SetIntensity changes the intensity (from 1 to 7) of the display
func (d *MAX7219Driver) SetIntensity(level byte) error {
	if level > 15 {
		level = 15
	}
	return d.All(MAX7219Intensity, level)
}

// ClearAll turns off all LEDs of all modules
func (d *MAX7219Driver) ClearAll() error {
	var err error
	for i := 1; i <= 8; i++ {
		if e := d.All(byte(i), 0); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return err
}

// ClearOne turns off all LEDs of the given module
func (d *MAX7219Driver) ClearOne(which uint) error {
	var err error
	for i := 1; i <= 8; i++ {
		if e := d.One(which, byte(i), 0); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return err
}

// All sends the same data to all the modules
func (d *MAX7219Driver) All(address byte, data byte) error {
	if err := d.pinCS.Off(); err != nil {
		return err
	}
	var c uint
	for c = 0; c < d.count; c++ {
		if err := d.send(address); err != nil {
			return err
		}
		if err := d.send(data); err != nil {
			return err
		}
	}
	return d.pinCS.On()
}

// One sends data to a specific module
func (d *MAX7219Driver) One(which uint, address byte, data byte) error {
	if err := d.pinCS.Off(); err != nil {
		return err
	}
	var c uint
	for c = 0; c < d.count; c++ {
		if c == which {
			if err := d.send(address); err != nil {
				return err
			}
			if err := d.send(data); err != nil {
				return err
			}
		} else {
			if err := d.send(0); err != nil {
				return err
			}
			if err := d.send(0); err != nil {
				return err
			}
		}
	}
	return d.pinCS.On()
}

// initialize initializes the max7219, it uses a SPI-like communication protocol
func (d *MAX7219Driver) initialize() error {
	if err := d.pinData.On(); err != nil {
		return err
	}
	if err := d.pinClock.On(); err != nil {
		return err
	}
	if err := d.pinCS.On(); err != nil {
		return err
	}

	if err := d.All(MAX7219ScanLimit, 0x07); err != nil {
		return err
	}
	if err := d.All(MAX7219DecodeMode, 0x00); err != nil {
		return err
	}
	if err := d.All(MAX7219Shutdown, 0x01); err != nil {
		return err
	}
	if err := d.All(MAX7219DisplayTest, 0x00); err != nil {
		return err
	}
	if err := d.ClearAll(); err != nil {
		return err
	}
	return d.All(MAX7219Intensity, 0x0f)
}

// send writes data on the module
func (d *MAX7219Driver) send(data byte) error {
	var i byte
	for i = 8; i > 0; i-- {
		mask := byte(0x01 << (i - 1))

		if err := d.pinClock.Off(); err != nil {
			return err
		}
		if data&mask > 0 {
			if err := d.pinData.On(); err != nil {
				return err
			}
		} else {
			if err := d.pinData.Off(); err != nil {
				return err
			}
		}
		if err := d.pinClock.On(); err != nil {
			return err
		}
	}

	return nil
}
