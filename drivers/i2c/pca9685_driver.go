package i2c

import (
	"strconv"
	"time"

	"gobot.io/x/gobot"
)

const pca9685Address = 0x40

const (
	PCA9685_MODE1        = 0x00
	PCA9685_PRESCALE     = 0xFE
	PCA9685_SUBADR1      = 0x02
	PCA9685_SUBADR2      = 0x03
	PCA9685_SUBADR3      = 0x04
	PCA9685_LED0_ON_L    = 0x06
	PCA9685_LED0_ON_H    = 0x07
	PCA9685_LED0_OFF_L   = 0x08
	PCA9685_LED0_OFF_H   = 0x09
	PCA9685_ALLLED_ON_L  = 0xFA
	PCA9685_ALLLED_ON_H  = 0xFB
	PCA9685_ALLLED_OFF_L = 0xFC
	PCA9685_ALLLED_OFF_H = 0xFD
)

// PCA9685Driver is a Gobot Driver for the PCA9685 16-channel 12-bit
// PWM/Servo controller.
type PCA9685Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
}

// NewPCA9685Driver creates a new driver with specified i2c interface
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewPCA9685Driver(a Connector, options ...func(Config)) *PCA9685Driver {
	p := &PCA9685Driver{
		name:      gobot.DefaultName("PCA9685"),
		connector: a,
		Config:    NewConfig(),
	}

	for _, option := range options {
		option(p)
	}

	// TODO: add commands for API
	return p
}

// Name returns the Name for the Driver
func (p *PCA9685Driver) Name() string { return p.name }

// SetName sets the Name for the Driver
func (p *PCA9685Driver) SetName(n string) { p.name = n }

// Connection returns the connection for the Driver
func (p *PCA9685Driver) Connection() gobot.Connection { return p.connector.(gobot.Connection) }

// Start initializes the pca9685
func (p *PCA9685Driver) Start() (err error) {
	bus := p.GetBusOrDefault(p.connector.GetDefaultBus())
	address := p.GetAddressOrDefault(pca9685Address)

	p.connection, err = p.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{PCA9685_MODE1, 0x00}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{PCA9685_ALLLED_OFF_H, 0x10}); err != nil {
		return err
	}

	return
}

// Halt stops the device
func (p *PCA9685Driver) Halt() (err error) {
	_, err = p.connection.Write([]byte{PCA9685_ALLLED_OFF_H, 0x10})
	return
}

// SetPWM sets a specific channel to a pwm value from 0-4096
func (p *PCA9685Driver) SetPWM(channel int, on uint16, off uint16) (err error) {
	if _, err := p.connection.Write([]byte{byte(PCA9685_LED0_ON_L + 4*channel), byte(on), byte(on >> 8), byte(off), byte(off >> 8)}); err != nil {
		return err
	}

	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (p *PCA9685Driver) SetPWMFreq(freq float32) error {
	freq *= 0.9

	var prescalevel float32 = 25000000
	prescalevel /= 4096
	prescalevel /= freq
	prescalevel -= 1
	prescale := byte(prescalevel + 0.5)

	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1)}); err != nil {
		return err
	}
	data := make([]byte, 1)
	oldmode, err := p.connection.Read(data)
	if err != nil {
		return err
	}

	newmode := (oldmode & 0x7F) | 0x10
	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(newmode)}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_PRESCALE), prescale}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(oldmode)}); err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)
	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(oldmode | 0xa1)}); err != nil {
		return err
	}

	return nil
}

// PwmWrite writes a PWM signal to the specified pin
func (p *PCA9685Driver) PwmWrite(pin string, val byte) (err error) {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 255), 0, 4096)
	return p.SetPWM(i, 0, uint16(v))
}

// ServoWrite writes a servo signal to the specified pin.
// Valid values are from 0-180.
func (p *PCA9685Driver) ServoWrite(pin string, val byte) (err error) {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 180), 200, 500)
	return p.SetPWM(i, 0, uint16(v))
}
