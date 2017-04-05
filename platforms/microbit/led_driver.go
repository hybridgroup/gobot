package microbit

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// LEDDriver is the Gobot driver for the Microbit's LED array
type LEDDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

const (
	// BLE services
	ledService = "e95dd91d251d470aa062fa1922dfa9a8"

	// BLE characteristics
	ledMatrixStateCharacteristic    = "e95d7b77251d470aa062fa1922dfa9a8"
	ledTextCharacteristic           = "e95d93ee251d470aa062fa1922dfa9a8"
	ledScrollingDelayCharacteristic = "e95d0d2d251d470aa062fa1922dfa9a8"
)

// NewLEDDriver creates a Microbit LEDDriver
func NewLEDDriver(a ble.BLEConnector) *LEDDriver {
	n := &LEDDriver{
		name:       gobot.DefaultName("Microbit LED"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the BLE connection
func (b *LEDDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *LEDDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *LEDDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *LEDDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *LEDDriver) Start() (err error) {
	return
}

// Halt stops LED driver (void)
func (b *LEDDriver) Halt() (err error) {
	return
}

// ReadMatrix read the current LED matrix state
func (b *LEDDriver) ReadMatrix() (data []byte, err error) {
	data, err = b.adaptor().ReadCharacteristic(ledMatrixStateCharacteristic)
	return
}

// WriteMatrix writes an array of 5 bytes to set the LED matrix
func (b *LEDDriver) WriteMatrix(data []byte) (err error) {
	err = b.adaptor().WriteCharacteristic(ledMatrixStateCharacteristic, data)
	return
}

// WriteText writes a text message to the Microbit LED matrix
func (b *LEDDriver) WriteText(msg string) (err error) {
	err = b.adaptor().WriteCharacteristic(ledTextCharacteristic, []byte(msg))
	return err
}

func (b *LEDDriver) ReadScrollingDelay() (delay uint16, err error) {
	return
}

func (b *LEDDriver) WriteScrollingDelay(delay uint16) (err error) {
	buf := []byte{byte(delay)}
	err = b.adaptor().WriteCharacteristic(ledScrollingDelayCharacteristic, buf)
	return
}

// Blank clears the LEDs on the Microbit
func (b *LEDDriver) Blank() (err error) {
	buf := []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	err = b.WriteMatrix(buf)
	return
}

// Solid turns on all of the Microbit LEDs
func (b *LEDDriver) Solid() (err error) {
	buf := []byte{0x1F, 0x1F, 0x1F, 0x1F, 0x1F}
	err = b.WriteMatrix(buf)
	return
}

// UpRightArrow displays an arrow pointing upwards and to the right on the Microbit LEDs
func (b *LEDDriver) UpRightArrow() (err error) {
	buf := []byte{0x0F, 0x03, 0x05, 0x09, 0x10}
	err = b.WriteMatrix(buf)
	return
}

// UpLeftArrow displays an arrow pointing upwards and to the left on the Microbit LEDs
func (b *LEDDriver) UpLeftArrow() (err error) {
	buf := []byte{0x1E, 0x18, 0x14, 0x12, 0x01}
	err = b.WriteMatrix(buf)
	return
}

// DownRightArrow displays an arrow pointing down and to the right on the Microbit LEDs
func (b *LEDDriver) DownRightArrow() (err error) {
	buf := []byte{0x10, 0x09, 0x05, 0x03, 0x0F}
	err = b.WriteMatrix(buf)
	return
}

// DownLeftArrow displays an arrow pointing down and to the left on the Microbit LEDs
func (b *LEDDriver) DownLeftArrow() (err error) {
	buf := []byte{0x01, 0x12, 0x14, 0x18, 0x1E}
	err = b.WriteMatrix(buf)
	return
}

// Dimond displays a dimond on the Microbit LEDs
func (b *LEDDriver) Dimond() (err error) {
	buf := []byte{0x04, 0x0A, 0x11, 0x0A, 0x04}
	err = b.WriteMatrix(buf)
	return
}

// Smile displays a smile on the Microbit LEDs
func (b *LEDDriver) Smile() (err error) {
	buf := []byte{0x0A, 0x0A, 0x00, 0x11, 0x0E}
	err = b.WriteMatrix(buf)
	return
}

// Wink displays a wink on the Microbit LEDs
func (b *LEDDriver) Wink() (err error) {
	buf := []byte{0x08, 0x0B, 0x00, 0x11, 0x0E}
	err = b.WriteMatrix(buf)
	return
}
