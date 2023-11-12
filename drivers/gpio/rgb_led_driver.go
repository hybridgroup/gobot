package gpio

import (
	"log"

	"gobot.io/x/gobot/v2"
)

// RgbLedDriver represents a digital RGB Led
type RgbLedDriver struct {
	pinRed     string
	redColor   byte
	pinGreen   string
	greenColor byte
	pinBlue    string
	blueColor  byte
	name       string
	connection DigitalWriter
	high       bool
	gobot.Commander
}

// NewRgbLedDriver return a new RgbLedDriver given a DigitalWriter and
// 3 pins: redPin, greenPin, and bluePin
//
// Adds the following API Commands:
//
//	"SetRGB" - See RgbLedDriver.SetRGB
//	"Toggle" - See RgbLedDriver.Toggle
//	"On" - See RgbLedDriver.On
//	"Off" - See RgbLedDriver.Off
func NewRgbLedDriver(a DigitalWriter, redPin string, greenPin string, bluePin string) *RgbLedDriver {
	d := &RgbLedDriver{
		name:       gobot.DefaultName("RGBLED"),
		pinRed:     redPin,
		pinGreen:   greenPin,
		pinBlue:    bluePin,
		connection: a,
		high:       false,
		Commander:  gobot.NewCommander(),
	}

	//nolint:forcetypeassert // ok here
	d.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r := byte(params["r"].(int))
		g := byte(params["g"].(int))
		b := byte(params["b"].(int))
		return d.SetRGB(r, g, b)
	})

	d.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		return d.Toggle()
	})

	d.AddCommand("On", func(params map[string]interface{}) interface{} {
		return d.On()
	})

	d.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return d.Off()
	})

	return d
}

// Start implements the Driver interface
func (d *RgbLedDriver) Start() error { return nil }

// Halt implements the Driver interface
func (d *RgbLedDriver) Halt() error { return nil }

// Name returns the RGBLEDDrivers name
func (d *RgbLedDriver) Name() string { return d.name }

// SetName sets the RGBLEDDrivers name
func (d *RgbLedDriver) SetName(n string) { d.name = n }

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

// Connection returns the RgbLedDriver Connection
func (d *RgbLedDriver) Connection() gobot.Connection {
	if conn, ok := d.connection.(gobot.Connection); ok {
		return conn
	}

	log.Printf("%s has no gobot connection\n", d.name)
	return nil
}

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
	if writer, ok := d.connection.(PwmWriter); ok {
		return writer.PwmWrite(pin, level)
	}
	return ErrPwmWriteUnsupported
}

// SetRGB sets the Red Green Blue value of the LED.
func (d *RgbLedDriver) SetRGB(r, g, b byte) error {
	d.redColor = r
	d.greenColor = g
	d.blueColor = b

	return d.On()
}
