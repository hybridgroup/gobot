package i2c

import (
	"math"

	"gobot.io/x/gobot"
)

const (
	defaultAddress = 0x1e // default I2C Address
	registerA      = 0x0  // Address of Configuration register A
	registerB      = 0x01 // Address of Configuration register B
	registerMode   = 0x02 // Address of node register
	xAxisH         = 0x03 // Address of X-axis MSB data register
	zAxisH         = 0x05 // Address of Z-axis MSB data register
	yAxisH         = 0x07 // Address of Y-axis MSB data register
)

// HMC8553LDriver is a Driver for a HMC6352 digital compass
type HMC8553LDriver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewHMC8553LDriver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewHMC8553LDriver(a Connector, options ...func(Config)) *HMC8553LDriver {
	hmc := &HMC8553LDriver{
		name:      gobot.DefaultName("HMC8553L"),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(hmc)
	}

	return hmc
}

// Name returns the name for this Driver
func (h *HMC8553LDriver) Name() string { return h.name }

// SetName sets the name for this Driver
func (h *HMC8553LDriver) SetName(n string) { h.name = n }

// Connection returns the connection for this Driver
func (h *HMC8553LDriver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initializes the HMC8553L
func (h *HMC8553LDriver) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(defaultAddress)
	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}
	if err := h.connection.WriteByteData(registerA, 0x70); err != nil {
		return err
	}
	if err := h.connection.WriteByteData(registerB, 0xa0); err != nil {
		return err
	}
	if err := h.connection.WriteByteData(registerMode, 0); err != nil {
		return err
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *HMC8553LDriver) Halt() (err error) { return }

// ReadRawData reads the raw values from the X, Y, and Z registers
func (h *HMC8553LDriver) ReadRawData() (x int16, y int16, z int16, err error) {
	unsignedX, err := h.connection.ReadWordData(xAxisH)
	if err != nil {
		return
	}
	unsignedY, err := h.connection.ReadWordData(yAxisH)
	if err != nil {
		return
	}
	unsignedZ, err := h.connection.ReadWordData(zAxisH)
	if err != nil {
		return
	}
	return unsignedToSigned(unsignedX), unsignedToSigned(unsignedY), unsignedToSigned(unsignedZ), nil
}

// Heading returns the current heading in radians
func (h *HMC8553LDriver) Heading() (heading float64, err error) {
	var x, y int16
	x, y, _, err = h.ReadRawData()
	if err != nil {
		return
	}
	heading = math.Atan2(float64(y), float64(x))
	if heading > 2*math.Pi {
		heading -= 2 * math.Pi
	}
	if heading < 0 {
		heading += 2 * math.Pi
	}
	return
}

func unsignedToSigned(unsignedValue uint16) int16 {
	if unsignedValue > 32768 {
		return int16(unsignedValue) - ^int16(0) - ^int16(0) - 2
	}
	return int16(unsignedValue)
}
