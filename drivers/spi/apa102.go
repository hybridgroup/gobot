package spi

import (
	"image/color"

	"gobot.io/x/gobot"
)

// APA102Driver is a driver for the APA102 programmable RGB LEDs
type APA102Driver struct {
	name       string
	connector  Connector
	connection Connection

	vals []color.RGBA
}

// NewAPA102Driver creates a new Gobot Driver for APA102 RGB LEDs.
//
// Params:
//      a *Adaptor - the Adaptor to use with this Driver
//		count int - how many LEDs are in the array controlled by this driver
//
func NewAPA102Driver(a Connector, count int) *APA102Driver {
	d := &APA102Driver{
		name:      gobot.DefaultName("APA102"),
		connector: a,
		vals:      make([]color.RGBA, count),
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
	bus := d.connector.GetSpiDefaultBus()
	mode := d.connector.GetSpiDefaultMode()
	maxSpeed := d.connector.GetSpiDefaultMaxSpeed()
	d.connection, err = d.connector.GetSpiConnection(bus, mode, maxSpeed)
	if err != nil {
		return err
	}
	return nil
}

// Halt stops the driver.
func (d *APA102Driver) Halt() (err error) {
	d.connection.Close()
	return
}

// SetRGBA sets the ith LED's color to the given RGBA value.
// A subsequent call to Draw is required to transmit values
// to the LED strip.
func (d *APA102Driver) SetRGBA(i int, v color.RGBA) {
	d.vals[i] = v
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
		tx[j] = 0xe0 + byte(c.R)
		tx[j+1] = byte(c.B)
		tx[j+2] = byte(c.G)
		tx[j+3] = byte(c.R)
	}

	// end frame with at least n/2 0xff vals
	for i := (n + 1) * 4; i < len(tx); i++ {
		tx[i] = 0xff
	}

	return d.connection.Tx(tx, nil)
}
