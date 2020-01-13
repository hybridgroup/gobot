package spi

import (
	"image/color"
	"math"

	"gobot.io/x/gobot"
)

// APA102Driver is a driver for the APA102 programmable RGB LEDs.
type APA102Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander

	vals       []color.RGBA
	brightness uint8
}

// NewAPA102Driver creates a new Gobot Driver for APA102 RGB LEDs.
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver.
//      count int - how many LEDs are in the array controlled by this driver.
//      bright - the default brightness to apply for all LEDs (must be between 0 and 31).
//
// Optional params:
//      spi.WithBus(int):    	  bus to use with this driver.
//      spi.WithChip(int):    	chip to use with this driver.
//      spi.WithMode(int):    	mode to use with this driver.
//      spi.WithBits(int):    	number of bits to use with this driver.
//      spi.WithSpeed(int64):   speed in Hz to use with this driver.
//
func NewAPA102Driver(a Connector, count int, bright uint8, options ...func(Config)) *APA102Driver {
	d := &APA102Driver{
		name:       gobot.DefaultName("APA102"),
		connector:  a,
		vals:       make([]color.RGBA, count),
		brightness: uint8(math.Min(float64(bright), 31)),
		Config:     NewConfig(),
	}
	for _, option := range options {
		option(d)
	}
	return d
}

// Name returns the name of the device.
func (d *APA102Driver) Name() string { return d.name }

// SetName sets the name of the device.
func (d *APA102Driver) SetName(n string) { d.name = n }

// Connection returns the Connection of the device.
func (d *APA102Driver) Connection() gobot.Connection { return d.connection.(gobot.Connection) }

// Start initializes the driver.
func (d *APA102Driver) Start() (err error) {
	bus := d.GetBusOrDefault(d.connector.GetSpiDefaultBus())
	chip := d.GetChipOrDefault(d.connector.GetSpiDefaultChip())
	mode := d.GetModeOrDefault(d.connector.GetSpiDefaultMode())
	bits := d.GetBitsOrDefault(d.connector.GetSpiDefaultBits())
	maxSpeed := d.GetSpeedOrDefault(d.connector.GetSpiDefaultMaxSpeed())

	d.connection, err = d.connector.GetSpiConnection(bus, chip, mode, bits, maxSpeed)
	if err != nil {
		return err
	}
	return nil
}

// Halt stops the driver.
func (d *APA102Driver) Halt() (err error) {
	return
}

// SetRGBA sets the ith LED's color to the given RGBA value.
// A subsequent call to Draw is required to transmit values
// to the LED strip.
func (d *APA102Driver) SetRGBA(i int, v color.RGBA) {
	d.vals[i] = v
}

// SetBrightness sets the ith LED's brightness to the given value.
// Must be between 0 and 31.
func (d *APA102Driver) SetBrightness(i uint8) {
	d.brightness = uint8(math.Min(float64(i), 31))
}

// Brightness return driver brightness value.
func (d *APA102Driver) Brightness() uint8 {
	return d.brightness
}

// Draw displays the RGBA values set on the actual LED strip.
func (d *APA102Driver) Draw() error {
	// TODO(jbd): dotstar allows other RGBA alignments, support those layouts.
	n := len(d.vals)

	tx := make([]byte, 4*(n+1)+(n/2+1))
	tx[0] = 0x00
	tx[1] = 0x00
	tx[2] = 0x00
	tx[3] = 0x00

	for i, c := range d.vals {
		j := (i + 1) * 4
		if c.A != 0 {
			tx[j] = 0xe0 + byte(math.Min(float64(c.A), 31))
		} else {
			tx[j] = 0xe0 + byte(d.brightness)
		}
		tx[j+1] = c.B
		tx[j+2] = c.G
		tx[j+3] = c.R
	}

	// end frame with at least n/2 0xff vals
	for i := (n + 1) * 4; i < len(tx); i++ {
		tx[i] = 0xff
	}

	return d.connection.Tx(tx, nil)
}
