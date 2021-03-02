package raspi

import (
	"errors"
	"fmt"
	//"os"
	"os/exec"

	"gobot.io/x/gobot"
	//"gobot.io/x/gobot/sysfs"
)

type PigPWM struct {
	pin    string
	dc     uint32
	period uint32
}

// NewPwmPin returns a new PigPWM
func NewPigPWM(pin string) *PigPWM {
	return &PigPWM{
		pin: pin}
}

// Export exports the pin for use by the Raspberry Pi
func (p *PigPWM) Export() error {
	return nil
}

// Unexport unexports the pin and releases the pin from the operating system
func (p *PigPWM) Unexport() error {
	return nil
}

// Enable enables/disables the PWM pin
func (p *PigPWM) Enable(e bool) (err error) {
	return nil
}

// Polarity returns the polarity either normal or inverted
func (p *PigPWM) Polarity() (polarity string, err error) {
	return "normal", nil
}

// InvertPolarity does not do anything when using PiBlaster
func (p *PigPWM) InvertPolarity(invert bool) (err error) {
	return nil
}

// Period returns the current PWM period for pin
func (p *PigPWM) Period() (period uint32, err error) {
	if p.period == 0 {
		return p.period, errors.New("Raspi PWM pin period not set")
	}

	return p.period, nil
}

// SetPeriod uses PiBlaster setting and cannot be changed once set
func (p *PigPWM) SetPeriod(period uint32) (err error) {
	if p.period != 0 {
		return errors.New("Cannot set the period of individual PWM pins on Raspi")
	}
	p.period = period
	return nil
}

// DutyCycle returns the duty cycle for the pin
func (p *PigPWM) DutyCycle() (duty uint32, err error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *PigPWM) SetDutyCycle(duty uint32) (err error) {
	if p.period == 0 {
		return errors.New("Raspi PWM pin period not set")
	}

	if duty > p.period {
		return errors.New("Duty cycle exceeds period.")
	}

	p.dc = duty

	val := gobot.FromScale(float64(p.dc), 0, float64(p.period))
	if val < .05 {
		val = .05
	}

	byteValue := int(255 * val)

	return pigs("p", p.pin, fmt.Sprintf("%d", byteValue))
}

func pigs(args ...string) (err error) {
	fmt.Printf("hi %v\n", args)
	cmd := exec.Command("pigs", args...)
	output, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	return fmt.Errorf("error running pigs command: %v %v", err, string(output))
}
