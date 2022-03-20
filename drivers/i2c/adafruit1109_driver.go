package i2c

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const adafruit1109Debug = false

type adafruit1109PortPin struct {
	port string
	pin  uint8
}

// have to implement DigitalWriter, DigitalReader interface
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

// Adafruit1109Driver is a driver for the 2x16 LCD display with RGB backlit and 5 keys from adafruit, designed for Pi.
// The display is driven by the HD44780, and all is connected by i2c port expander MCP23017.
// https://www.adafruit.com/product/1109
//
// Because both are already implemented in gobot, we creates a wrapper for using existing implementation.
// So, for the documentation of the parameters, have a look at this drivers.
//
// Tests are done with a tinkerboard.
func NewAdafruit1109Driver(a Connector, options ...func(Config)) *Adafruit1109Driver {
	options = append(options, WithMCP23017AutoIODirOff(1))
	mcp := NewMCP23017Driver(a, options...)
	m := &Adafruit1109Driver{
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
		D4: m.dataPinD4.String(),
		D5: m.dataPinD5.String(),
		D6: m.dataPinD6.String(),
		D7: m.dataPinD7.String(),
	}

	//rwPin := "B_6" not mapped in HD44780 driver
	// at test initialization, there seems rows and columns be switched
	// but inside the driver the row is used as row and col as column
	rows := 2
	columns := 16
	lcd := gpio.NewHD44780Driver(m, columns, rows, gpio.HD44780_4BITMODE, m.rsPin.String(), m.enPin.String(), dataPins)
	lcd.SetRWPin(m.rwPin.String())
	m.HD44780Driver = lcd
	return m
}

// gobot.Connection interface
func (m *Adafruit1109Driver) Name() string {
	return fmt.Sprintf("%s_%s_%s", m.name, m.MCP23017Driver.Name(), m.HD44780Driver.Name())
}
func (m *Adafruit1109Driver) SetName(n string) { m.name = n }

// gobot.Device interface
func (m *Adafruit1109Driver) Connection() gobot.Connection { return m.MCP23017Driver.Connection() }
func (m *Adafruit1109Driver) Halt() (err error)            { return m.MCP23017Driver.Halt() }

func (m *Adafruit1109Driver) Start() (err error) {
	if adafruit1109Debug {
		log.Printf("## MCP.Start ##")
	}
	if err = m.MCP23017Driver.Start(); err != nil {
		return err
	}

	// set all to output (inputs will be set by initButton)
	for pin := uint8(0); pin <= 7; pin++ {
		if err := m.PinMode(pin, 0, "A"); err != nil {
			return err
		}
		if err := m.PinMode(pin, 0, "B"); err != nil {
			return err
		}
	}

	// button pins are inputs, has inverse logic and needs pull up
	if err := m.adafruit1109InitButton(m.selectPin); err != nil {
		return err
	}
	if err := m.adafruit1109InitButton(m.upPin); err != nil {
		return err
	}
	if err := m.adafruit1109InitButton(m.downPin); err != nil {
		return err
	}
	if err := m.adafruit1109InitButton(m.leftPin); err != nil {
		return err
	}
	if err := m.adafruit1109InitButton(m.rightPin); err != nil {
		return err
	}

	// lets start with neutral background
	if err = m.SetRGB(true, true, true); err != nil {
		return err
	}
	// set rw pin to write
	if err := m.writePin(m.rwPin, 0x00); err != nil {
		return err
	}
	if adafruit1109Debug {
		log.Printf("## HD.Start ##")
	}
	return m.HD44780Driver.Start()
}

// DigitalWriter interface
// This is called by HD44780 driver to set one gpio output. We redirect the call to the i2c driver MCP23017.
// The given id is the same as defined in dataPins and has the syntax "<port>_<pin>".
func (m *Adafruit1109Driver) DigitalWrite(id string, val byte) (err error) {
	portio := adafruit1109ParseId(id)
	return m.writePin(portio, val)
}

// DigitalReader interface
// This is called by HD44780 driver to read one gpio input. We redirect the call to the i2c driver MCP23017.
// The given id is the same as defined in dataPins and has the syntax "<port>_<pin>".
func (m *Adafruit1109Driver) DigitalRead(id string) (val int, err error) {
	portio := adafruit1109ParseId(id)
	uval, err := m.readPin(portio)
	if err != nil {
		return 0, err
	}
	return int(uval), err
}

// Connector interface, haven't found any adaptor which implements this with more content
func (m *Adafruit1109Driver) Connect() (err error)  { return }
func (m *Adafruit1109Driver) Finalize() (err error) { return }

// SetRGB sets the Red Green Blue value of backlit.
// The MCP23017 variant don't support PWM and have inverted logic
func (m *Adafruit1109Driver) SetRGB(r, g, b bool) error {
	if adafruit1109Debug {
		log.Printf("## SetRGB %t, %t, %t ##", r, g, b)
	}
	rio := m.redPin
	gio := m.greenPin
	bio := m.bluePin
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

	if err := m.writePin(rio, rval); err != nil {
		return err
	}

	if err := m.writePin(gio, gval); err != nil {
		return err
	}

	if err := m.writePin(bio, bval); err != nil {
		return err
	}
	return nil
}

func (m *Adafruit1109Driver) SelectButton() (uint8, error) {
	return m.readPin(m.selectPin)
}

func (m *Adafruit1109Driver) UpButton() (uint8, error) {
	return m.readPin(m.upPin)
}

func (m *Adafruit1109Driver) DownButton() (uint8, error) {
	return m.readPin(m.downPin)
}

func (m *Adafruit1109Driver) LeftButton() (uint8, error) {
	return m.readPin(m.leftPin)
}

func (m *Adafruit1109Driver) RightButton() (uint8, error) {
	return m.readPin(m.rightPin)
}

func (m *Adafruit1109Driver) writePin(ap adafruit1109PortPin, val uint8) (err error) {
	return m.WriteGPIO(ap.pin, val, ap.port)
}

func (m *Adafruit1109Driver) readPin(ap adafruit1109PortPin) (uint8, error) {
	return m.ReadGPIO(ap.pin, ap.port)
}

func (ap *adafruit1109PortPin) String() string {
	return fmt.Sprintf("%s_%d", ap.port, ap.pin)
}

func adafruit1109ParseId(id string) adafruit1109PortPin {
	items := strings.Split(id, "_")
	io := uint8(0)
	if io64, err := strconv.ParseUint(items[1], 10, 32); err == nil {
		io = uint8(io64)
	}
	return adafruit1109PortPin{port: items[0], pin: io}
}

func (m *Adafruit1109Driver) adafruit1109InitButton(p adafruit1109PortPin) error {
	// make an input
	if err := m.PinMode(p.pin, 1, p.port); err != nil {
		return err
	}
	// add pull up resistors
	if err := m.SetPullUp(p.pin, 1, p.port); err != nil {
		return err
	}
	// invert polarity
	if err := m.SetGPIOPolarity(p.pin, 1, p.port); err != nil {
		return err
	}
	return nil
}
