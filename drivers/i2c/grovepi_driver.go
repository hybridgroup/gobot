package i2c

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// default is for grovepi4 installer
const grovePiDefaultAddress = 0x04

// commands, see:
// * https://www.dexterindustries.com/GrovePi/programming/grovepi-protocol-adding-custom-sensors/
// * https://github.com/DexterInd/GrovePi/tree/1.3.0/Script/multi_grovepi_installer/grovepi4.py
const (
	commandReadDigital         = 1
	commandWriteDigital        = 2
	commandReadAnalog          = 3
	commandWriteAnalog         = 4
	commandSetPinMode          = 5
	commandReadUltrasonic      = 7
	commandReadFirmwareVersion = 8
	commandReadDHT             = 40
)

// GrovePiDriver is a driver for the GrovePi+ for I²C bus interface.
// https://www.dexterindustries.com/grovepi/
//
// To use this driver with the GrovePi, it must be running the firmware >= 1.4.0 and the system version >=3.
// https://github.com/DexterInd/GrovePi/tree/1.3.0/README.md
type GrovePiDriver struct {
	*Driver
	pins map[int]string
}

// NewGrovePiDriver creates a new driver with specified i2c interface
// Params:
//
//	conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewGrovePiDriver(c Connector, options ...func(Config)) *GrovePiDriver {
	d := &GrovePiDriver{
		Driver: NewDriver(c, "GrovePi", grovePiDefaultAddress),
		pins:   make(map[int]string),
	}

	for _, option := range options {
		option(d)
	}

	// TODO: add commands for API
	return d
}

// Connect is here to implement the Adaptor interface.
func (d *GrovePiDriver) Connect() error {
	return nil
}

// Finalize is here to implement the Adaptor interface.
func (d *GrovePiDriver) Finalize() error {
	return nil
}

// AnalogRead returns value from analog pin implementing the AnalogReader interface.
func (d *GrovePiDriver) AnalogRead(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	pinNum, err := d.preparePin(pin, "input")
	if err != nil {
		return 0, err
	}

	buf := []byte{commandReadAnalog, byte(pinNum), 0, 0}
	if _, err := d.connection.Write(buf); err != nil {
		return 0, err
	}

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 3)
	if err = d.readForCommand(commandReadAnalog, data); err != nil {
		return 0, err
	}

	return int(data[1])*256 + int(data[2]), nil
}

// DigitalRead performs a read on a digital pin.
func (d *GrovePiDriver) DigitalRead(pin string) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	pinNum, err := d.preparePin(pin, "input")
	if err != nil {
		return 0, err
	}

	buf := []byte{commandReadDigital, byte(pinNum), 0, 0}
	if _, err := d.connection.Write(buf); err != nil {
		return 0, err
	}

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 2)
	if err = d.readForCommand(commandReadDigital, data); err != nil {
		return 0, err
	}

	return int(data[1]), nil
}

// UltrasonicRead performs a read on an ultrasonic pin with duration >=2 millisecond.
func (d *GrovePiDriver) UltrasonicRead(pin string, duration int) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if duration < 2 {
		duration = 2
	}

	pinNum, err := d.preparePin(pin, "input")
	if err != nil {
		return 0, err
	}

	buf := []byte{commandReadUltrasonic, byte(pinNum), 0, 0}
	if _, err = d.connection.Write(buf); err != nil {
		return 0, err
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)

	data := make([]byte, 3)
	if err := d.readForCommand(commandReadUltrasonic, data); err != nil {
		return 0, err
	}

	return int(data[1])*255 + int(data[2]), nil
}

// FirmwareVersionRead returns the GrovePi firmware version.
func (d *GrovePiDriver) FirmwareVersionRead() (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	buf := []byte{commandReadFirmwareVersion, 0, 0, 0}
	if _, err := d.connection.Write(buf); err != nil {
		return "", err
	}

	time.Sleep(2 * time.Millisecond)

	data := make([]byte, 4)
	if err := d.readForCommand(commandReadFirmwareVersion, data); err != nil {
		return "", err
	}

	return fmt.Sprintf("%v.%v.%v", data[1], data[2], data[3]), nil
}

// DHTRead performs a read temperature and humidity sensors with duration >=2 millisecond.
// DHT11 (blue): sensorType=0
// DHT22 (white): sensorTyp=1
//
//nolint:nonamedreturns // is sufficient here
func (d *GrovePiDriver) DHTRead(pin string, sensorType byte, duration int) (temp float32, hum float32, err error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if duration < 2 {
		duration = 2
	}

	pinNum, err := d.preparePin(pin, "input")
	if err != nil {
		return 0, 0, err
	}

	buf := []byte{commandReadDHT, byte(pinNum), sensorType, 0}
	if _, err = d.connection.Write(buf); err != nil {
		return 0, 0, err
	}
	time.Sleep(time.Duration(duration) * time.Millisecond)

	data := make([]byte, 9)
	if err = d.readForCommand(commandReadDHT, data); err != nil {
		return 0, 0, err
	}

	temp = float32Of4BytesLittleEndian(data[1:5])
	if temp > 150 {
		temp = 150
	}
	if temp < -100 {
		temp = -100
	}

	hum = float32Of4BytesLittleEndian(data[5:9])
	if hum > 100 {
		hum = 100
	}
	if hum < 0 {
		hum = 0
	}

	return temp, hum, err
}

// DigitalWrite writes a value to a specific digital pin implementing the DigitalWriter interface.
func (d *GrovePiDriver) DigitalWrite(pin string, val byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	pinNum, err := d.preparePin(pin, "output")
	if err != nil {
		return err
	}

	buf := []byte{commandWriteDigital, byte(pinNum), val, 0}
	if _, err := d.connection.Write(buf); err != nil {
		return err
	}

	time.Sleep(2 * time.Millisecond)

	_, err = d.connection.ReadByte()
	return err
}

// AnalogWrite writes PWM aka analog to the GrovePi analog pin implementing the AnalogWriter interface.
func (d *GrovePiDriver) AnalogWrite(pin string, val int) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	pinNum, err := d.preparePin(pin, "output")
	if err != nil {
		return err
	}

	buf := []byte{commandWriteAnalog, byte(pinNum), byte(val), 0}
	if _, err := d.connection.Write(buf); err != nil {
		return err
	}

	time.Sleep(2 * time.Millisecond)

	_, err = d.connection.ReadByte()
	return err
}

// SetPinMode sets the pin mode to input or output.
func (d *GrovePiDriver) SetPinMode(pin byte, mode string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return d.setPinMode(pin, mode)
}

func getPin(pin string) string {
	if len(pin) > 1 {
		if strings.ToUpper(pin[0:1]) == "A" || strings.ToUpper(pin[0:1]) == "D" {
			return pin[1:]
		}
	}

	return pin
}

func (d *GrovePiDriver) setPinMode(pin byte, mode string) error {
	var b []byte
	if mode == "output" {
		b = []byte{commandSetPinMode, pin, 1, 0}
	} else {
		b = []byte{commandSetPinMode, pin, 0, 0}
	}
	if _, err := d.connection.Write(b); err != nil {
		return err
	}

	time.Sleep(2 * time.Millisecond)

	_, err := d.connection.ReadByte()
	return err
}

func (d *GrovePiDriver) ensurePinMode(pinNum int, mode string) error {
	if dir, ok := d.pins[pinNum]; !ok || dir != mode {
		if err := d.setPinMode(byte(pinNum), mode); err != nil {
			return err
		}
		d.pins[pinNum] = mode
	}
	return nil
}

func (d *GrovePiDriver) preparePin(pin string, mode string) (int, error) {
	pin = getPin(pin)
	pinNum, err := strconv.Atoi(pin)
	if err != nil {
		return -1, err
	}

	if err := d.ensurePinMode(pinNum, mode); err != nil {
		return -1, err
	}

	return pinNum, nil
}

func (d *GrovePiDriver) readForCommand(command byte, data []byte) error {
	cnt, err := d.connection.Read(data)
	if err != nil {
		return err
	}
	if len(data) != cnt {
		return fmt.Errorf("read count mismatch (%d should be %d)", cnt, len(data))
	}
	if data[0] != command {
		return fmt.Errorf("answer (%d) was not for command (%d)", data[0], command)
	}
	return nil
}

func float32Of4BytesLittleEndian(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
