package gpio

import "gobot.io/x/gobot"

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
//	"SetRGB" - See RgbLedDriver.SetRGB
//	"Toggle" - See RgbLedDriver.Toggle
//	"On" - See RgbLedDriver.On
//	"Off" - See RgbLedDriver.Off
func NewRgbLedDriver(a DigitalWriter, redPin string, greenPin string, bluePin string) *RgbLedDriver {
	l := &RgbLedDriver{
		name:       gobot.DefaultName("RGBLED"),
		pinRed:     redPin,
		pinGreen:   greenPin,
		pinBlue:    bluePin,
		connection: a,
		high:       false,
		Commander:  gobot.NewCommander(),
	}

	l.AddCommand("SetRGB", func(params map[string]interface{}) interface{} {
		r := byte(params["r"].(int))
		g := byte(params["g"].(int))
		b := byte(params["b"].(int))
		return l.SetRGB(r, g, b)
	})

	l.AddCommand("Toggle", func(params map[string]interface{}) interface{} {
		return l.Toggle()
	})

	l.AddCommand("On", func(params map[string]interface{}) interface{} {
		return l.On()
	})

	l.AddCommand("Off", func(params map[string]interface{}) interface{} {
		return l.Off()
	})

	return l
}

// Start implements the Driver interface
func (l *RgbLedDriver) Start() (err error) { return }

// Halt implements the Driver interface
func (l *RgbLedDriver) Halt() (err error) { return }

// Name returns the RGBLEDDrivers name
func (l *RgbLedDriver) Name() string { return l.name }

// SetName sets the RGBLEDDrivers name
func (l *RgbLedDriver) SetName(n string) { l.name = n }

// Pin returns the RgbLedDrivers pins
func (l *RgbLedDriver) Pin() string { return "r=" + l.pinRed + ", g=" + l.pinGreen + ", b=" + l.pinBlue }

// RedPin returns the RgbLedDrivers redPin
func (l *RgbLedDriver) RedPin() string { return l.pinRed }

// GreenPin returns the RgbLedDrivers redPin
func (l *RgbLedDriver) GreenPin() string { return l.pinGreen }

// BluePin returns the RgbLedDrivers bluePin
func (l *RgbLedDriver) BluePin() string { return l.pinBlue }

// Connection returns the RgbLedDriver Connection
func (l *RgbLedDriver) Connection() gobot.Connection {
	return l.connection.(gobot.Connection)
}

// State return true if the led is On and false if the led is Off
func (l *RgbLedDriver) State() bool {
	return l.high
}

// On sets the led's pins to their various states
func (l *RgbLedDriver) On() (err error) {
	if err = l.SetLevel(l.pinRed, l.redColor); err != nil {
		return
	}

	if err = l.SetLevel(l.pinGreen, l.greenColor); err != nil {
		return
	}

	if err = l.SetLevel(l.pinBlue, l.blueColor); err != nil {
		return
	}

	l.high = true
	return
}

// Off sets the led to black.
func (l *RgbLedDriver) Off() (err error) {
	if err = l.SetLevel(l.pinRed, 0); err != nil {
		return
	}

	if err = l.SetLevel(l.pinGreen, 0); err != nil {
		return
	}

	if err = l.SetLevel(l.pinBlue, 0); err != nil {
		return
	}

	l.high = false
	return
}

// Toggle sets the led to the opposite of it's current state
func (l *RgbLedDriver) Toggle() (err error) {
	if l.State() {
		err = l.Off()
	} else {
		err = l.On()
	}
	return
}

// SetLevel sets the led to the specified color level
func (l *RgbLedDriver) SetLevel(pin string, level byte) (err error) {
	if writer, ok := l.connection.(PwmWriter); ok {
		return writer.PwmWrite(pin, level)
	}
	return ErrPwmWriteUnsupported
}

// SetRGB sets the Red Green Blue value of the LED.
func (l *RgbLedDriver) SetRGB(r, g, b byte) error {
	l.redColor = r
	l.greenColor = g
	l.blueColor = b

	return l.On()
}
