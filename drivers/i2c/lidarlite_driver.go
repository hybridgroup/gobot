package i2c

import (
	"time"
)

const lidarliteDefaultAddress = 0x62

// LIDARLiteDriver is the Gobot driver for the LIDARLite I2C LIDAR device.
type LIDARLiteDriver struct {
	*Driver
}

// NewLIDARLiteDriver creates a new driver for the LIDARLite I2C LIDAR device.
//
// Params:
//		c Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewLIDARLiteDriver(c Connector, options ...func(Config)) *LIDARLiteDriver {
	l := &LIDARLiteDriver{
		Driver: NewDriver(c, "LIDARLite", lidarliteDefaultAddress),
	}

	for _, option := range options {
		option(l)
	}

	// TODO: add commands to API
	return l
}

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
