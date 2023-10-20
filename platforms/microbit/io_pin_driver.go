package microbit

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/ble"
)

// IOPinDriver is the Gobot driver for the Microbit's built-in digital and
// analog I/O
type IOPinDriver struct {
	name       string
	adMask     int
	ioMask     int
	connection gobot.Connection
	gobot.Eventer
}

const (
	// BLE services
	// ioPinService = "e95d127b251d470aa062fa1922dfa9a8"

	// BLE characteristics
	pinDataCharacteristic     = "e95d8d00251d470aa062fa1922dfa9a8"
	pinADConfigCharacteristic = "e95d5899251d470aa062fa1922dfa9a8"
	pinIOConfigCharacteristic = "e95db9fe251d470aa062fa1922dfa9a8"
)

// PinData has the read data for a specific digital pin
type PinData struct {
	pin   uint8
	value uint8
}

// NewIOPinDriver creates a Microbit IOPinDriver
func NewIOPinDriver(a ble.BLEConnector) *IOPinDriver {
	n := &IOPinDriver{
		name:       gobot.DefaultName("Microbit IO Pins"),
		connection: a,
		Eventer:    gobot.NewEventer(),
	}

	return n
}

// Connection returns the BLE connection
func (b *IOPinDriver) Connection() gobot.Connection { return b.connection }

// Name returns the Driver Name
func (b *IOPinDriver) Name() string { return b.name }

// SetName sets the Driver Name
func (b *IOPinDriver) SetName(n string) { b.name = n }

// adaptor returns BLE adaptor
func (b *IOPinDriver) adaptor() ble.BLEConnector {
	return b.Connection().(ble.BLEConnector)
}

// Start tells driver to get ready to do work
func (b *IOPinDriver) Start() error {
	if _, err := b.ReadPinADConfig(); err != nil {
		return err
	}
	_, err := b.ReadPinIOConfig()
	return err
}

// Halt stops driver (void)
func (b *IOPinDriver) Halt() error { return nil }

// ReadAllPinData reads and returns the pin data for all pins
func (b *IOPinDriver) ReadAllPinData() (pins []PinData) {
	c, _ := b.adaptor().ReadCharacteristic(pinDataCharacteristic)
	buf := bytes.NewBuffer(c)
	pinsData := make([]PinData, buf.Len()/2)

	for i := 0; i < buf.Len()/2; i++ {
		pinData := PinData{}
		pinData.pin, _ = buf.ReadByte()
		pinData.value, _ = buf.ReadByte()
		pinsData[i] = pinData
	}

	return pinsData
}

// WritePinData writes the pin data for a single pin
func (b *IOPinDriver) WritePinData(pin string, data byte) error {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	buf := []byte{byte(i), data}
	err = b.adaptor().WriteCharacteristic(pinDataCharacteristic, buf)
	return err
}

// ReadPinADConfig reads and returns the pin A/D config mask for all pins
func (b *IOPinDriver) ReadPinADConfig() (int, error) {
	c, err := b.adaptor().ReadCharacteristic(pinADConfigCharacteristic)
	if err != nil {
		return 0, err
	}
	var result byte
	for i := 0; i < 4; i++ {
		result |= c[i] << uint(i)
	}

	b.adMask = int(result)
	return int(result), nil
}

// WritePinADConfig writes the pin A/D config mask for all pins
func (b *IOPinDriver) WritePinADConfig(config int) error {
	b.adMask = config
	data := &bytes.Buffer{}
	if err := binary.Write(data, binary.LittleEndian, uint32(config)); err != nil {
		return err
	}

	return b.adaptor().WriteCharacteristic(pinADConfigCharacteristic, data.Bytes())
}

// ReadPinIOConfig reads and returns the pin IO config mask for all pins
func (b *IOPinDriver) ReadPinIOConfig() (int, error) {
	c, err := b.adaptor().ReadCharacteristic(pinIOConfigCharacteristic)
	if err != nil {
		return 0, err
	}

	var result byte
	for i := 0; i < 4; i++ {
		result |= c[i] << uint(i)
	}

	b.ioMask = int(result)
	return int(result), nil
}

// WritePinIOConfig writes the pin I/O config mask for all pins
func (b *IOPinDriver) WritePinIOConfig(config int) error {
	b.ioMask = config
	data := &bytes.Buffer{}
	if err := binary.Write(data, binary.LittleEndian, uint32(config)); err != nil {
		return err
	}

	return b.adaptor().WriteCharacteristic(pinIOConfigCharacteristic, data.Bytes())
}

// DigitalRead reads from a pin
func (b *IOPinDriver) DigitalRead(pin string) (int, error) {
	p, err := validatedPin(pin)
	if err != nil {
		return 0, err
	}

	if err := b.ensureDigital(p); err != nil {
		return 0, err
	}
	if err := b.ensureInput(p); err != nil {
		return 0, err
	}

	pins := b.ReadAllPinData()
	return int(pins[p].value), nil
}

// DigitalWrite writes to a pin
func (b *IOPinDriver) DigitalWrite(pin string, level byte) error {
	p, err := validatedPin(pin)
	if err != nil {
		return err
	}

	if err := b.ensureDigital(p); err != nil {
		return err
	}
	if err := b.ensureOutput(p); err != nil {
		return err
	}

	return b.WritePinData(pin, level)
}

// AnalogRead reads from a pin
func (b *IOPinDriver) AnalogRead(pin string) (int, error) {
	p, err := validatedPin(pin)
	if err != nil {
		return 0, err
	}

	if err := b.ensureAnalog(p); err != nil {
		return 0, err
	}
	if err := b.ensureInput(p); err != nil {
		return 0, err
	}

	pins := b.ReadAllPinData()
	return int(pins[p].value), nil
}

func (b *IOPinDriver) ensureDigital(pin int) error {
	if hasBit(b.adMask, pin) {
		return b.WritePinADConfig(clearBit(b.adMask, pin))
	}

	return nil
}

func (b *IOPinDriver) ensureAnalog(pin int) error {
	if !hasBit(b.adMask, pin) {
		return b.WritePinADConfig(setBit(b.adMask, pin))
	}

	return nil
}

func (b *IOPinDriver) ensureInput(pin int) error {
	if !hasBit(b.ioMask, pin) {
		return b.WritePinIOConfig(setBit(b.ioMask, pin))
	}

	return nil
}

func (b *IOPinDriver) ensureOutput(pin int) error {
	if hasBit(b.ioMask, pin) {
		return b.WritePinIOConfig(clearBit(b.ioMask, pin))
	}

	return nil
}

func validatedPin(pin string) (int, error) {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return 0, err
	}

	if i < 0 || i > 2 {
		return 0, errors.New("Invalid pin.")
	}

	return i, nil
}

// via http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go
// Sets the bit at pos in the integer n.
func setBit(n int, pos int) int {
	n |= (1 << uint(pos))
	return n
}

// Test if the bit at pos is set in the integer n.
func hasBit(n int, pos int) bool {
	val := n & (1 << uint(pos))
	return (val > 0)
}

// Clears the bit at pos in n.
func clearBit(n int, pos int) int {
	return n &^ (1 << uint(pos))
}
