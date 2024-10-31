package i2c

import (
	"gobot.io/x/gobot/v2"
)

// Adafruit2327Driver is a driver for Adafruit 16-Channel PWM/Servo HAT & Bonnet - a Raspberry Pi add-on, based on
// PCA9685. This driver just wraps the PCA9685Driver.
// Stacking 62 of them is possible (addresses 0x40..0x7E), for controlling up to 992 servos.
// datasheet:
// https://cdn-learn.adafruit.com/downloads/pdf/adafruit-16-channel-pwm-servo-hat-for-raspberry-pi.pdf
type Adafruit2327Driver struct {
	*PCA9685Driver
}

// NewAdafruit2327Driver initializes a new driver for PWM servos.
// Params:
//
//	c Connector - the Adaptor to use with this Driver
//
// Optional params:
//
//	i2c.WithBus(int):	    bus to use with this driver
//	i2c.WithAddress(int):	address to use with this driver
func NewAdafruit2327Driver(c Connector, options ...func(Config)) *Adafruit2327Driver {
	pca := NewPCA9685Driver(c, options...) // the default address of the driver is already 0x40
	pca.SetName(gobot.DefaultName("Adafruit2327ServoHat"))
	d := &Adafruit2327Driver{
		PCA9685Driver: pca,
	}

	// TODO: add API funcs
	return d
}

// SetServoMotorFreq sets the frequency for the currently addressed PWM Servo HAT.
func (a *Adafruit2327Driver) SetServoMotorFreq(freq float64) error {
	return a.SetPWMFreq(float32(freq))
}

// SetServoMotorPulse is a convenience function to specify the 'tick' value,
// between 0-4095, when the signal will turn on, and when it will turn off.
func (a *Adafruit2327Driver) SetServoMotorPulse(channel byte, on, off int32) error {
	return a.SetPWM(int(channel), uint16(on), uint16(off)) //nolint:gosec // TODO: fix later
}
