package firmata

import (
	"strconv"

	"gobot.io/x/gobot"
)

const (
	NEOPIXEL_CMD = 0x51
	// turn strip off
	NEOPIXEL_OFF = 0x00
	// configure the strip
	NEOPIXEL_CONFIG = 0x01
	// show currently set pixels
	NEOPIXEL_SHOW = 0x02
	// set the color value of pixel n using 32bit packed color value
	NEOPIXEL_SET_PIXEL = 0x03
	// TODO: set color of whole strip
	NEOPIXEL_SET_STRIP = 0x04
	// TODO: shift all the pixels n places along the strip
	NEOPIXEL_SHIFT = 0x05
)

// NeopixelDriver represents a connection to a NeoPixel
type NeopixelDriver struct {
	name       string
	pin        string
	pixelCount uint16
	connection *Adaptor
	gobot.Eventer
}

// NewNeopixelDriver returns a new NeopixelDriver
func NewNeopixelDriver(a *Adaptor, pin string, pixelCount uint16) *NeopixelDriver {
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
	i, _ := strconv.Atoi(neo.pin)
	return neo.SetConfig(uint8(i), neo.pixelCount)
}

// Halt stops the NeopixelDriver
func (neo *NeopixelDriver) Halt() (err error) {
	return neo.Off()
}

// Name returns the Driver's name
func (neo *NeopixelDriver) Name() string { return neo.name }

// SetName sets the Driver's name
func (neo *NeopixelDriver) SetName(n string) { neo.name = n }

// Pin returns the Driver's pin
func (neo *NeopixelDriver) Pin() string { return neo.pin }

// Connection returns the Driver's Connection
func (neo *NeopixelDriver) Connection() gobot.Connection { return neo.connection }

// Off turns off all the Neopixels in the strip
func (neo *NeopixelDriver) Off() error {
	return neo.connection.WriteSysex([]byte{NEOPIXEL_CMD, NEOPIXEL_OFF})
}

// Show activates all the Neopixels in the strip
func (neo *NeopixelDriver) Show() error {
	return neo.connection.WriteSysex([]byte{NEOPIXEL_CMD, NEOPIXEL_SHOW})
}

// SetPixel sets the color of one specific Neopixel in the strip
func (neo *NeopixelDriver) SetPixel(pix uint16, color uint32) error {
	return neo.connection.WriteSysex([]byte{NEOPIXEL_CMD, NEOPIXEL_SET_PIXEL,
		byte(pix & 0x7F),
		byte((pix >> 7) & 0x7F),
		byte(color & 0x7F),
		byte((color >> 7) & 0x7F),
		byte((color >> 14) & 0x7F),
		byte((color >> 21) & 0x7F),
	})
}

// SetConfig sets the config info for the Neopixel strip
func (neo *NeopixelDriver) SetConfig(pin uint8, len uint16) error {
	return neo.connection.WriteSysex([]byte{NEOPIXEL_CMD, NEOPIXEL_CONFIG,
		byte(pin),
		byte(len & 0x7F),
		byte((len >> 7) & 0x7F),
	})
}
