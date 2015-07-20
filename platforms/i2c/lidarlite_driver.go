package i2c

import (
	"github.com/hybridgroup/gobot"

	"time"
)

var _ gobot.Driver = (*LIDARLiteDriver)(nil)

const lidarliteAddress = 0x62

type LIDARLiteDriver struct {
	name       string
	connection I2c
}

// NewLIDARLiteDriver creates a new driver with specified name and i2c interface
func NewLIDARLiteDriver(a I2c, name string) *LIDARLiteDriver {
	return &LIDARLiteDriver{
		name:       name,
		connection: a,
	}
}

func (h *LIDARLiteDriver) Name() string                 { return h.name }
func (h *LIDARLiteDriver) Connection() gobot.Connection { return h.connection.(gobot.Connection) }

// Start initialized the LIDAR
func (h *LIDARLiteDriver) Start() (errs []error) {
	if err := h.connection.I2cStart(lidarliteAddress); err != nil {
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *LIDARLiteDriver) Halt() (errs []error) { return }

// Distance returns the current distance in cm
func (h *LIDARLiteDriver) Distance() (distance int, err error) {
	if err = h.connection.I2cWrite(lidarliteAddress, []byte{0x00, 0x04}); err != nil {
		return
	}
	<-time.After(20 * time.Millisecond)

	if err = h.connection.I2cWrite(lidarliteAddress, []byte{0x0F}); err != nil {
		return
	}

	upper, err := h.connection.I2cRead(lidarliteAddress, 1)
	if err != nil {
		return
	}

	if len(upper) != 1 {
		err = ErrNotEnoughBytes
		return
	}

	if err = h.connection.I2cWrite(lidarliteAddress, []byte{0x10}); err != nil {
		return
	}

	lower, err := h.connection.I2cRead(lidarliteAddress, 1)
	if err != nil {
		return
	}

	if len(lower) != 1 {
		err = ErrNotEnoughBytes
		return
	}

	distance = ((int(upper[0]) & 0xff) << 8) | (int(lower[0]) & 0xff)

	return
}
