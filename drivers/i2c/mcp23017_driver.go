package i2c

import (
	"log"
	"strings"

	"gobot.io/x/gobot"
)

const (
	// default address for device when a2/a1/a0 pins are all tied to ground
	mcp23017Address = 0x20
)

var (
	debug = false // toggle debugging information
)

// port contains all the registers for the device.
type port struct {
	IODIR   uint8 // I/O direction register: 0=output / 1=input
	IPOL    uint8 // input polarity register: 0=normal polarity / 1=inversed
	GPINTEN uint8 // interrupt on change control register: 0=disabled / 1=enabled
	DEFVAL  uint8 // default compare register for interrupt on change
	INTCON  uint8 // interrupt control register: bit set to 0= use defval bit value to compare pin value/ bit set to 1= pin value compared to previous pin value
	IOCON   uint8 // configuration register
	GPPU    uint8 // pull-up resistor configuration register: 0=enabled / 1=disabled
	INTF    uint8 // interrupt flag register: 0=no interrupt / 1=pin caused interrupt
	INTCAP  uint8 // interrupt capture register, captures pin values during interrupt: 0=logic low / 1=logic high
	GPIO    uint8 // port register, reading from this register reads the port
	OLAT    uint8 // output latch register, write modifies the pins: 0=logic low / 1=logic high
}

// A bank is made up of PortA and PortB pins.
// Port B pins are on the left side of the chip (starting with pin 1), while port A pins are on the right side.
type bank struct {
	PortA port
	PortB port
}

// MCP23017Config contains the device configuration for the IOCON register.
// These fields should only be set with values 0 or 1.
type MCP23017Config struct {
	Bank   uint8
	Mirror uint8
	Seqop  uint8
	Disslw uint8
	Haen   uint8
	Odr    uint8
	Intpol uint8
}

// MCP23017Driver contains the driver configuration parameters.
type MCP23017Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	MCPConf MCP23017Config
	gobot.Commander
	gobot.Eventer
}

// WithMCP23017Bank option sets the MCP23017Driver bank option
func WithMCP23017Bank(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Bank = val
		} else {
			panic("trying to set bank for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Mirror option sets the MCP23017Driver Mirror option
func WithMCP23017Mirror(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Mirror = val
		} else {
			panic("Trying to set Mirror for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Seqop option sets the MCP23017Driver Seqop option
func WithMCP23017Seqop(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Seqop = val
		} else {
			panic("Trying to set Seqop for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Disslw option sets the MCP23017Driver Disslw option
func WithMCP23017Disslw(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Disslw = val
		} else {
			panic("Trying to set Disslw for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Haen option sets the MCP23017Driver Haen option
func WithMCP23017Haen(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Haen = val
		} else {
			panic("Trying to set Haen for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Odr option sets the MCP23017Driver Odr option
func WithMCP23017Odr(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Odr = val
		} else {
			panic("Trying to set Odr for non-MCP23017Driver")
		}
	}
}

// WithMCP23017Intpol option sets the MCP23017Driver Intpol option
func WithMCP23017Intpol(val uint8) func(Config) {
	return func(c Config) {
		d, ok := c.(*MCP23017Driver)
		if ok {
			d.MCPConf.Intpol = val
		} else {
			panic("Trying to set Intpol for non-MCP23017Driver")
		}
	}
}

// NewMCP23017Driver creates a new Gobot Driver to the MCP23017 i2c port expander.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//		i2c.WithMCP23017Bank(int):	MCP23017 bank to use with this driver
//		i2c.WithMCP23017Mirror(int):	MCP23017 mirror to use with this driver
//		i2c.WithMCP23017Seqop(int):	MCP23017 seqop to use with this driver
//		i2c.WithMCP23017Disslw(int):	MCP23017 disslw to use with this driver
//		i2c.WithMCP23017Haen(int):	MCP23017 haen to use with this driver
//		i2c.WithMCP23017Odr(int):	MCP23017 odr to use with this driver
//		i2c.WithMCP23017Intpol(int):	MCP23017 intpol to use with this driver
//
func NewMCP23017Driver(a Connector, options ...func(Config)) *MCP23017Driver {
	m := &MCP23017Driver{
		name:      gobot.DefaultName("MCP23017"),
		connector: a,
		Config:    NewConfig(),
		MCPConf:   MCP23017Config{},
		Commander: gobot.NewCommander(),
		Eventer:   gobot.NewEventer(),
	}

	for _, option := range options {
		option(m)
	}

	m.AddCommand("WriteGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(uint8)
		val := params["val"].(uint8)
		port := params["port"].(string)
		err := m.WriteGPIO(pin, val, port)
		return map[string]interface{}{"err": err}
	})

	m.AddCommand("ReadGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(uint8)
		port := params["port"].(string)
		val, err := m.ReadGPIO(pin, port)
		return map[string]interface{}{"val": val, "err": err}
	})

	return m
}

// Name return the driver name.
func (m *MCP23017Driver) Name() string { return m.name }

// SetName set the driver name.
func (m *MCP23017Driver) SetName(n string) { m.name = n }

// Connection returns the I2c connection.
func (m *MCP23017Driver) Connection() gobot.Connection { return m.connector.(gobot.Connection) }

// Halt stops the driver.
func (m *MCP23017Driver) Halt() (err error) { return }

// Start writes the device configuration.
func (m *MCP23017Driver) Start() (err error) {
	bus := m.GetBusOrDefault(m.connector.GetDefaultBus())
	address := m.GetAddressOrDefault(mcp23017Address)

	m.connection, err = m.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}
	// Set IOCON register with MCP23017 configuration.
	ioconReg := m.getPort("A").IOCON // IOCON address is the same for Port A or B.
	ioconVal := m.MCPConf.getUint8Value()
	if _, err := m.connection.Write([]uint8{ioconReg, ioconVal}); err != nil {
		return err
	}
	return
}

// WriteGPIO writes a value to a gpio pin (0-7) and a port (A or B).
func (m *MCP23017Driver) WriteGPIO(pin uint8, val uint8, portStr string) (err error) {
	selectedPort := m.getPort(portStr)
	// read current value of IODIR register
	iodir, err := m.read(selectedPort.IODIR)
	if err != nil {
		return err
	}
	// set pin as output by clearing bit
	iodirVal := clearBit(iodir, uint8(pin))
	// write IODIR register bit
	err = m.write(selectedPort.IODIR, uint8(pin), uint8(iodirVal))
	if err != nil {
		return err
	}
	// read current value of OLAT register
	olat, err := m.read(selectedPort.OLAT)
	if err != nil {
		return err
	}
	// set or clear olat value, 0 is no output, 1 is an output
	var olatVal uint8
	if val == 0 {
		olatVal = clearBit(olat, uint8(pin))
	} else {
		olatVal = setBit(olat, uint8(pin))
	}
	// write OLAT register bit
	err = m.write(selectedPort.OLAT, uint8(pin), uint8(olatVal))
	if err != nil {
		return err
	}
	return nil
}

// ReadGPIO reads a value from a given gpio pin (0-7) and a
// port (A or B).
func (m *MCP23017Driver) ReadGPIO(pin uint8, portStr string) (val uint8, err error) {
	selectedPort := m.getPort(portStr)
	// read current value of IODIR register
	iodir, err := m.read(selectedPort.IODIR)
	if err != nil {
		return 0, err
	}
	// set pin as input by setting bit
	iodirVal := setBit(iodir, uint8(pin))
	// write IODIR register bit
	err = m.write(selectedPort.IODIR, uint8(pin), uint8(iodirVal))
	if err != nil {
		return 0, err
	}
	val, err = m.read(selectedPort.GPIO)
	if err != nil {
		return val, err
	}
	val = 1 << uint8(pin) & val
	if val > 1 {
		val = 1
	}
	return val, nil
}

// SetPullUp sets the pull up state of a given pin based on the value:
// val = 1 pull up enabled.
// val = 0 pull up disabled.
func (m *MCP23017Driver) SetPullUp(pin uint8, val uint8, portStr string) error {
	selectedPort := m.getPort(portStr)
	return m.write(selectedPort.GPPU, pin, val)
}

// SetGPIOPolarity will change a given pin's polarity based on the value:
// val = 1 opposite logic state of the input pin.
// val = 0 same logic state of the input pin.
func (m *MCP23017Driver) SetGPIOPolarity(pin uint8, val uint8, portStr string) (err error) {
	selectedPort := m.getPort(portStr)
	return m.write(selectedPort.IPOL, pin, val)
}

// write gets the value of the passed in register, and then overwrites
// the bit specified by the pin, with the given value.
func (m *MCP23017Driver) write(reg uint8, pin uint8, val uint8) (err error) {
	if debug {
		log.Printf("write: MCP address: 0x%X, register:0x%X,value: 0x%X\n", m.GetAddressOrDefault(mcp23017Address), reg, val)
	}
	if _, err = m.connection.Write([]uint8{reg, val}); err != nil {
		return err
	}
	return nil
}

// read get the data from a given register
func (m *MCP23017Driver) read(reg uint8) (val uint8, err error) {
	buf := []byte{0}
	if _, err := m.connection.Write([]uint8{reg}); err != nil {
		return val, err
	}
	bytesRead, err := m.connection.Read(buf)
	if err != nil {
		return val, err
	}
	if bytesRead != 1 {
		err = ErrNotEnoughBytes
		return
	}
	if debug {
		log.Printf("reading: MCP address: 0x%X, register:0x%X,value: 0x%X\n", m.GetAddressOrDefault(mcp23017Address), reg, buf)
	}
	return buf[0], nil
}

// getPort return the port (A or B) given a string and the bank.
// Port A is the default if an incorrect or no port is specified.
func (m *MCP23017Driver) getPort(portStr string) (selectedPort port) {
	portStr = strings.ToUpper(portStr)
	switch {
	case portStr == "A":
		return getBank(m.MCPConf.Bank).PortA
	case portStr == "B":
		return getBank(m.MCPConf.Bank).PortB
	default:
		return getBank(m.MCPConf.Bank).PortA
	}
}

// getUint8Value returns the configuration data as a packed value.
func (mc *MCP23017Config) getUint8Value() uint8 {
	return mc.Bank<<7 | mc.Mirror<<6 | mc.Seqop<<5 | mc.Disslw<<4 | mc.Haen<<3 | mc.Odr<<2 | mc.Intpol<<1
}

// setBit is used to set a bit at a given position to 1.
func setBit(n uint8, pos uint8) uint8 {
	n |= (1 << pos)
	return n
}

// clearBit is used to set a bit at a given position to 0.
func clearBit(n uint8, pos uint8) uint8 {
	mask := ^uint8(1 << pos)
	n &= mask
	return n
}

// getBank returns a bank's PortA and PortB registers given a bank number (0/1).
func getBank(bnk uint8) bank {
	if bnk == 0 {
		return bank{PortA: port{0x00, 0x02, 0x04, 0x06, 0x08, 0x0A, 0x0C, 0x0E, 0x10, 0x12, 0x14}, PortB: port{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F, 0x11, 0x13, 0x15}}
	}
	return bank{PortA: port{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}, PortB: port{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A}}
}
