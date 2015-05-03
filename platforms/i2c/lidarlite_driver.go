package i2c

import (
	"github.com/hybridgroup/gobot"

	"time"
)

var _ gobot.Driver = (*LIDARLiteDriver)(nil)

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
	if err := h.connection.I2cStart(0x62); err != nil {
		return []error{err}
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *LIDARLiteDriver) Halt() (errs []error) { return }

// Distance returns the current distance
func (h *LIDARLiteDriver) Distance() (distance int, err error) {
	if err = h.connection.I2cWrite([]byte{0x00, 0x04}); err != nil {
		return
	}
	<-time.After(20 * time.Millisecond)

	if err = h.connection.I2cWrite([]byte{0x8f}); err != nil {
		return
	}
	<-time.After(20 * time.Millisecond)

	ret, err := h.connection.I2cRead(2)
	if err != nil {
		return
	}
	if len(ret) == 2 {
		distance = (int(ret[1]) + int(ret[0])*256)
		return
	} else {
		err = ErrNotEnoughBytes
	}
	return
}
