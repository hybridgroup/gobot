package raspi

import (
	"errors"
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

// PWMPin is the Raspberry Pi implementation of the PWMPinner interface.
// It uses Pi Blaster.
type PWMPin struct {
	sys    *system.Accesser
	path   string
	pin    string
	dc     uint32
	period uint32
}

// NewPWMPin returns a new PWMPin
func NewPWMPin(sys *system.Accesser, path string, pin string) *PWMPin {
	return &PWMPin{
		sys:  sys,
		path: path,
		pin:  pin,
	}
}

// Export exports the pin for use by the Raspberry Pi
func (p *PWMPin) Export() error {
	return nil
}

// Unexport releases the pin from the operating system
func (p *PWMPin) Unexport() error {
	return p.writeValue(fmt.Sprintf("release %v\n", p.pin))
}

// Enabled returns always true for "enabled"
func (p *PWMPin) Enabled() (bool, error) {
	return true, nil
}

// SetEnabled do nothing for PiBlaster
func (p *PWMPin) SetEnabled(e bool) error {
	return nil
}

// Polarity returns always true for "normal"
func (p *PWMPin) Polarity() (bool, error) {
	return true, nil
}

// SetPolarity does not do anything when using PiBlaster
func (p *PWMPin) SetPolarity(bool) (err error) {
	return nil
}

// Period returns the cached PWM period for pin
func (p *PWMPin) Period() (uint32, error) {
	if p.period == 0 {
		return p.period, errors.New("Raspi PWM pin period not set")
	}

	return p.period, nil
}

// SetPeriod uses PiBlaster setting and cannot be changed once set
func (p *PWMPin) SetPeriod(period uint32) error {
	if p.period != 0 {
		return errors.New("Cannot set the period of individual PWM pins on Raspi")
	}
	p.period = period
	return nil
}

// DutyCycle returns the duty cycle for the pin
func (p *PWMPin) DutyCycle() (uint32, error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *PWMPin) SetDutyCycle(duty uint32) error {
	if p.period == 0 {
		return errors.New("Raspi PWM pin period not set")
	}

	if duty > p.period {
		return errors.New("Duty cycle exceeds period")
	}

	val := gobot.FromScale(float64(duty), 0, float64(p.period))
	// never go below minimum allowed duty for pi blaster
	// unless the duty equals to 0
	if val < 0.05 && val != 0 {
		val = 0.05
	}

	if err := p.writeValue(fmt.Sprintf("%v=%v\n", p.pin, val)); err != nil {
		return err
	}

	p.dc = duty
	return nil
}

func (p *PWMPin) writeValue(data string) (err error) {
	fi, err := p.sys.OpenFile(p.path, os.O_WRONLY|os.O_APPEND, 0o644)
	defer fi.Close() //nolint:staticcheck // for historical reasons

	if err != nil {
		return err
	}

	_, err = fi.WriteString(data)
	return
}
