package i2c

import "gobot.io/x/gobot"

const (
	// DefaultAddress is the default I2C address for the as5600
	as5600Address = 0x36

	as5600ZMCO = 0x00
	// zero position
	as5600ZPOSMSB = 0x01
	as5600ZPOSLSB = 0x02
	//maximum position
	as5600MPOSMSB = 0x03
	as5600MPOSLSB = 0x04
	// maximum angle
	as5600MANGMSB = 0x05
	as5600MANGLSB = 0x06
	// Customize the device
	as5600CONFMSB = 0x07
	as5600CONFLSB = 0x08
	//Unscaled and unmodified angle
	as5600RAWANGLEMSB = 0x0C
	as5600RAWANGLELSB = 0x0D
	// Scaled output value
	as5600ANGLEMSB = 0x0E
	as5600ANGLELSB = 0x0F
	// Current state (basically strength of magnet)
	as5600STATUS = 0x0B
	// Automatic Gain Control
	as5600AGC = 0x1A
	// Magnitude value of the internal cordic
	as5600MAGNITUDEMSB = 0x1B
	as5600MAGNITUDELSB = 0x1C
	// Burn commands
	as5600BURN = 0xFF
)

// AS5600 Power Mode
const (
	as5600PowerModeNorm = iota
	as5600PowerModeLPM1
	as5600PowerModeLPM2
	as5600PowerModeLPM3
)

// AS5600 Hysteresis
const (
	as5600HysteresisNorm = iota
	as5600Hysteresis1LSB
	as5600Hysteresis2LSB
	as5600Hysteresis3LSB
)

const (
	as5600OutputStageAnalogFull = iota
	as5600OutputStageAnalogReduced
	as5600OutputStageDigitalPWM
)

const (
	_ = 1 << iota
	_
	_
	as5600StatusMHBit
	as5600StatusMLBit
	as5600StatusMDBit
)

// Magnet status
const (
	as5600NoMagnetDetected = iota
	as5600MagnetTooWeak
	as5600MagnetOk
	as5600MagnetTooStrong
	as5600MagnetUnknown
)

// AS5600Driver is a Driver for a AS5600 magnetic encoder
type AS5600Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewAS5600Driver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewAS5600Driver(a Connector, options ...func(Config)) *AS5600Driver {
	as := &AS5600Driver{
		name:      gobot.DefaultName("AS5600"),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(as)
	}

	return as
}

// Name returns the name for this Driver
func (as *AS5600Driver) Name() string {

	return as.name
}

// SetName sets the name for this Driver
func (as *AS5600Driver) SetName(n string) {

	as.name = n
}

// Connection returns the connection for this Driver
func (as *AS5600Driver) Connection() gobot.Connection {

	return as.connector.(gobot.Connection)
}

// Start initializes the AS5600
func (as *AS5600Driver) Start() (err error) {

	bus := as.GetBusOrDefault(as.connector.GetDefaultBus())
	address := as.GetAddressOrDefault(as5600Address)

	if as.connection, err = as.connector.GetConnection(address, bus); err != nil {
		return err
	}

	return nil
}

// Halt returns true if devices is halted successfully
func (as *AS5600Driver) Halt() (err error) {

	return nil
}

// DetectMagnet returns if the magnet is detected
func (as *AS5600Driver) DetecMagnet() (bool, error) {
	var magStatus uint8
	var err error

	magStatus, err = as.connection.ReadByteData(as5600STATUS)
	if err != nil {
		return false, err
	}

	return (magStatus&as5600StatusMDBit != 0), err
}

// GetMagnetStrength returns if the magnet is detected
func (as *AS5600Driver) GetMagnetStrength() (int, error) {
	var magStatus uint8
	var err error

	magStatus, err = as.connection.ReadByteData(as5600STATUS)
	if err != nil {
		return as5600MagnetUnknown, err
	}
	if magStatus&as5600StatusMDBit != 0 {
		return as5600NoMagnetDetected, nil
	}
	if magStatus&as5600StatusMHBit != 0 {
		return as5600MagnetTooStrong, nil
	}
	if magStatus&as5600StatusMLBit != 0 {
		return as5600MagnetTooWeak, nil
	}

	return as5600MagnetOk, nil
}
