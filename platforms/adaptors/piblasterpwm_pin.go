package adaptors

import (
	"fmt"
	"os"
	"strconv"

	"gobot.io/x/gobot/v2/system"
)

const (
	piBlasterPath        = "/dev/pi-blaster"
	piBlasterMinDutyNano = 10000 // 10 us
)

// piBlasterPWMPin is the Raspberry Pi implementation of the PWMPinner interface.
// It uses Pi Blaster.
type piBlasterPWMPin struct {
	sys    *system.Accesser
	pin    string
	dc     uint32
	period uint32
}

// newPiBlasterPWMPin returns a new PWM pin for pi-blaster access.
func newPiBlasterPWMPin(sys *system.Accesser, pinNo int) *piBlasterPWMPin {
	return &piBlasterPWMPin{
		sys: sys,
		pin: strconv.Itoa(pinNo),
	}
}

// Export exports the pin for use by the Raspberry Pi
func (p *piBlasterPWMPin) Export() error {
	return nil
}

// Unexport releases the pin from the operating system
func (p *piBlasterPWMPin) Unexport() error {
	return p.writeValue(fmt.Sprintf("release %v\n", p.pin))
}

// Enabled returns always true for "enabled"
func (p *piBlasterPWMPin) Enabled() (bool, error) {
	return true, nil
}

// SetEnabled do nothing for PiBlaster
func (p *piBlasterPWMPin) SetEnabled(e bool) error {
	return nil
}

// Polarity returns always true for "normal"
func (p *piBlasterPWMPin) Polarity() (bool, error) {
	return true, nil
}

// SetPolarity does not do anything when using PiBlaster
func (p *piBlasterPWMPin) SetPolarity(bool) error {
	return nil
}

// Period returns the cached PWM period for pin
func (p *piBlasterPWMPin) Period() (uint32, error) {
	return p.period, nil
}

// SetPeriod uses PiBlaster setting and cannot be changed. We allow setting once here to define a base period for
// ServoWrite(). see https://github.com/sarfata/pi-blaster#how-to-adjust-the-frequency-and-the-resolution-of-the-pwm
func (p *piBlasterPWMPin) SetPeriod(period uint32) error {
	if p.period != 0 {
		return fmt.Errorf("the period of PWM pins needs to be set to '%d' in pi-blaster source code", period)
	}
	p.period = period
	return nil
}

// DutyCycle returns the duty cycle for the pin
func (p *piBlasterPWMPin) DutyCycle() (uint32, error) {
	return p.dc, nil
}

// SetDutyCycle writes the duty cycle to the pin
func (p *piBlasterPWMPin) SetDutyCycle(dutyNanos uint32) error {
	if p.period == 0 {
		return fmt.Errorf("pi-blaster PWM pin period not set while try to set duty cycle to '%d'", dutyNanos)
	}

	if dutyNanos > p.period {
		return fmt.Errorf("the duty cycle (%d) exceeds period (%d) for pi-blaster", dutyNanos, p.period)
	}

	// never go below minimum allowed duty for pi blaster unless the duty equals to 0
	if dutyNanos < piBlasterMinDutyNano && dutyNanos != 0 {
		dutyNanos = piBlasterMinDutyNano
		fmt.Printf("duty cycle value limited to '%d' ns for pi-blaster", dutyNanos)
	}

	duty := float64(dutyNanos) / float64(p.period)
	if err := p.writeValue(fmt.Sprintf("%v=%v\n", p.pin, duty)); err != nil {
		return err
	}

	p.dc = dutyNanos
	return nil
}

func (p *piBlasterPWMPin) writeValue(data string) error {
	fi, err := p.sys.OpenFile(piBlasterPath, os.O_WRONLY|os.O_APPEND, 0o644)
	defer fi.Close() //nolint:staticcheck // for historical reasons

	if err != nil {
		return err
	}

	_, err = fi.WriteString(data)
	return err
}
