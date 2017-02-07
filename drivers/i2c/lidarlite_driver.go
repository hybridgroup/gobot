package i2c

import (
	"gobot.io/x/gobot"

	"time"
)

const lidarliteAddress = 0x62

type LIDARLiteDriver struct {
	name       string
	connector  I2cConnector
	connection I2cConnection
}

// NewLIDARLiteDriver creates a new driver with specified i2c interface
func NewLIDARLiteDriver(a I2cConnector) *LIDARLiteDriver {
	return &LIDARLiteDriver{
		name:      gobot.DefaultName("LIDARLite"),
		connector: a,
	}
}

func (h *LIDARLiteDriver) Name() string                 { return h.name }
func (h *LIDARLiteDriver) SetName(n string)             { h.name = n }
func (h *LIDARLiteDriver) Connection() gobot.Connection { return h.connector.(gobot.Connection) }

// Start initialized the LIDAR
func (h *LIDARLiteDriver) Start() (err error) {
	bus := h.connector.I2cGetDefaultBus()
	h.connection, err = h.connector.I2cGetConnection(lidarliteAddress, bus)
	if err != nil {
		return err
	}
	return
}

// Halt returns true if devices is halted successfully
func (h *LIDARLiteDriver) Halt() (err error) { return }

// Distance returns the current distance in cm
func (h *LIDARLiteDriver) Distance() (distance int, err error) {
	if _, err = h.connection.Write([]byte{0x00, 0x04}); err != nil {
		return
	}
	time.Sleep(20 * time.Millisecond)

	if _, err = h.connection.Write([]byte{0x0F}); err != nil {
		return
	}

	upper := []byte{0}
	bytesRead, err := h.connection.Read(upper)
	if err != nil {
		return
	}

	if bytesRead != 1 {
		err = ErrNotEnoughBytes
		return
	}

	if _, err = h.connection.Write([]byte{0x10}); err != nil {
		return
	}

	lower := []byte{0}
	bytesRead, err = h.connection.Read(lower)
	if err != nil {
		return
	}

	if bytesRead != 1 {
		err = ErrNotEnoughBytes
		return
	}

	distance = ((int(upper[0]) & 0xff) << 8) | (int(lower[0]) & 0xff)

	return
}
