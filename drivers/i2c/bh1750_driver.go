package i2c

import (
	"errors"
	"time"
)

const bh1750DefaultAddress = 0x23

const (
	BH1750_POWER_DOWN                 = 0x00
	BH1750_POWER_ON                   = 0x01
	BH1750_RESET                      = 0x07
	BH1750_CONTINUOUS_HIGH_RES_MODE   = 0x10
	BH1750_CONTINUOUS_HIGH_RES_MODE_2 = 0x11
	BH1750_CONTINUOUS_LOW_RES_MODE    = 0x13
	BH1750_ONE_TIME_HIGH_RES_MODE     = 0x20
	BH1750_ONE_TIME_HIGH_RES_MODE_2   = 0x21
	BH1750_ONE_TIME_LOW_RES_MODE      = 0x23
)

// BH1750Driver is a driver for the BH1750 digital Ambient Light Sensor IC for I²C bus interface.
type BH1750Driver struct {
	*Driver
	mode byte
}

// NewBH1750Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewBH1750Driver(c Connector, options ...func(Config)) *BH1750Driver {
	h := &BH1750Driver{
		Driver: NewDriver(c, "BH1750", bh1750DefaultAddress),
		mode:   BH1750_CONTINUOUS_HIGH_RES_MODE,
	}
	h.afterStart = h.initialize

	for _, option := range options {
		option(h)
	}

	// TODO: add commands for API
	return h
}

// RawSensorData returns the raw value from the bh1750
func (h *BH1750Driver) RawSensorData() (level int, err error) {
	buf := []byte{0, 0}
	bytesRead, err := h.connection.Read(buf)
	if bytesRead != 2 {
		err = errors.New("wrong number of bytes read")
		return
	}
	if err != nil {
		return
	}
	level = int(buf[0])<<8 | int(buf[1])

	return
}

// Lux returns the adjusted value from the bh1750
func (h *BH1750Driver) Lux() (lux int, err error) {
	lux, err = h.RawSensorData()
	lux = int(float64(lux) / 1.2)

	return
}

func (h *BH1750Driver) initialize() error {
	err := h.connection.WriteByte(h.mode)
	time.Sleep(10 * time.Microsecond)
	if err != nil {
		return err
	}
	return nil
}
