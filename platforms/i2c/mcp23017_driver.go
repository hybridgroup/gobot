/*
Copyright (c) 2015 Ulises Flynn

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package i2c

import (
	"bytes"
	"encoding/binary"
	"log"
	"strings"
	"time"

	"github.com/hybridgroup/gobot"
)

var (
	Debug = true // Set this to true to see debugging information
	// Register this Driver
	_ gobot.Driver = (*MCP23017Driver)(nil)
)

// Port contains all the registers for the device.
type port struct {
	IODIR   byte // I/O direction register: 0=output / 1=input
	IPOL    byte // Input polarity register: 0=normal polarity / 1=inversed
	GPINTEN byte // Interrupt on change control register: 0=disabled / 1=enabled
	DEFVAL  byte // Default compare register for interrupt on change
	INTCON  byte // Interrupt control register: bit set to 0= use defval bit value to compare pin value/ bit set to 1= pin value compared to previous pin value
	IOCON   byte // Configuration register
	GPPU    byte // Pull-up resistor configuration register: 0=enabled / 1=disabled
	INTF    byte // Interrupt flag register: 0=no interrupt / 1=pin caused interrupt
	INTCAP  byte // Interrupt capture register, captures pin values during interrupt: 0=logic low / 1=logic high
	GPIO    byte // Port register, reading from this register reads the port
	OLAT    byte // Output latch register, write modifies the pins: 0=logic low / 1=logic high
}

// Registers in the MCP23017 have different address based on which bank is used.
// Each bank is made up of PortA and PortB registers.
type bank struct {
	PortA port
	PortB port
}

// getBank returns a bank's PortA and PortB registers given a bank number (0/1).
func getBank(bnk uint8) bank {
	if bnk == 0 {
		return bank{PortA: port{0x00, 0x02, 0x04, 0x06, 0x08, 0x0A, 0x0C, 0x0E, 0x10, 0x12, 0x14}, PortB: port{0x01, 0x03, 0x05, 0x07, 0x09, 0x0B, 0x0D, 0x0F, 0x11, 0x13, 0x15}}
	}
	return bank{PortA: port{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}, PortB: port{0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A}}
}

// MCP23017Config contains the device configuration settings to be loaded at startup.
type MCP23017Config struct {
	Bank   uint8
	Mirror uint8
	Seqop  uint8
	Disslw uint8
	Haen   uint8
	Odr    uint8
	Intpol uint8
}

// getBytes returns the mcp23017 configuration as a slice of bytes.
func (conf *MCP23017Config) getBytes() []byte {
	return []byte{conf.Bank, conf.Mirror, conf.Seqop, conf.Disslw, conf.Haen, conf.Odr, conf.Intpol}
}

// MCP23107Driver contains the driver configuration parameters.
type MCP23017Driver struct {
	name            string
	connection      I2c
	conf            MCP23017Config
	mcp23017Address int
	interval        time.Duration
	gobot.Commander
	gobot.Eventer
}

// NewMCP23017Driver creates a new driver with specified name and i2c interface.
func NewMCP23017Driver(a I2c, name string, conf MCP23017Config, deviceAddress int, v ...time.Duration) *MCP23017Driver {
	m := &MCP23017Driver{
		name:            name,
		connection:      a,
		conf:            conf,
		mcp23017Address: deviceAddress,
		Commander:       gobot.NewCommander(),
		Eventer:         gobot.NewEventer(),
	}

	m.AddCommand("WriteGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(float64)
		val := params["val"].(float64)
		port := params["port"].(string)
		return m.WriteGPIO(pin, val, port)
	})

	m.AddCommand("ReadGPIO", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(float64)
		port := params["port"].(string)
		val, err := m.ReadGPIO(pin, port)
		return map[string]interface{}{"val": val, "err": err}
	})

	return m
}

func (m *MCP23017Driver) Name() string { return m.name }

func (m *MCP23017Driver) Connection() gobot.Connection { return m.connection.(gobot.Connection) }

func (m *MCP23017Driver) Halt() (err []error) { return }

// Start writes initialization bytes and reads.
func (m *MCP23017Driver) Start() (errs []error) {
	if err := m.connection.I2cStart(m.mcp23017Address); err != nil {
		return []error{err}
	}
	// Set IOCON register with the given configuration.
	selectedPort := m.getPort("A") // IOCON address is the same for Port A or B.
	var ioval uint8
	buf := bytes.NewReader(m.conf.getBytes())
	err := binary.Read(buf, binary.LittleEndian, &ioval)
	if err != nil {
		return []error{err}
	}
	if err := m.connection.I2cWrite(m.mcp23017Address, []byte{selectedPort.IOCON, ioval}); err != nil {
		return []error{err}
	}
	return
}

// WriteGPIO writes a value to a gpio pin (0-7) and a
// port (A or B).
func (m *MCP23017Driver) WriteGPIO(pin float64, val float64, portStr string) (err error) {
	selectedPort := m.getPort(portStr)
	// Set IODIR register bit for given pin to an output.
	if err := m.write(selectedPort.IODIR, uint8(pin), 0); err != nil {
		return err
	}
	// Set OLAT register to write a value to the given pin.
	if err := m.write(selectedPort.OLAT, uint8(pin), uint8(val)); err != nil {
		return err
	}
	return nil
}

// ReadGPIO reads a value from a given gpio pin (0-7) and a
// port (A or B).
func (m *MCP23017Driver) ReadGPIO(pin float64, portStr string) (val bool, err error) {
	selectedPort := m.getPort(portStr)
	gpio, err := m.read(selectedPort.GPIO)
	if err != nil {
		return false, err
	}
	return ((1 << uint8(pin) & gpio) != 0), nil
}

// SetPullUp sets the pull up state of a given pin based on the value:
// val = 1 pull up enabled.
// val = 0 pull up disabled.
func (m *MCP23017Driver) SetPullUp(pin uint8, val byte, portStr string) error {
	selectedPort := m.getPort(portStr)
	if err := m.write(selectedPort.GPPU, pin, val); err != nil {
		return err
	}
	return nil
}

// SetGPIOPolarity will change a given pin's polarity based on the value:
// val = 1 opposite logic state of the input pin.
// val = 0 same logic state of the input pin.
func (m *MCP23017Driver) SetGPIOPolarity(pin uint8, val byte, portStr string) (err error) {
	selectedPort := m.getPort(portStr)
	if err := m.write(selectedPort.IPOL, pin, val); err != nil {
		return err
	}
	return nil
}

// write gets the value of the passed in register, and then overwrites
// the bit specified by the pin, with the given value.
func (m *MCP23017Driver) write(reg byte, pin uint8, val byte) (err error) {
	var ioval uint8
	iodir, err := m.read(reg)
	if err != nil {
		return err
	}
	if val == 0 {
		ioval = clearBit(iodir, uint8(pin))
	} else if val == 1 {
		ioval = setBit(iodir, uint8(pin))
	}
	if err = m.connection.I2cWrite(m.mcp23017Address, []byte{reg, ioval}); err != nil {
		return err
	}
	return nil
}

// Read returns the values in the given register.
func (m *MCP23017Driver) read(reg byte) (val uint8, err error) {
	bytesToRead := int(reg)
	v, err := m.connection.I2cRead(m.mcp23017Address, bytesToRead+1)
	if err != nil {
		return val, err
	}
	if Debug {
		log.Printf("Register addr:0x%X val: 0x%X\n", reg, v[bytesToRead])
	}
	return v[bytesToRead], nil
}

// getPort return the port (A or B) given a string and the bank.
// Port A is the default if an incorrect or no port is specified.
func (m *MCP23017Driver) getPort(portStr string) (selectedPort port) {
	portStr = strings.ToUpper(portStr)
	switch {
	case portStr == "A":
		return getBank(m.conf.Bank).PortA
	case portStr == "B":
		return getBank(m.conf.Bank).PortB
	default:
		return getBank(m.conf.Bank).PortA
	}
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
