package system

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

const pwmDebug = false

const (
	pwmPinErrorPattern    = "%s() failed for id %s with %v"   //nolint:gosec // false positive
	pwmPinSetErrorPattern = "%s(%v) failed for id %s with %v" //nolint:gosec // false positive
)

// pwmPinSysFs represents a PWM pin
type pwmPinSysFs struct {
	path                       string
	pin                        string
	polarityNormalIdentifier   string
	polarityInvertedIdentifier string

	sfa *sysfsFileAccess
}

// newPWMPinSysfs returns a new pwmPin, working with sysfs file access.
func newPWMPinSysfs(sfa *sysfsFileAccess, path string, pin int, polNormIdent string, polInvIdent string) *pwmPinSysFs {
	p := &pwmPinSysFs{
		path:                       path,
		pin:                        strconv.Itoa(pin),
		polarityNormalIdentifier:   polNormIdent,
		polarityInvertedIdentifier: polInvIdent,

		sfa: sfa,
	}
	return p
}

// Export exports the pin for use by the operating system by writing the pin to the "Export" path
func (p *pwmPinSysFs) Export() error {
	if err := p.sfa.write(p.pwmExportPath(), []byte(p.pin)); err != nil {
		// If EBUSY then the pin has already been exported, we suppress the error
		var pathError *os.PathError
		if !(errors.As(err, &pathError) && errors.Is(err, Syscall_EBUSY)) {
			return fmt.Errorf(pwmPinErrorPattern, "Export", p.pin, err)
		}
	}

	// Pause to avoid race condition in case there is any udev rule
	// that changes file permissions on newly exported PWMPin. This
	// is a common circumstance when running as a non-root user.
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Unexport releases the pin from the operating system by writing the pin to the "Unexport" path
func (p *pwmPinSysFs) Unexport() error {
	if err := p.sfa.write(p.pwmUnexportPath(), []byte(p.pin)); err != nil {
		return fmt.Errorf(pwmPinErrorPattern, "Unexport", p.pin, err)
	}
	return nil
}

// Enabled reads and returns the enabled state of the pin
func (p *pwmPinSysFs) Enabled() (bool, error) {
	val, err := p.sfa.readInteger(p.pwmEnablePath())
	if err != nil {
		return false, fmt.Errorf(pwmPinErrorPattern, "Enabled", p.pin, err)
	}

	return val > 0, nil
}

// SetEnabled writes enable(1) or disable(0) status. For most platforms this is only possible if period was
// already set to > 0. Regardless of setting to enable or disable, a "write error: Invalid argument" will occur.
func (p *pwmPinSysFs) SetEnabled(enable bool) error {
	enableVal := 0
	if enable {
		enableVal = 1
	}
	if err := p.sfa.writeInteger(p.pwmEnablePath(), enableVal); err != nil {
		if pwmDebug {
			p.printState()
		}
		return fmt.Errorf(pwmPinSetErrorPattern, "SetEnabled", enable, p.pin, err)
	}

	return nil
}

// Polarity reads and returns false if the polarity is inverted, otherwise true
func (p *pwmPinSysFs) Polarity() (bool, error) {
	buf, err := p.sfa.read(p.pwmPolarityPath())
	if err != nil {
		return true, fmt.Errorf(pwmPinErrorPattern, "Polarity", p.pin, err)
	}
	if len(buf) == 0 {
		return true, nil
	}

	ps := string(bytes.TrimRight(buf, "\n"))
	if ps == p.polarityNormalIdentifier {
		return true, nil
	}
	if ps == p.polarityInvertedIdentifier {
		return false, nil
	}

	return true, fmt.Errorf("unknown value (%s) in Polarity", ps)
}

// SetPolarity writes the polarity as normal if called with true and as inverted if called with false
func (p *pwmPinSysFs) SetPolarity(normal bool) error {
	enabled, _ := p.Enabled()
	if enabled {
		return fmt.Errorf("Cannot set PWM polarity when enabled")
	}
	value := p.polarityNormalIdentifier
	if !normal {
		value = p.polarityInvertedIdentifier
	}
	if err := p.sfa.write(p.pwmPolarityPath(), []byte(value)); err != nil {
		if pwmDebug {
			p.printState()
		}
		return fmt.Errorf(pwmPinSetErrorPattern, "SetPolarity", value, p.pin, err)
	}
	return nil
}

// Period returns the current period in nanoseconds
func (p *pwmPinSysFs) Period() (uint32, error) {
	val, err := p.sfa.readInteger(p.pwmPeriodPath())
	if err != nil {
		return 0, fmt.Errorf(pwmPinErrorPattern, "Period", p.pin, err)
	}
	//nolint:gosec // TODO: fix later
	return uint32(val), nil
}

// SetPeriod writes the current period in nanoseconds
func (p *pwmPinSysFs) SetPeriod(period uint32) error {
	if err := p.sfa.writeInteger(p.pwmPeriodPath(), int(period)); err != nil {

		if pwmDebug {
			p.printState()
		}
		return fmt.Errorf(pwmPinSetErrorPattern, "SetPeriod", period, p.pin, err)
	}
	return nil
}

// DutyCycle reads and returns the duty cycle in nanoseconds
func (p *pwmPinSysFs) DutyCycle() (uint32, error) {
	val, err := p.sfa.readInteger(p.pwmDutyCyclePath())
	if err != nil {
		return 0, fmt.Errorf(pwmPinErrorPattern, "DutyCycle", p.pin, err)
	}
	//nolint:gosec // TODO: fix later
	return uint32(val), err
}

// SetDutyCycle writes the duty cycle in nanoseconds
func (p *pwmPinSysFs) SetDutyCycle(duty uint32) error {
	if err := p.sfa.writeInteger(p.pwmDutyCyclePath(), int(duty)); err != nil {
		if pwmDebug {
			p.printState()
		}
		return fmt.Errorf(pwmPinSetErrorPattern, "SetDutyCycle", duty, p.pin, err)
	}
	return nil
}

// pwmExportPath returns export path
func (p *pwmPinSysFs) pwmExportPath() string {
	return path.Join(p.path, "export")
}

// pwmUnexportPath returns unexport path
func (p *pwmPinSysFs) pwmUnexportPath() string {
	return path.Join(p.path, "unexport")
}

// pwmDutyCyclePath returns duty_cycle path for specified pin
func (p *pwmPinSysFs) pwmDutyCyclePath() string {
	return path.Join(p.path, "pwm"+p.pin, "duty_cycle")
}

// pwmPeriodPath returns period path for specified pin
func (p *pwmPinSysFs) pwmPeriodPath() string {
	return path.Join(p.path, "pwm"+p.pin, "period")
}

// pwmEnablePath returns enable path for specified pin
func (p *pwmPinSysFs) pwmEnablePath() string {
	return path.Join(p.path, "pwm"+p.pin, "enable")
}

// pwmPolarityPath returns polarity path for specified pin
func (p *pwmPinSysFs) pwmPolarityPath() string {
	return path.Join(p.path, "pwm"+p.pin, "polarity")
}

func (p *pwmPinSysFs) printState() {
	enabled, _ := p.Enabled()
	polarity, _ := p.Polarity()
	period, _ := p.Period()
	duty, _ := p.DutyCycle()

	fmt.Println("Print state of all PWM variables...")
	fmt.Printf("Enable: %v, ", enabled)
	fmt.Printf("Polarity: %v, ", polarity)
	fmt.Printf("Period: %v, ", period)
	fmt.Printf("DutyCycle: %v, ", duty)
	var powerPercent float64
	if enabled {
		if polarity {
			powerPercent = float64(duty) / float64(period) * 100
		} else {
			powerPercent = float64(period) / float64(duty) * 100
		}
	}
	fmt.Printf("Power: %.1f\n", powerPercent)
}
