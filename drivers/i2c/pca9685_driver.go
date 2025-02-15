package i2c

import (
	"strconv"
	"time"

	"gobot.io/x/gobot/v2"
)

const pca9685DefaultAddress = 0x40

const (
	pca9685Mode1Reg      = 0x00
	pca9685Mode2Reg      = 0x01
	pca9685Subadr1Reg    = 0x02
	pca9685Subadr2Reg    = 0x03
	pca9685Subadr3Reg    = 0x04
	pca9685Led0OnLReg    = 0x06
	pca9685Led0OnHReg    = 0x07
	pca9685Led0OffLReg   = 0x08
	pca9685Led0OffHReg   = 0x09
	pca9685AllLedOnLReg  = 0xFA
	pca9685AllLedOnHReg  = 0xFB
	pca9685AllLedOffLReg = 0xFC
	pca9685AllLedOffHReg = 0xFD
	pca9685PrescaleReg   = 0xFE

	pca9685Mode1RegRestartBit    = 0x80 // bit 7 - 0: restart disabled (default)
	pca9685Mode1RegSleepBit      = 0x10 // bit 4 - 0: normal, 1: low power (default)
	pca9685Mode1RegAllCallBit    = 0x01 // bit 0 - 0: no response to all-call, 1: respond to all-call (default)
	pca9685Mode2RegInvertBit     = 0x10 // bit 4 - 0: outputs not inverted (default), 1: outputs inverted
	pca9685Mode2RegOutdrvBit     = 0x04 // bit 2 - 0: open-drain, 1: totem-pole (default)
	pca9685AllLedOffHRegShutDown = 0x10 // bit 4 - 1: orderly shut down
)

// PCA9685Driver is a Gobot Driver for the PCA9685 16-channel 12-bit PWM/Servo controller.
//
// For example, here is the Adafruit board that uses this chip:
// https://www.adafruit.com/product/815
type PCA9685Driver struct {
	*Driver
}

// NewPCA9685Driver creates a new driver with specified i2c interface
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewPCA9685Driver(c Connector, options ...func(Config)) *PCA9685Driver {
	p := &PCA9685Driver{
		Driver: NewDriver(c, "PCA9685", pca9685DefaultAddress),
	}
	p.afterStart = p.initialize
	p.beforeHalt = p.shutdown

	for _, option := range options {
		option(p)
	}

	//nolint:forcetypeassert // ok here
	p.AddCommand("PwmWrite", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, _ := strconv.Atoi(params["val"].(string))
		return p.PwmWrite(pin, byte(val))
	})
	//nolint:forcetypeassert // ok here
	p.AddCommand("ServoWrite", func(params map[string]interface{}) interface{} {
		pin := params["pin"].(string)
		val, _ := strconv.Atoi(params["val"].(string))
		return p.ServoWrite(pin, byte(val))
	})
	//nolint:forcetypeassert // ok here
	p.AddCommand("SetPWM", func(params map[string]interface{}) interface{} {
		channel, _ := strconv.Atoi(params["channel"].(string))
		on, _ := strconv.Atoi(params["on"].(string))
		off, _ := strconv.Atoi(params["off"].(string))
		return p.SetPWM(channel, uint16(on), uint16(off)) //nolint:gosec // TODO: fix later
	})
	//nolint:forcetypeassert // ok here
	p.AddCommand("SetPWMFreq", func(params map[string]interface{}) interface{} {
		freq, _ := strconv.ParseFloat(params["freq"].(string), 32)
		return p.SetPWMFreq(float32(freq))
	})

	return p
}

// SetPWM sets a specific channel to a pwm value from 0-4095.
// Params:
//
//	channel int - the channel to send the pulse
//	on uint16 - the time to start the pulse
//	off uint16 - the time to stop the pulse
//
// Most typically you set "on" to a zero value, and then set "off" to your desired duty.
func (p *PCA9685Driver) SetPWM(channel int, on uint16, off uint16) error {
	if _, err := p.connection.Write([]byte{byte(pca9685Led0OnLReg + 4*channel), byte(on) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(pca9685Led0OnHReg + 4*channel), byte(on >> 8)}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(pca9685Led0OffLReg + 4*channel), byte(off) & 0xFF}); err != nil {
		return err
	}

	_, err := p.connection.Write([]byte{byte(pca9685Led0OffHReg + 4*channel), byte(off >> 8)})
	return err
}

// SetAllPWM sets all channels to a pwm value from 0-4095.
// Params:
//
//	on uint16 - the time to start the pulse
//	off uint16 - the time to stop the pulse
//
// Most typically you set "on" to a zero value, and then set "off" to your desired duty.
func (p *PCA9685Driver) SetAllPWM(on uint16, off uint16) error {
	if _, err := p.connection.Write([]byte{byte(pca9685AllLedOnLReg), byte(on) & 0xFF}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(pca9685AllLedOnHReg), byte(on >> 8)}); err != nil {
		return err
	}

	if _, err := p.connection.Write([]byte{byte(pca9685AllLedOffLReg), byte(off) & 0xFF}); err != nil {
		return err
	}

	_, err := p.connection.Write([]byte{byte(pca9685AllLedOffHReg), byte(off >> 8)})
	return err
}

// SetPWMFreq sets the PWM frequency in Hz between 24Hz and 1526Hz, the default is 200Hz.
func (p *PCA9685Driver) SetPWMFreq(freq float32) error {
	// internal IC oscillator frequency is 25 MHz
	var prescalevel float32 = 25000000
	// find frequency of PWM waveform
	prescalevel /= 4096
	// ratio between desired frequency and maximum
	prescalevel /= freq
	prescalevel--
	// round value to nearest whole
	prescale := byte(prescalevel + 0.5)

	if _, err := p.connection.Write([]byte{byte(pca9685Mode1Reg)}); err != nil {
		return err
	}
	oldmode, err := p.connection.ReadByte()
	if err != nil {
		return err
	}

	// put oscillator in sleep mode, clear the restart bit 7 here to prevent unneeded restart
	sleepMode := (oldmode &^ pca9685Mode1RegRestartBit) | pca9685Mode1RegSleepBit
	if _, err := p.connection.Write([]byte{byte(pca9685Mode1Reg), sleepMode}); err != nil {
		return err
	}
	// write prescaler value
	if _, err := p.connection.Write([]byte{byte(pca9685PrescaleReg), prescale}); err != nil {
		return err
	}
	// put back to old settings, ensure no sleep
	noSleepMode := oldmode &^ pca9685Mode1RegSleepBit
	if _, err := p.connection.Write([]byte{byte(pca9685Mode1Reg), noSleepMode}); err != nil {
		return err
	}

	// wait >500us according to data sheet
	time.Sleep(5 * time.Millisecond)

	// initiate a restart
	restartMode := oldmode | pca9685Mode1RegRestartBit
	_, err = p.connection.Write([]byte{byte(pca9685Mode1Reg), restartMode})
	return err
}

// PwmWrite writes a PWM signal to the specified channel aka "pin".
// Value values are from 0-255, to conform to the PwmWriter interface.
// If you need finer control, please look at SetPWM().
func (p *PCA9685Driver) PwmWrite(pin string, val byte) error {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 255), 0, 4095)
	return p.SetPWM(i, 0, uint16(v))
}

// ServoWrite writes a servo signal to the specified channel aka "pin".
// Valid values are from 0-180, to conform to the ServoWriter interface.
// If you need finer control, please look at SetPWM().
func (p *PCA9685Driver) ServoWrite(pin string, val byte) error {
	i, err := strconv.Atoi(pin)
	if err != nil {
		return err
	}
	v := gobot.ToScale(gobot.FromScale(float64(val), 0, 180), 200, 500)
	return p.SetPWM(i, 0, uint16(v))
}

// initialize the driver according to the data sheet section "7.3.1.1 Restart mode"
// * ensure the sleep bit is unset
// * wait > 500us
// * write a logic 1 to bit 7 (RESTART) of register "MODE1"
func (p *PCA9685Driver) initialize() error {
	if err := p.SetAllPWM(0, 0); err != nil {
		return err
	}

	// set not inverted (default), outputs change on stop (default), OE reaction to 0 (default), totem-pole (default)
	if _, err := p.connection.Write([]byte{pca9685Mode2Reg, pca9685Mode2RegOutdrvBit}); err != nil {
		return err
	}
	// reset of sleep bit together with set of no restart (default), internal clock (default), no AI (default),
	// no response to sub address 1, 2 or 3 (default), activate response to all-call (default)
	if _, err := p.connection.Write([]byte{pca9685Mode1Reg, pca9685Mode1RegAllCallBit}); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	// initiate the restart
	if _, err := p.connection.Write([]byte{byte(pca9685Mode1Reg)}); err != nil {
		return err
	}
	oldmode, err := p.connection.ReadByte()
	if err != nil {
		return err
	}
	oldmode = oldmode | byte(pca9685Mode1RegRestartBit)

	if _, err := p.connection.Write([]byte{pca9685Mode1Reg, oldmode}); err != nil {
		return err
	}

	time.Sleep(5 * time.Millisecond)

	return nil
}

func (p *PCA9685Driver) shutdown() error {
	_, err := p.connection.Write([]byte{pca9685AllLedOffHReg, pca9685AllLedOffHRegShutDown})
	return err
}
