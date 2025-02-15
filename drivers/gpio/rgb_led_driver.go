package gpio

import (
	"gobot.io/x/gobot/v2"
)

// RgbLedDriver represents a digital RGB Led
type RgbLedDriver struct {
	*driver
	pinRed     string
	redColor   byte
	pinGreen   string
	greenColor byte
	pinBlue    string
	blueColor  byte
	high       bool
}

// NewRgbLedDriver return a new RgbLedDriver given a PwmWriter and 3 pins: redPin, greenPin, and bluePin
//
// Supported options:
//
//	"WithName"
//
// Adds the following API Commands:
//
//	"SetRGB" - See RgbLedDriver.SetRGB
//	"Toggle" - See RgbLedDriver.Toggle
//	"On" - See RgbLedDriver.On
//	"Off" - See RgbLedDriver.Off
func NewRgbLedDriver(a PwmWriter, redPin string, greenPin string, bluePin string, opts ...interface{}) *RgbLedDriver {
	//nolint:forcetypeassert // no error return value, so there is no better way
	d := &RgbLedDriver{
		driver:   newDriver(a.(gobot.Connection), "RGBLED", opts...),
		pinRed:   redPin,
		pinGreen: greenPin,
		pinBlue:  bluePin,
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r := byte(params["r"].(int))
		g := byte(params["g"].(int))
		b := byte(params["b"].(int))
		return d.SetRGB(r, g, b)
	})

	d.AddCommand("Toggle", func(_ map[string]interface{}) interface{} {
		return d.Toggle()
	})

	d.AddCommand("On", func(_ map[string]interface{}) interface{} {
		return d.On()
	})

	d.AddCommand("Off", func(_ map[string]interface{}) interface{} {
		return d.Off()
	})

	return d
}

// Pin returns the RgbLedDrivers pins
func (d *RgbLedDriver) Pin() string {
	return "r=" + d.pinRed + ", g=" + d.pinGreen + ", b=" + d.pinBlue
}

// RedPin returns the RgbLedDrivers redPin
func (d *RgbLedDriver) RedPin() string { return d.pinRed }

// GreenPin returns the RgbLedDrivers redPin
func (d *RgbLedDriver) GreenPin() string { return d.pinGreen }

// BluePin returns the RgbLedDrivers bluePin
func (d *RgbLedDriver) BluePin() string { return d.pinBlue }

// State return true if the led is On and false if the led is Off
func (d *RgbLedDriver) State() bool {
	return d.high
}

// On sets the led's pins to their various states
func (d *RgbLedDriver) On() error {
	if err := d.SetLevel(d.pinRed, d.redColor); err != nil {
		return err
	}

	if err := d.SetLevel(d.pinGreen, d.greenColor); err != nil {
		return err
	}

	if err := d.SetLevel(d.pinBlue, d.blueColor); err != nil {
		return err
	}

	d.high = true
	return nil
}

// Off sets the led to black.
func (d *RgbLedDriver) Off() error {
	if err := d.SetLevel(d.pinRed, 0); err != nil {
		return err
	}

	if err := d.SetLevel(d.pinGreen, 0); err != nil {
		return err
	}

	if err := d.SetLevel(d.pinBlue, 0); err != nil {
		return err
	}

	d.high = false
	return nil
}

// Toggle sets the led to the opposite of it's current state
func (d *RgbLedDriver) Toggle() error {
	if d.State() {
		return d.Off()
	}

	return d.On()
}

// SetLevel sets the led to the specified color level
func (d *RgbLedDriver) SetLevel(pin string, level byte) error {
	return d.pwmWrite(pin, level)
}

// SetRGB sets the Red Green Blue value of the LED.
func (d *RgbLedDriver) SetRGB(r, g, b byte) error {
	d.redColor = r
	d.greenColor = g
	d.blueColor = b

	return d.On()
}
