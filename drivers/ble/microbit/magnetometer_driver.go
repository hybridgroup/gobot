package microbit

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// magnetometerService = "e95df2d8251d470aa062fa1922dfa9a8"
	magnetometerChara = "e95dfb11251d470aa062fa1922dfa9a8"

	MagnetometerEvent = "magnetometer"
)

// MagnetometerDriver is the Gobot driver for the Microbit's built-in magnetometer
type MagnetometerDriver struct {
	*ble.Driver
	gobot.Eventer
}

type MagnetometerData struct {
	X float32
	Y float32
	Z float32
}

// NewMagnetometerDriver creates a Microbit MagnetometerDriver
func NewMagnetometerDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *MagnetometerDriver {
	d := &MagnetometerDriver{
		Eventer: gobot.NewEventer(),
	}
	d.Driver = ble.NewDriver(a, "Microbit Magnetometer", d.initialize, nil, opts...)

	d.AddEvent(MagnetometerEvent)

	return d
}

// initialize tells driver to get ready to do work
func (d *MagnetometerDriver) initialize() error {
	// subscribe to magnetometer notifications
	return d.Adaptor().Subscribe(magnetometerChara, func(data []byte) {
		a := struct{ x, y, z int16 }{x: 0, y: 0, z: 0}

		buf := bytes.NewBuffer(data)
		if err := binary.Read(buf, binary.LittleEndian, &a.x); err != nil {
			panic(err)
		}
		if err := binary.Read(buf, binary.LittleEndian, &a.y); err != nil {
			panic(err)
		}
		if err := binary.Read(buf, binary.LittleEndian, &a.z); err != nil {
			panic(err)
		}

		result := &MagnetometerData{
			X: float32(a.x) / 1000.0,
			Y: float32(a.y) / 1000.0,
			Z: float32(a.z) / 1000.0,
		}

		d.Publish(d.Event(MagnetometerEvent), result)
	})
}
