package jetson

import (
	"errors"
	"fmt"
	"os"
	"path"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/system"
)

const (
	minimumPeriod = 5334
	minimumRate   = 0.05
)

// PWMPin is the Jetson Nano implementation of the PWMPinner interface.
// It uses gpio pwm.
type PWMPin struct {
	sys     *system.Accesser
	path    string
	fn      string
	dc      uint32
	period  uint32
	enabled bool
}

// NewPWMPin returns a new PWMPin
// pin32 pwm0, pin33 pwm2
func NewPWMPin(sys *system.Accesser, path string, fn string) *PWMPin {
	p := &PWMPin{
		sys:  sys,
		path: path,
		fn:   fn,
	}
	return p
}

// Export exports the pin for use by the Jetson Nano
func (p *PWMPin) Export() error {
	return p.writeFile("export", p.fn)
}

// Unexport releases the pin from the operating system
func (p *PWMPin) Unexport() error {
	return p.writeFile("unexport", p.fn)
}

// Enabled returns the cached enabled state of the PWM pin
func (p *PWMPin) Enabled() (bool, error) {
	return p.enabled, nil
}

// SetEnabled enables/disables the PWM pin
func (p *PWMPin) SetEnabled(e bool) error {
	if err := p.writeFile(fmt.Sprintf("pwm%s/enable", p.fn), fmt.Sprintf("%v", bool2int(e))); err != nil {
		return err
	}
	p.enabled = e
	return nil
}

// Polarity returns always the polarity "true" for normal
func (p *PWMPin) Polarity() (bool, error) {
	return true, nil
}

// SetPolarity does not do anything when using Jetson Nano
func (p *PWMPin) SetPolarity(bool) error {
	return nil
}

// Period returns the cached PWM period for pin
func (p *PWMPin) Period() (period uint32, err error) {
	if p.period == 0 {
		return p.period, errors.New("Jetson PWM pin period not set")
	}

	return p.period, nil
}

// SetPeriod uses Jetson Nano setting and cannot be changed once set
func (p *PWMPin) SetPeriod(period uint32) error {
	if p.period != 0 {
		return errors.New("Cannot set the period of individual PWM pins on Jetson")
	}
	// JetsonNano Minimum period
	if period < minimumPeriod {
		return errors.New("Cannot set the period more then minimum")
	}
	if err := p.writeFile(fmt.Sprintf("pwm%s/period", p.fn), fmt.Sprintf("%v", period)); err != nil {
		return err
	}
	p.period = period
	return nil
}

// DutyCycle returns the cached duty cycle for the pin
func (p *PWMPin) DutyCycle() (uint32, error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *PWMPin) SetDutyCycle(duty uint32) error {
	if p.period == 0 {
		return errors.New("Jetson PWM pin period not set")
	}

	if duty > p.period {
		return errors.New("Duty cycle exceeds period")
	}

	rate := gobot.FromScale(float64(duty), 0, float64(p.period))
	// never go below minimum allowed duty because very short duty
	if rate < minimumRate {
		duty = uint32(minimumRate * float64(p.period) / 100)
	}
	if err := p.writeFile(fmt.Sprintf("pwm%s/duty_cycle", p.fn), fmt.Sprintf("%v", duty)); err != nil {
		return err
	}

	p.dc = duty
	return nil
}

func (p *PWMPin) writeFile(subpath string, value string) error {
	sysfspath := path.Join(p.path, subpath)
	fi, err := p.sys.OpenFile(sysfspath, os.O_WRONLY|os.O_APPEND, 0o644)
	defer fi.Close() //nolint:staticcheck // for historical reasons

	if err != nil {
		return err
	}

	_, err = fi.WriteString(value)
	return err
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
