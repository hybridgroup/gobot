package microbit

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
)

const (
	// accelerometerService = "e95d0753251d470aa062fa1922dfa9a8"
	accelerometerChara = "e95dca4b251d470aa062fa1922dfa9a8"

	AccelerometerEvent = "accelerometer"
)

// AccelerometerDriver is the Gobot driver for the Microbit's built-in accelerometer
type AccelerometerDriver struct {
	*ble.Driver
	gobot.Eventer
}

type AccelerometerData struct {
	X float32
	Y float32
	Z float32
}

// NewAccelerometerDriver creates a  AccelerometerDriver
func NewAccelerometerDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *AccelerometerDriver {
	d := &AccelerometerDriver{
		Eventer: gobot.NewEventer(),
	}
	d.Driver = ble.NewDriver(a, "Microbit Accelerometer", d.initialize, nil, opts...)

	d.AddEvent(AccelerometerEvent)

	return d
}

// initialize tells driver to get ready to do work
func (d *AccelerometerDriver) initialize() error {
	// subscribe to accelerometer notifications
	return d.Adaptor().Subscribe(accelerometerChara, func(data []byte) {
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

		result := &AccelerometerData{
			X: float32(a.x) / 1000.0,
			Y: float32(a.y) / 1000.0,
			Z: float32(a.z) / 1000.0,
		}

		d.Publish(d.Event(AccelerometerEvent), result)
	})
}
