package gpio

import (
	"gobot.io/x/gobot"
)

// NeopixelDriver represents a connection to a NeoPixel
type NeopixelDriver struct {
	name       string
	pin        string
	pixelCount uint16
	connection DigitalWriter
	gobot.Eventer
}

// NewNeopixelDriver returns a new NeopixelDriver
func NewNeopixelDriver(a DigitalWriter, pin string, pixelCount uint16) *NeopixelDriver {
	neo := &NeopixelDriver{
		name:       gobot.DefaultName("Neopixel"),
		connection: a,
		pin:        pin,
		pixelCount: pixelCount,
		Eventer:    gobot.NewEventer(),
	}

	return neo
}

// Start starts up the NeopixelDriver
func (neo *NeopixelDriver) Start() (err error) {
	return
}

// Halt stops the NeopixelDriver
func (neo *NeopixelDriver) Halt() (err error) {
	return
}

// Name returns the Driver's name
func (neo *NeopixelDriver) Name() string { return neo.name }

// SetName sets the Driver's name
func (neo *NeopixelDriver) SetName(n string) { neo.name = n }

// Pin returns the Driver's pin
func (neo *NeopixelDriver) Pin() string { return neo.pin }

// Connection returns the Driver's Connection
func (neo *NeopixelDriver) Connection() gobot.Connection {
	return neo.connection.(gobot.Connection)
}

// Off turns off all the Neopixels in the strip
func (neo *NeopixelDriver) Off() error {
	return nil
}

// Show activates all the Neopixels in the strip
func (neo *NeopixelDriver) Show() error {
	return nil
}

// SetPixel sets the color of one specific Neopixel in the strip
func (neo *NeopixelDriver) SetPixel(pix uint16, color uint32) error {
	return nil
}

// SetConfig sets the config info for the Neopixel strip
func (neo *NeopixelDriver) SetConfig(pin uint8, len uint16) error {
	return nil
}
