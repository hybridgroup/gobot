package microbit

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// AccelerometerDriver is the Gobot driver for the Microbit's built-in accelerometer
type AccelerometerDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

type RawAccelerometerData struct {
	X int16
	Y int16
	Z int16
}

type AccelerometerData struct {
	X float32
	Y float32
	Z float32
}

const (
	// BLE services
	accelerometerService = "e95d0753251d470aa062fa1922dfa9a8"

	// BLE characteristics
	accelerometerCharacteristic = "e95dca4b251d470aa062fa1922dfa9a8"

	// Accelerometer event
	Accelerometer = "accelerometer"
)

// NewAccelerometerDriver creates a Microbit AccelerometerDriver
func NewAccelerometerDriver(a ble.BLEConnector) *AccelerometerDriver {
	n := &AccelerometerDriver{
		name:       gobot.DefaultName("Microbit Accelerometer"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent(Accelerometer)

	return n
}

// Connection returns the BLE connection
func (b *AccelerometerDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *AccelerometerDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *AccelerometerDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *AccelerometerDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *AccelerometerDriver) Start() (err error) {
	// subscribe to accelerometer notifications
	b.adaptor().Subscribe(accelerometerCharacteristic, func(data []byte, e error) {
		a := &RawAccelerometerData{X: 0, Y: 0, Z: 0}

		buf := bytes.NewBuffer(data)
		binary.Read(buf, binary.LittleEndian, &a.X)
		binary.Read(buf, binary.LittleEndian, &a.Y)
		binary.Read(buf, binary.LittleEndian, &a.Z)

		result := &AccelerometerData{
			X: float32(a.X) / 1000.0,
			Y: float32(a.Y) / 1000.0,
			Z: float32(a.Z) / 1000.0}

		b.Publish(b.Event(Accelerometer), result)
	})

	return
}

// Halt stops LED driver (void)
func (b *AccelerometerDriver) Halt() (err error) {
	return
}
