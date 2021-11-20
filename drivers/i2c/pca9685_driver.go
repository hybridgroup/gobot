package i2c

import (
	"strconv"
	"time"

	"gobot.io/x/gobot"
)

const pca9685Address = 0x40

const (
	PCA9685_MODE1        = 0x00
	PCA9685_MODE2        = 0x01
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

	PCA9685_RESTART = 0x80
	PCA9685_SLEEP   = 0x10
	PCA9685_ALLCALL = 0x01
	PCA9685_INVRT   = 0x10
	PCA9685_OUTDRV  = 0x04
)

// PCA9685Driver is a Gobot Driver for the PCA9685 16-channel 12-bit PWM/Servo controller.
//
// For example, here is the Adafruit board that uses this chip:
// https://www.adafruit.com/product/815
//
type PCA9685Driver struct {
	name       string
	connector  Connector
	connection Connection
	Config
	gobot.Commander
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
		Commander: gobot.NewCommander(),
	}

	for _, option := range options {
		option(p)
	}

	p.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, _ := strconv.Atoi(params["val"].(string))
		return p.PwmWrite(pin, byte(val))
	})
	p.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, _ := strconv.Atoi(params["val"].(string))
		return p.ServoWrite(pin, byte(val))
	})
	p.AddCommand("SetPWM", func(params map[string]interface{}) interface{} {
		channel, _ := strconv.Atoi(params["channel"].(string))
		on, _ := strconv.Atoi(params["on"].(string))
		off, _ := strconv.Atoi(params["off"].(string))
		return p.SetPWM(channel, uint16(on), uint16(off))
	})
	p.AddCommand("SetPWMFreq", func(params map[string]interface{}) interface{} {
		freq, _ := strconv.ParseFloat(params["freq"].(string), 32)
		return p.SetPWMFreq(float32(freq))
	})

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

	if err := p.SetAllPWM(0, 0); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{PCA9685_MODE2, PCA9685_OUTDRV}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{PCA9685_MODE1, PCA9685_ALLCALL}); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1)}); err != nil {
		return err
	}
	oldmode, err := p.connection.ReadByte()
	if err != nil {
		return err
	}
	oldmode = oldmode &^ byte(PCA9685_SLEEP)

	if _, err := p.connection.Write([]byte{PCA9685_MODE1, oldmode}); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	return
}

// Halt stops the device
func (p *PCA9685Driver) Halt() (err error) {
	_, err = p.connection.Write([]byte{PCA9685_ALLLED_OFF_H, 0x10})
	return
}

// SetPWM sets a specific channel to a pwm value from 0-4095.
// Params:
//		channel int - the channel to send the pulse
//		on uint16 - the time to start the pulse
//		off uint16 - the time to stop the pulse
//
// Most typically you set "on" to a zero value, and then set "off" to your desired duty.
//
func (p *PCA9685Driver) SetPWM(channel int, on uint16, off uint16) (err error) {
	if _, err := p.connection.Write([]byte{byte(PCA9685_LED0_ON_L + 4*channel), byte(on) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_LED0_ON_H + 4*channel), byte(on >> 8)}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_LED0_OFF_L + 4*channel), byte(off) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_LED0_OFF_H + 4*channel), byte(off >> 8)}); err != nil {
		return err
	}

	return
}

// SetAllPWM sets all channels to a pwm value from 0-4095.
// Params:
//		on uint16 - the time to start the pulse
//		off uint16 - the time to stop the pulse
//
// Most typically you set "on" to a zero value, and then set "off" to your desired duty.
//
func (p *PCA9685Driver) SetAllPWM(on uint16, off uint16) (err error) {
	if _, err := p.connection.Write([]byte{byte(PCA9685_ALLLED_ON_L), byte(on) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_ALLLED_ON_H), byte(on >> 8)}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_ALLLED_OFF_L), byte(off) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(PCA9685_ALLLED_OFF_H), byte(off >> 8)}); err != nil {
		return err
	}

	return
}

// SetPWMFreq sets the PWM frequency in Hz
func (p *PCA9685Driver) SetPWMFreq(freq float32) error {
	// IC oscillator frequency is 25 MHz
	var prescalevel float32 = 25000000
	// Find frequency of PWM waveform
	prescalevel /= 4096
	// Ratio between desired frequency and maximum
	prescalevel /= freq
	prescalevel -= 1
	// Round value to nearest whole
	prescale := byte(prescalevel + 0.5)

	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1)}); err != nil {
		return err
	}
	oldmode, err := p.connection.ReadByte()
	if err != nil {
		return err
	}

	// Put oscillator in sleep mode, clear bit 7 here to avoid overwriting
	// previous setting
	newmode := (oldmode & 0x7F) | 0x10
	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(newmode)}); err != nil {
		return err
	}
	// Write prescaler value
	if _, err := p.connection.Write([]byte{byte(PCA9685_PRESCALE), prescale}); err != nil {
		return err
	}
	// Put back to old settings
	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(oldmode)}); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	// Enable response to All Call address, enable auto-increment, clear restart
	if _, err := p.connection.Write([]byte{byte(PCA9685_MODE1), byte(oldmode | 0x80)}); err != nil {
		return err
	}

	return nil
}

// PwmWrite writes a PWM signal to the specified channel aka "pin".
// Value values are from 0-255, to conform to the PwmWriter interface.
// If you need finer control, please look at SetPWM().
//
func (p *PCA9685Driver) PwmWrite(pin string, val byte) (err error) {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 255), 0, 4095)
	return p.SetPWM(i, 0, uint16(v))
}

// ServoWrite writes a servo signal to the specified channel aka "pin".
// Valid values are from 0-180, to conform to the ServoWriter interface.
// If you need finer control, please look at SetPWM().
//
func (p *PCA9685Driver) ServoWrite(pin string, val byte) (err error) {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 180), 200, 500)
	return p.SetPWM(i, 0, uint16(v))
}
