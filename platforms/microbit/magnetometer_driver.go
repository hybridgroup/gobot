package microbit

import (
	"bytes"
	"encoding/binary"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

// MagnetometerDriver is the Gobot driver for the Microbit's built-in magnetometer
type MagnetometerDriver struct {
	name       string
	connection gobot.Connection
	gobot.Eventer
}

type RawMagnetometerData struct {
	X int16
	Y int16
	Z int16
}

type MagnetometerData struct {
	X float32
	Y float32
	Z float32
}

const (
	// BLE services
	magnetometerService = "e95df2d8251d470aa062fa1922dfa9a8"

	// BLE characteristics
	magnetometerCharacteristic = "e95dfb11251d470aa062fa1922dfa9a8"

	// Magnetometer event
	Magnetometer = "magnetometer"
)

// NewMagnetometerDriver creates a Microbit MagnetometerDriver
func NewMagnetometerDriver(a ble.BLEConnector) *MagnetometerDriver {
	n := &MagnetometerDriver{
		name:       gobot.DefaultName("Microbit Magnetometer"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	n.AddEvent(Magnetometer)

	return n
}

// Connection returns the BLE connection
func (b *MagnetometerDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *MagnetometerDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *MagnetometerDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *MagnetometerDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *MagnetometerDriver) Start() (err error) {
	// subscribe to magnetometer notifications
	b.adaptor().Subscribe(magnetometerCharacteristic, func(data []byte, e error) {
		a := &RawMagnetometerData{X: 0, Y: 0, Z: 0}

		buf := bytes.NewBuffer(data)
		binary.Read(buf, binary.LittleEndian, &a.X)
		binary.Read(buf, binary.LittleEndian, &a.Y)
		binary.Read(buf, binary.LittleEndian, &a.Z)

		result := &MagnetometerData{
			X: float32(a.X) / 1000.0,
			Y: float32(a.Y) / 1000.0,
			Z: float32(a.Z) / 1000.0}

		b.Publish(b.Event(Magnetometer), result)
	})

	return
}

// Halt stops LED driver (void)
func (b *MagnetometerDriver) Halt() (err error) {
	return
}
