package gpio

import (
	"time"

	"gobot.io/x/gobot"
)

// NeopixelDriver represents a connection to a NeoPixel
type NeopixelDriver struct {
	name       string
	pin        string
	pixelCount uint16
	pixels     []uint32
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
		pixels:     make([]uint32, pixelCount),
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
	for i := 0; i < int(neo.pixelCount); i++ {
		r := byte((neo.pixels[i] >> 16) & 0xff)
		g := byte((neo.pixels[i] >> 8) & 0xff)
		b := byte(neo.pixels[i] & 0xff)

		neo.writeByte(g)
		neo.writeByte(r)
		neo.writeByte(b)
	}

	neo.reset()
	return nil
}

// SetPixel sets the color of one specific Neopixel in the strip
func (neo *NeopixelDriver) SetPixel(pix uint16, color uint32) error {
	neo.pixels[pix] = color
	return nil
}

// SetConfig sets the config info for the Neopixel strip
func (neo *NeopixelDriver) SetConfig(pin uint8, len uint16) error {
	return nil
}

// writePixel
func (neo *NeopixelDriver) writeByte(b byte) error {
	for i := 0; i < 8; i++ {
		if hasBit(b, i) {
			neo.write1()
		} else {
			neo.write0()
		}
	}
	return nil
}

// write0
func (neo *NeopixelDriver) write0() error {
	neo.connection.DigitalWrite(neo.pin, 1)
	time.Sleep(350 * time.Nanosecond)
	neo.connection.DigitalWrite(neo.pin, 0)
	time.Sleep(800 * time.Nanosecond)
	return nil
}

// write1
func (neo *NeopixelDriver) write1() error {
	neo.connection.DigitalWrite(neo.pin, 1)
	time.Sleep(700 * time.Nanosecond)
	neo.connection.DigitalWrite(neo.pin, 0)
	time.Sleep(600 * time.Nanosecond)
	return nil
}

// frame reset
func (neo *NeopixelDriver) reset() error {
	neo.connection.DigitalWrite(neo.pin, 0)
	time.Sleep(300 * time.Microsecond)
	return nil
}

// check to see if a particular bit is set
func hasBit(n byte, pos int) bool {
	val := n & (1 << uint(pos))
	return (val > 0)
}
