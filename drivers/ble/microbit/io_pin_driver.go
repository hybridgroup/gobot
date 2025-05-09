package microbit

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/common/bit"
)

const (
	// ioPinService = "e95d127b251d470aa062fa1922dfa9a8"
	pinDataChara     = "e95d8d00251d470aa062fa1922dfa9a8"
	pinADConfigChara = "e95d5899251d470aa062fa1922dfa9a8"
	pinIOConfigChara = "e95db9fe251d470aa062fa1922dfa9a8"
)

// IOPinDriver is the Gobot driver for the Microbit's built-in digital and analog I/O
type IOPinDriver struct {
	*ble.Driver
	adMask int
	ioMask int
	gobot.Eventer
}

// pinData has the read data for a specific digital pin
type pinData struct {
	pin   uint8
	value uint8
}

// NewIOPinDriver creates a new driver
func NewIOPinDriver(a gobot.BLEConnector, opts ...ble.OptionApplier) *IOPinDriver {
	d := &IOPinDriver{
		Eventer: gobot.NewEventer(),
	}

	d.Driver = ble.NewDriver(a, "Microbit IO Pins", d.initialize, nil, opts...)

	return d
}

// initialize tells driver to get ready to do work
func (d *IOPinDriver) initialize() error {
	if _, err := d.ReadPinADConfig(); err != nil {
		return err
	}
	_, err := d.ReadPinIOConfig()
	return err
}

// WritePinData writes the pin data for a single pin
func (d *IOPinDriver) WritePinData(pin string, data byte) error {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}

	buf := []byte{byte(i), data}
	err = d.Adaptor().WriteCharacteristic(pinDataChara, buf)
	return err
}

// ReadPinADConfig reads and returns the pin A/D config mask for all pins
func (d *IOPinDriver) ReadPinADConfig() (int, error) {
	c, err := d.Adaptor().ReadCharacteristic(pinADConfigChara)
	if err != nil {
		return 0, err
	}
	var result byte
	for i := 0; i < 4; i++ {
		result |= c[i] << uint(i) //nolint:gosec // ok here
	}

	d.adMask = int(result)
	return int(result), nil
}

// WritePinADConfig writes the pin A/D config mask for all pins
func (d *IOPinDriver) WritePinADConfig(config int) error {
	d.adMask = config
	data := &bytes.Buffer{}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(data, binary.LittleEndian, uint32(config)); err != nil {
		return err
	}

	return d.Adaptor().WriteCharacteristic(pinADConfigChara, data.Bytes())
}

// ReadPinIOConfig reads and returns the pin IO config mask for all pins
func (d *IOPinDriver) ReadPinIOConfig() (int, error) {
	c, err := d.Adaptor().ReadCharacteristic(pinIOConfigChara)
	if err != nil {
		return 0, err
	}

	var result byte
	for i := 0; i < 4; i++ {
		result |= c[i] << uint(i) //nolint:gosec // ok here
	}

	d.ioMask = int(result)
	return int(result), nil
}

// WritePinIOConfig writes the pin I/O config mask for all pins
func (d *IOPinDriver) WritePinIOConfig(config int) error {
	d.ioMask = config
	data := &bytes.Buffer{}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(data, binary.LittleEndian, uint32(config)); err != nil {
		return err
	}

	return d.Adaptor().WriteCharacteristic(pinIOConfigChara, data.Bytes())
}

// DigitalRead reads from a pin
func (d *IOPinDriver) DigitalRead(pin string) (int, error) {
	p, err := validatedPin(pin)
	if err != nil {
		return 0, err
	}

	if err := d.ensureDigital(p); err != nil {
		return 0, err
	}
	if err := d.ensureInput(p); err != nil {
		return 0, err
	}

	pins, err := d.readAllPinData()
	if err != nil {
		return 0, err
	}

	return int(pins[p].value), nil
}

// DigitalWrite writes to a pin
func (d *IOPinDriver) DigitalWrite(pin string, level byte) error {
	p, err := validatedPin(pin)
	if err != nil {
		return err
	}

	if err := d.ensureDigital(p); err != nil {
		return err
	}
	if err := d.ensureOutput(p); err != nil {
		return err
	}

	return d.WritePinData(pin, level)
}

// AnalogRead reads from a pin
func (d *IOPinDriver) AnalogRead(pin string) (int, error) {
	p, err := validatedPin(pin)
	if err != nil {
		return 0, err
	}

	if err := d.ensureAnalog(p); err != nil {
		return 0, err
	}
	if err := d.ensureInput(p); err != nil {
		return 0, err
	}

	pins, err := d.readAllPinData()
	if err != nil {
		return 0, err
	}

	return int(pins[p].value), nil
}

func (d *IOPinDriver) ensureDigital(pin int) error {
	//nolint:gosec // TODO: fix later
	if bit.IsSet(d.adMask, uint8(pin)) {
		return d.WritePinADConfig(bit.Clear(d.adMask, uint8(pin)))
	}

	return nil
}

func (d *IOPinDriver) ensureAnalog(pin int) error {
	//nolint:gosec // TODO: fix later
	if !bit.IsSet(d.adMask, uint8(pin)) {
		return d.WritePinADConfig(bit.Set(d.adMask, uint8(pin)))
	}

	return nil
}

func (d *IOPinDriver) ensureInput(pin int) error {
	//nolint:gosec // TODO: fix later
	if !bit.IsSet(d.ioMask, uint8(pin)) {
		return d.WritePinIOConfig(bit.Set(d.ioMask, uint8(pin)))
	}

	return nil
}

func (d *IOPinDriver) ensureOutput(pin int) error {
	//nolint:gosec // TODO: fix later
	if bit.IsSet(d.ioMask, uint8(pin)) {
		return d.WritePinIOConfig(bit.Clear(d.ioMask, uint8(pin))) //nolint:gosec // TODO: fix later
	}

	return nil
}

func (d *IOPinDriver) readAllPinData() ([]pinData, error) {
	c, _ := d.Adaptor().ReadCharacteristic(pinDataChara)
	buf := bytes.NewBuffer(c)
	pinsData := make([]pinData, buf.Len()/2)

	for i := 0; i < buf.Len()/2; i++ {
		pin, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		value, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		pinData := pinData{
			pin:   pin,
			value: value,
		}
		pinsData[i] = pinData
	}

	return pinsData, nil
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
