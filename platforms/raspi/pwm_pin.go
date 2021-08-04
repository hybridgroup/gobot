package raspi

import (
	"errors"
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/sysfs"
)

// PWMPin is the Raspberry Pi implementation of the PWMPinner interface.
// It uses Pi Blaster.
type PWMPin struct {
	pin    string
	dc     uint32
	period uint32
}

// NewPwmPin returns a new PWMPin
func NewPWMPin(pin string) *PWMPin {
	return &PWMPin{
		pin: pin}
}

// Export exports the pin for use by the Raspberry Pi
func (p *PWMPin) Export() error {
	return nil
}

// Unexport unexports the pin and releases the pin from the operating system
func (p *PWMPin) Unexport() error {
	return p.piBlaster(fmt.Sprintf("release %v\n", p.pin))
}

// Enable enables/disables the PWM pin
func (p *PWMPin) Enable(e bool) (err error) {
	return nil
}

// Polarity returns the polarity either normal or inverted
func (p *PWMPin) Polarity() (polarity string, err error) {
	return "normal", nil
}

// InvertPolarity does not do anything when using PiBlaster
func (p *PWMPin) InvertPolarity(invert bool) (err error) {
	return nil
}

// Period returns the current PWM period for pin
func (p *PWMPin) Period() (period uint32, err error) {
	if p.period == 0 {
		return p.period, errors.New("Raspi PWM pin period not set")
	}

	return p.period, nil
}

// SetPeriod uses PiBlaster setting and cannot be changed once set
func (p *PWMPin) SetPeriod(period uint32) (err error) {
	if p.period != 0 {
		return errors.New("Cannot set the period of individual PWM pins on Raspi")
	}
	p.period = period
	return nil
}

// DutyCycle returns the duty cycle for the pin
func (p *PWMPin) DutyCycle() (duty uint32, err error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *PWMPin) SetDutyCycle(duty uint32) (err error) {
	if p.period == 0 {
		return errors.New("Raspi PWM pin period not set")
	}

	if duty > p.period {
		return errors.New("Duty cycle exceeds period.")
	}
	p.dc = duty

	val := gobot.FromScale(float64(p.dc), 0, float64(p.period))
	
	// never go below minimum allowed duty for pi blaster
	// unless the duty equals to 0
	if val < 0.05 && val != 0 {
		val = 0.05
	}
	return p.piBlaster(fmt.Sprintf("%v=%v\n", p.pin, val))
}

func (p *PWMPin) piBlaster(data string) (err error) {
	fi, err := sysfs.OpenFile("/dev/pi-blaster", os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(data)
	return
}
