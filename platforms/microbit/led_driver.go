package microbit

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
)

// LEDDriver is the Gobot driver for the Microbit's LED array
type LEDDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

const (
	// BLE services
	// ledService = "e95dd91d251d470aa062fa1922dfa9a8"

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
	//nolint:forcetypeassert // ok here
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *LEDDriver) Start() error { return nil }

// Halt stops LED driver (void)
func (b *LEDDriver) Halt() error { return nil }

// ReadMatrix read the current LED matrix state
func (b *LEDDriver) ReadMatrix() ([]byte, error) {
	return b.adaptor().ReadCharacteristic(ledMatrixStateCharacteristic)
}

// WriteMatrix writes an array of 5 bytes to set the LED matrix
func (b *LEDDriver) WriteMatrix(data []byte) error {
	return b.adaptor().WriteCharacteristic(ledMatrixStateCharacteristic, data)
}

// WriteText writes a text message to the Microbit LED matrix
func (b *LEDDriver) WriteText(msg string) error {
	return b.adaptor().WriteCharacteristic(ledTextCharacteristic, []byte(msg))
}

func (b *LEDDriver) ReadScrollingDelay() (uint16, error) {
	return 0, nil
}

func (b *LEDDriver) WriteScrollingDelay(delay uint16) error {
	buf := []byte{byte(delay)}
	return b.adaptor().WriteCharacteristic(ledScrollingDelayCharacteristic, buf)
}

// Blank clears the LEDs on the Microbit
func (b *LEDDriver) Blank() error {
	buf := []byte{0x00, 0x00, 0x00, 0x00, 0x00}
	return b.WriteMatrix(buf)
}

// Solid turns on all of the Microbit LEDs
func (b *LEDDriver) Solid() error {
	buf := []byte{0x1F, 0x1F, 0x1F, 0x1F, 0x1F}
	return b.WriteMatrix(buf)
}

// UpRightArrow displays an arrow pointing upwards and to the right on the Microbit LEDs
func (b *LEDDriver) UpRightArrow() error {
	buf := []byte{0x0F, 0x03, 0x05, 0x09, 0x10}
	return b.WriteMatrix(buf)
}

// UpLeftArrow displays an arrow pointing upwards and to the left on the Microbit LEDs
func (b *LEDDriver) UpLeftArrow() error {
	buf := []byte{0x1E, 0x18, 0x14, 0x12, 0x01}
	return b.WriteMatrix(buf)
}

// DownRightArrow displays an arrow pointing down and to the right on the Microbit LEDs
func (b *LEDDriver) DownRightArrow() error {
	buf := []byte{0x10, 0x09, 0x05, 0x03, 0x0F}
	return b.WriteMatrix(buf)
}

// DownLeftArrow displays an arrow pointing down and to the left on the Microbit LEDs
func (b *LEDDriver) DownLeftArrow() error {
	buf := []byte{0x01, 0x12, 0x14, 0x18, 0x1E}
	return b.WriteMatrix(buf)
}

// Dimond displays a dimond on the Microbit LEDs
func (b *LEDDriver) Dimond() error {
	buf := []byte{0x04, 0x0A, 0x11, 0x0A, 0x04}
	return b.WriteMatrix(buf)
}

// Smile displays a smile on the Microbit LEDs
func (b *LEDDriver) Smile() error {
	buf := []byte{0x0A, 0x0A, 0x00, 0x11, 0x0E}
	return b.WriteMatrix(buf)
}

// Wink displays a wink on the Microbit LEDs
func (b *LEDDriver) Wink() error {
	buf := []byte{0x08, 0x0B, 0x00, 0x11, 0x0E}
	return b.WriteMatrix(buf)
}
