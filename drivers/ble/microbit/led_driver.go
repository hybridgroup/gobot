package microbit

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// ledService = "e95dd91d251d470aa062fa1922dfa9a8"
	ledMatrixStateChara    = "e95d7b77251d470aa062fa1922dfa9a8"
	ledTextChara           = "e95d93ee251d470aa062fa1922dfa9a8"
	ledScrollingDelayChara = "e95d0d2d251d470aa062fa1922dfa9a8"
)

// LEDDriver is the Gobot driver for the Microbit's LED array
type LEDDriver struct {
	*ble.Driver
	gobot.Eventer
}

// NewLEDDriver creates a Microbit LEDDriver
func NewLEDDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *LEDDriver {
	d := &LEDDriver{
		Eventer: gobot.NewEventer(),
	}

	d.Driver = ble.NewDriver(a, "Microbit LED", nil, nil, opts...)

	return d
}

// ReadMatrix read the current LED matrix state
func (b *LEDDriver) ReadMatrix() ([]byte, error) {
	return b.Adaptor().ReadCharacteristic(ledMatrixStateChara)
}

// WriteMatrix writes an array of 5 bytes to set the LED matrix
func (b *LEDDriver) WriteMatrix(data []byte) error {
	return b.Adaptor().WriteCharacteristic(ledMatrixStateChara, data)
}

// WriteText writes a text message to the Microbit LED matrix
func (b *LEDDriver) WriteText(msg string) error {
	return b.Adaptor().WriteCharacteristic(ledTextChara, []byte(msg))
}

func (b *LEDDriver) ReadScrollingDelay() (uint16, error) {
	return 0, nil
}

func (b *LEDDriver) WriteScrollingDelay(delay uint16) error {
	buf := []byte{byte(delay)}
	return b.Adaptor().WriteCharacteristic(ledScrollingDelayChara, buf)
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
