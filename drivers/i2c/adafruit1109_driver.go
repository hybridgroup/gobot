package i2c

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

const adafruit1109Debug = false

type adafruit1109PortPin struct {
	port string
	pin  uint8
}

// Adafruit1109Driver is a driver for the 2x16 LCD display with RGB backlit and 5 keys from adafruit, designed for Pi.
// The display is driven by the HD44780, and all is connected by i2c port expander MCP23017.
// https://www.adafruit.com/product/1109
//
// Have to implement DigitalWriter, DigitalReader interface
type Adafruit1109Driver struct {
	name string
	*MCP23017Driver
	redPin    adafruit1109PortPin
	greenPin  adafruit1109PortPin
	bluePin   adafruit1109PortPin
	selectPin adafruit1109PortPin
	upPin     adafruit1109PortPin
	downPin   adafruit1109PortPin
	leftPin   adafruit1109PortPin
	rightPin  adafruit1109PortPin
	rwPin     adafruit1109PortPin
	rsPin     adafruit1109PortPin
	enPin     adafruit1109PortPin
	dataPinD4 adafruit1109PortPin
	dataPinD5 adafruit1109PortPin
	dataPinD6 adafruit1109PortPin
	dataPinD7 adafruit1109PortPin
	*gpio.HD44780Driver
}

// NewAdafruit1109Driver creates is a new driver for the 2x16 LCD display with RGB backlit and 5 keys.
//
// Because HD44780 and MCP23017 are already implemented in gobot, this is a wrapper for using existing implementation.
// So, for the documentation of the parameters, have a look at this drivers.
//
// Tests are done with a Tinkerboard.
func NewAdafruit1109Driver(a Connector, options ...func(Config)) *Adafruit1109Driver {
	options = append(options, WithMCP23017AutoIODirOff(1))
	mcp := NewMCP23017Driver(a, options...)
	d := &Adafruit1109Driver{
		name:           gobot.DefaultName("Adafruit1109"),
		MCP23017Driver: mcp,
		redPin:         adafruit1109PortPin{"A", 6},
		greenPin:       adafruit1109PortPin{"A", 7},
		bluePin:        adafruit1109PortPin{"B", 0},
		selectPin:      adafruit1109PortPin{"A", 0},
		upPin:          adafruit1109PortPin{"A", 3},
		downPin:        adafruit1109PortPin{"A", 2},
		leftPin:        adafruit1109PortPin{"A", 4},
		rightPin:       adafruit1109PortPin{"A", 1},
		rwPin:          adafruit1109PortPin{"B", 6},
		rsPin:          adafruit1109PortPin{"B", 7},
		enPin:          adafruit1109PortPin{"B", 5},
		dataPinD4:      adafruit1109PortPin{"B", 4},
		dataPinD5:      adafruit1109PortPin{"B", 3},
		dataPinD6:      adafruit1109PortPin{"B", 2},
		dataPinD7:      adafruit1109PortPin{"B", 1},
	}
	// mapping for HD44780 to MCP23017 port and IO, 4-Bit data
	dataPins := gpio.HD44780DataPin{
		D4: d.dataPinD4.String(),
		D5: d.dataPinD5.String(),
		D6: d.dataPinD6.String(),
		D7: d.dataPinD7.String(),
	}

	// rwPin := "B_6" not mapped in HD44780 driver
	// at test initialization, there seems rows and columns be switched
	// but inside the driver the row is used as row and col as column
	rows := 2
	columns := 16
	lcd := gpio.NewHD44780Driver(d, columns, rows, gpio.HD44780_4BITMODE, d.rsPin.String(), d.enPin.String(), dataPins,
		gpio.WithHD44780RWPin(d.rwPin.String()))

	d.HD44780Driver = lcd
	return d
}

// Connect implements the adaptor.Connector interface.
// Haven't found any adaptor which implements this with more content.
func (d *Adafruit1109Driver) Connect() error { return nil }

// Finalize implements the adaptor.Connector interface.
// Haven't found any adaptor which implements this with more content.
func (d *Adafruit1109Driver) Finalize() error { return nil }

// Name implements the gobot.Device interface
func (d *Adafruit1109Driver) Name() string {
	return fmt.Sprintf("%s_%s_%s", d.name, d.MCP23017Driver.Name(), d.HD44780Driver.Name())
}

// SetName implements the gobot.Device interface.
func (d *Adafruit1109Driver) SetName(n string) { d.name = n }

// Connection implements the gobot.Device interface.
func (d *Adafruit1109Driver) Connection() gobot.Connection { return d.MCP23017Driver.Connection() }

// Start implements the gobot.Device interface.
func (d *Adafruit1109Driver) Start() error {
	if adafruit1109Debug {
		log.Printf("## MCP.Start ##")
	}
	if err := d.MCP23017Driver.Start(); err != nil {
		return err
	}

	// set all to output (inputs will be set by initButton)
	for pin := uint8(0); pin <= 7; pin++ {
		if err := d.SetPinMode(pin, "A", 0); err != nil {
			return err
		}
		if err := d.SetPinMode(pin, "B", 0); err != nil {
			return err
		}
	}

	// button pins are inputs, has inverse logic and needs pull up
	if err := d.adafruit1109InitButton(d.selectPin); err != nil {
		return err
	}
	if err := d.adafruit1109InitButton(d.upPin); err != nil {
		return err
	}
	if err := d.adafruit1109InitButton(d.downPin); err != nil {
		return err
	}
	if err := d.adafruit1109InitButton(d.leftPin); err != nil {
		return err
	}
	if err := d.adafruit1109InitButton(d.rightPin); err != nil {
		return err
	}

	// lets start with neutral background
	if err := d.SetRGB(true, true, true); err != nil {
		return err
	}
	// set rw pin to write
	if err := d.writePin(d.rwPin, 0x00); err != nil {
		return err
	}
	if adafruit1109Debug {
		log.Printf("## HD.Start ##")
	}
	return d.HD44780Driver.Start()
}

// Halt implements the gobot.Device interface.
func (d *Adafruit1109Driver) Halt() error {
	// we try halt on each device, not stopping on the first error
	var errors []string

	if err := d.HD44780Driver.Halt(); err != nil {
		errors = append(errors, err.Error())
	}
	// switch off the background light
	if err := d.SetRGB(false, false, false); err != nil {
		errors = append(errors, err.Error())
	}
	// must be after HD44780Driver
	if err := d.MCP23017Driver.Halt(); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("Halt the driver %s", strings.Join(errors, ", "))
	}

	return nil
}

// DigitalWrite implements the DigitalWriter interface
// This is called by HD44780 driver to set one gpio output. We redirect the call to the i2c driver MCP23017.
// The given id is the same as defined in dataPins and has the syntax "<port>_<pin>".
func (d *Adafruit1109Driver) DigitalWrite(id string, val byte) error {
	portio := adafruit1109ParseID(id)
	return d.writePin(portio, val)
}

// DigitalRead implements the DigitalReader interface
// This is called by HD44780 driver to read one gpio input. We redirect the call to the i2c driver MCP23017.
// The given id is the same as defined in dataPins and has the syntax "<port>_<pin>".
func (d *Adafruit1109Driver) DigitalRead(id string) (int, error) {
	portio := adafruit1109ParseID(id)
	uval, err := d.readPin(portio)
	if err != nil {
		return 0, err
	}
	return int(uval), err
}

// SetRGB sets the Red Green Blue value of backlit.
// The MCP23017 variant don't support PWM and have inverted logic
func (d *Adafruit1109Driver) SetRGB(r, g, b bool) error {
	if adafruit1109Debug {
		log.Printf("## SetRGB %t, %t, %t ##", r, g, b)
	}
	rio := d.redPin
	gio := d.greenPin
	bio := d.bluePin
	rval := uint8(0x1)
	gval := uint8(0x1)
	bval := uint8(0x1)
	if r {
		rval = 0x00
	}
	if g {
		gval = 0x00
	}
	if b {
		bval = 0x00
	}

	if err := d.writePin(rio, rval); err != nil {
		return err
	}

	if err := d.writePin(gio, gval); err != nil {
		return err
	}

	if err := d.writePin(bio, bval); err != nil {
		return err
	}
	return nil
}

// SelectButton reads the state of the "select" button (1=pressed).
func (d *Adafruit1109Driver) SelectButton() (uint8, error) {
	return d.readPin(d.selectPin)
}

// UpButton reads the state of the "up" button (1=pressed).
func (d *Adafruit1109Driver) UpButton() (uint8, error) {
	return d.readPin(d.upPin)
}

// DownButton reads the state of the "down" button (1=pressed).
func (d *Adafruit1109Driver) DownButton() (uint8, error) {
	return d.readPin(d.downPin)
}

// LeftButton reads the state of the "left" button (1=pressed).
func (d *Adafruit1109Driver) LeftButton() (uint8, error) {
	return d.readPin(d.leftPin)
}

// RightButton reads the state of the "right" button (1=pressed).
func (d *Adafruit1109Driver) RightButton() (uint8, error) {
	return d.readPin(d.rightPin)
}

func (d *Adafruit1109Driver) writePin(ap adafruit1109PortPin, val uint8) error {
	return d.WriteGPIO(ap.pin, ap.port, val)
}

func (d *Adafruit1109Driver) readPin(ap adafruit1109PortPin) (uint8, error) {
	return d.ReadGPIO(ap.pin, ap.port)
}

func (ap *adafruit1109PortPin) String() string {
	return fmt.Sprintf("%s_%d", ap.port, ap.pin)
}

func adafruit1109ParseID(id string) adafruit1109PortPin {
	items := strings.Split(id, "_")
	io := uint8(0)
	if io64, err := strconv.ParseUint(items[1], 10, 32); err == nil {
		io = uint8(io64) //nolint:gosec // TODO: fix later
	}
	return adafruit1109PortPin{port: items[0], pin: io}
}

func (d *Adafruit1109Driver) adafruit1109InitButton(p adafruit1109PortPin) error {
	// make an input
	if err := d.SetPinMode(p.pin, p.port, 1); err != nil {
		return err
	}
	// add pull up resistors
	if err := d.SetPullUp(p.pin, p.port, 1); err != nil {
		return err
	}
	// invert polarity
	if err := d.SetGPIOPolarity(p.pin, p.port, 1); err != nil {
		return err
	}
	return nil
}
