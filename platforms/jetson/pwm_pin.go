package jetson

import (
	"errors"
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/sysfs"
)

const (
	minimumPeriod = 5334
	minimumRate   = 0.05
)

// PWMPin is the Jetson Nano implementation of the PWMPinner interface.
// It uses gpio pwm.
type PWMPin struct {
	pin    string
	fn     string
	dc     uint32
	period uint32
}

// NewPwmPin returns a new PWMPin
// pin32 pwm0, pin33 pwm2
func NewPWMPin(pin string) (p *PWMPin, err error) {
	if val, ok := pwms[pin]; ok {
		p = &PWMPin{pin: pin, fn: val}
	} else {
		err = errors.New("Not a valid pin")
	}

	return
}

// Export exports the pin for use by the Jetson Nano
func (p *PWMPin) Export() error {
	fi, err := sysfs.OpenFile("/sys/class/pwm/pwmchip0/export", os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(p.fn)

	return err
}

// Unexport unexports the pin and releases the pin from the operating system
func (p *PWMPin) Unexport() error {
	fi, err := sysfs.OpenFile("/sys/class/pwm/pwmchip0/unexport", os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(p.fn)

	return err
}

// Enable enables/disables the PWM pin
func (p *PWMPin) Enable(e bool) (err error) {
	fi, err := sysfs.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%s/enable", p.fn), os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(fmt.Sprintf("%v", bool2int(e)))

	return err
}

// Polarity returns the polarity either normal or inverted
func (p *PWMPin) Polarity() (polarity string, err error) {
	return "normal", nil
}

// SetPolarity does not do anything when using Jetson Nano
func (p *PWMPin) SetPolarity(value string) (err error) {
	return nil
}

// InvertPolarity does not do anything when using Jetson Nano
func (p *PWMPin) InvertPolarity(invert bool) (err error) {
	return nil
}

// Period returns the current PWM period for pin
func (p *PWMPin) Period() (period uint32, err error) {
	if p.period == 0 {
		return p.period, errors.New("Jetson PWM pin period not set")
	}

	return p.period, nil
}

// SetPeriod uses Jetson Nano setting and cannot be changed once set
func (p *PWMPin) SetPeriod(period uint32) (err error) {
	if p.period != 0 {
		return errors.New("Cannot set the period of individual PWM pins on Jetson")
	}
	// JetsonNano Minimum period
	if period < minimumPeriod {
		return errors.New("Cannot set the period more Then minimum.")
	}

	p.period = period
	fi, err := sysfs.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%s/period", p.fn), os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(fmt.Sprintf("%v", p.period))

	return nil
}

// DutyCycle returns the duty cycle for the pin
func (p *PWMPin) DutyCycle() (duty uint32, err error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *PWMPin) SetDutyCycle(duty uint32) (err error) {
	if p.period == 0 {
		return errors.New("Jetson PWM pin period not set")
	}

	if duty > p.period {
		return errors.New("Duty cycle exceeds period.")
	}
	p.dc = duty

	rate := gobot.FromScale(float64(p.dc), 0, float64(p.period))

	// never go below minimum allowed duty becuse very short duty
	if rate < minimumRate {
		duty = uint32(minimumRate * float64(p.period) / 100)
		p.dc = duty
	}

	fi, err := sysfs.OpenFile(fmt.Sprintf("/sys/class/pwm/pwmchip0/pwm%s/duty_cycle", p.fn), os.O_WRONLY|os.O_APPEND, 0644)
	defer fi.Close()

	if err != nil {
		return err
	}

	_, err = fi.WriteString(fmt.Sprintf("%v", duty))

	return
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
