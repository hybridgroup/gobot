package i2c

import (
	"gobot.io/x/gobot"
)

const mma7660Address = 0x4c

const (
	MMA7660_X              = 0x00
	MMA7660_Y              = 0x01
	MMA7660_Z              = 0x02
	MMA7660_TILT           = 0x03
	MMA7660_SRST           = 0x04
	MMA7660_SPCNT          = 0x05
	MMA7660_INTSU          = 0x06
	MMA7660_MODE           = 0x07
	MMA7660_STAND_BY       = 0x00
	MMA7660_ACTIVE         = 0x01
	MMA7660_SR             = 0x08
	MMA7660_AUTO_SLEEP_120 = 0x00
	MMA7660_AUTO_SLEEP_64  = 0x01
	MMA7660_AUTO_SLEEP_32  = 0x02
	MMA7660_AUTO_SLEEP_16  = 0x03
	MMA7660_AUTO_SLEEP_8   = 0x04
	MMA7660_AUTO_SLEEP_4   = 0x05
	MMA7660_AUTO_SLEEP_2   = 0x06
	MMA7660_AUTO_SLEEP_1   = 0x07
	MMA7660_PDET           = 0x09
	MMA7660_PD             = 0x0A
)

type MMA7660Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewMMA7660Driver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewMMA7660Driver(a Connector, options ...func(Config)) *MMA7660Driver {
	m := &MMA7660Driver{
		name:      gobot.DefaultName("MMA7660"),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(m)
	}

	// TODO: add commands for API
	return m
}

// Name returns the Name for the Driver
func (h *MMA7660Driver) Name() string { return h.name }

// SetName sets the Name for the Driver
func (h *MMA7660Driver) SetName(n string) { h.name = n }

// Connection returns the connection for the Driver
func (h *MMA7660Driver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initialized the mma7660
func (h *MMA7660Driver) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(mma7660Address)

	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{MMA7660_MODE, MMA7660_STAND_BY}); err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{MMA7660_SR, MMA7660_AUTO_SLEEP_32}); err != nil {
		return err
	}

	if _, err := h.connection.Write([]byte{MMA7660_MODE, MMA7660_ACTIVE}); err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (h *MMA7660Driver) Halt() (err error) { return }

// Acceleration returns the acceleration of the provided x, y, z
func (h *MMA7660Driver) Acceleration(x, y, z float64) (ax, ay, az float64) {
	return x / 21.0, y / 21.0, z / 21.0
}

// XYZ returns the raw x,y and z axis from the mma7660
func (h *MMA7660Driver) XYZ() (x float64, y float64, z float64, err error) {
	buf := []byte{0, 0, 0}
	bytesRead, err := h.connection.Read(buf)
	if err != nil {
		return
	}

	if bytesRead != 3 {
		err = ErrNotEnoughBytes
		return
	}

	for _, val := range buf {
		if ((val >> 6) & 0x01) == 1 {
			err = ErrNotReady
			return
		}
	}

	x = float64((int8(buf[0]) << 2)) / 4.0
	y = float64((int8(buf[1]) << 2)) / 4.0
	z = float64((int8(buf[2]) << 2)) / 4.0

	return
}
