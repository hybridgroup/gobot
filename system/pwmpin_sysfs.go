package system

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"
)

const (
	pwmPinErrorPattern    = "%s() failed for pin %s with %v"
	pwmPinSetErrorPattern = "%s(%v) failed for pin %s with %v"
)

const (
	polarityNormal   = "normal"
	polarityInverted = "inverted"
)

const pwmDebug = false

// PWMPinner is the interface for system PWM interactions
type PWMPinner interface {
	// Export exports the pin for use by the operating system
	Export() error
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Enable enables/disables the PWM pin
	// TODO: rename to "SetEnable(bool)" according to golang style and allow "Enable()" to be the getter function
	Enable(bool) (err error)
	// Polarity returns the polarity either normal or inverted
	Polarity() (polarity string, err error)
	// SetPolarity writes value to pwm polarity path
	SetPolarity(value string) (err error)
	// InvertPolarity sets the polarity to inverted if called with true
	InvertPolarity(invert bool) (err error)
	// Period returns the current PWM period for pin
	Period() (period uint32, err error)
	// SetPeriod sets the current PWM period for pin
	SetPeriod(period uint32) (err error)
	// DutyCycle returns the duty cycle for the pin
	DutyCycle() (duty uint32, err error)
	// SetDutyCycle writes the duty cycle to the pin
	SetDutyCycle(duty uint32) (err error)
}

// PWMPinnerProvider is the interface that an Adaptor should implement to allow
// clients to obtain access to any PWMPin's available on that board.
type PWMPinnerProvider interface {
	PWMPin(string) (PWMPinner, error)
}

// PWMPin represents a PWM pin
type PWMPin struct {
	pin     string
	path    string
	enabled bool
	write   func(fs filesystem, path string, data []byte) (i int, err error)
	read    func(fs filesystem, path string) ([]byte, error)
	fs      filesystem
}

// NewPWMPin returns a new pwmPin
func (a *Accesser) NewPWMPin(path string, pin int) *PWMPin {
	return &PWMPin{
		pin:     strconv.Itoa(pin),
		enabled: false,
		path:    path, //"/sys/class/pwm/pwmchip0"
		read:    readPwmFile,
		write:   writePwmFile,
		fs:      a.fs,
	}
}

// Export writes pin to pwm export path
func (p *PWMPin) Export() error {
	_, err := p.write(p.fs, p.pwmExportPath(), []byte(p.pin))
	if err != nil {
		// If EBUSY then the pin has already been exported
		e, ok := err.(*os.PathError)
		if !ok || e.Err != syscall.EBUSY {
			return fmt.Errorf(pwmPinErrorPattern, "Export", p.pin, err)
		}
	}

	// Pause to avoid race condition in case there is any udev rule
	// that changes file permissions on newly exported PWMPin. This
	// is a common circumstance when running as a non-root user.
	time.Sleep(100 * time.Millisecond)

	return nil
}

// Unexport writes pin to pwm unexport path
func (p *PWMPin) Unexport() (err error) {
	if _, err = p.write(p.fs, p.pwmUnexportPath(), []byte(p.pin)); err != nil {
		err = fmt.Errorf(pwmPinErrorPattern, "Unexport", p.pin, err)
	}
	return
}

// enable returns current enable value
func (p *PWMPin) enable() (enabled bool, err error) {
	buf, err := p.read(p.fs, p.pwmEnablePath())
	if err != nil {
		return enabled, fmt.Errorf(pwmPinErrorPattern, "enable", p.pin, err)
	}
	if len(buf) == 0 {
		return enabled, nil
	}

	v := bytes.TrimRight(buf, "\n")
	val, e := strconv.Atoi(string(v))
	return val > 0, e
}

// Enable writes value to pwm enable path
func (p *PWMPin) Enable(enable bool) (err error) {
	if p.enabled != enable {
		p.enabled = enable
		enableVal := 0
		if enable {
			enableVal = 1
		}
		if _, err = p.write(p.fs, p.pwmEnablePath(), []byte(fmt.Sprintf("%v", enableVal))); err != nil {
			err = fmt.Errorf(pwmPinSetErrorPattern, "Enable", enable, p.pin, err)
			if pwmDebug {
				p.printState()
			}
		}
	}
	return
}

// Polarity returns current polarity value
func (p *PWMPin) Polarity() (polarity string, err error) {
	buf, err := p.read(p.fs, p.pwmPolarityPath())
	if err != nil {
		return polarity, fmt.Errorf(pwmPinErrorPattern, "Polarity", p.pin, err)
	}
	if len(buf) == 0 {
		return "", nil
	}

	return string(bytes.TrimRight(buf, "\n")), nil
}

// SetPolarity writes value to pwm polarity path
func (p *PWMPin) SetPolarity(value string) (err error) {
	if p.enabled {
		return fmt.Errorf("Cannot set PWM polarity when enabled")
	}
	if _, err = p.write(p.fs, p.pwmPolarityPath(), []byte(value)); err != nil {
		err = fmt.Errorf(pwmPinSetErrorPattern, "SetPolarity", value, p.pin, err)
		if pwmDebug {
			p.printState()
		}
	}
	return
}

// InvertPolarity sets the polarity to "inverted" when 'true' is given, otherwise to "normal"
func (p *PWMPin) InvertPolarity(invert bool) (err error) {
	polarity := polarityNormal
	if invert {
		polarity = polarityInverted
	}
	return p.SetPolarity(polarity)
}

// Period reads from pwm period path and returns value in nanoseconds
func (p *PWMPin) Period() (period uint32, err error) {
	buf, err := p.read(p.fs, p.pwmPeriodPath())
	if err != nil {
		return period, fmt.Errorf(pwmPinErrorPattern, "Period", p.pin, err)
	}
	if len(buf) == 0 {
		return 0, nil
	}

	v := bytes.TrimRight(buf, "\n")
	val, e := strconv.Atoi(string(v))
	return uint32(val), e
}

// SetPeriod sets pwm period in nanoseconds
func (p *PWMPin) SetPeriod(period uint32) (err error) {
	if _, err = p.write(p.fs, p.pwmPeriodPath(), []byte(fmt.Sprintf("%v", period))); err != nil {
		err = fmt.Errorf(pwmPinSetErrorPattern, "SetPeriod", period, p.pin, err)
		if pwmDebug {
			p.printState()
		}
	}
	return
}

// DutyCycle reads from pwm duty cycle path and returns value in nanoseconds
func (p *PWMPin) DutyCycle() (duty uint32, err error) {
	buf, err := p.read(p.fs, p.pwmDutyCyclePath())
	if err != nil {
		return duty, fmt.Errorf(pwmPinErrorPattern, "DutyCycle", p.pin, err)
	}

	v := bytes.TrimRight(buf, "\n")
	val, e := strconv.Atoi(string(v))
	return uint32(val), e
}

// SetDutyCycle writes value to pwm duty cycle path
// duty is in nanoseconds
func (p *PWMPin) SetDutyCycle(duty uint32) (err error) {
	if _, err = p.write(p.fs, p.pwmDutyCyclePath(), []byte(fmt.Sprintf("%v", duty))); err != nil {
		err = fmt.Errorf(pwmPinSetErrorPattern, "SetDutyCycle", duty, p.pin, err)
		if pwmDebug {
			p.printState()
		}
	}
	return
}

// pwmExportPath returns export path
func (p *PWMPin) pwmExportPath() string {
	return path.Join(p.path, "export")
}

// pwmUnexportPath returns unexport path
func (p *PWMPin) pwmUnexportPath() string {
	return path.Join(p.path, "unexport")
}

// pwmDutyCyclePath returns duty_cycle path for specified pin
func (p *PWMPin) pwmDutyCyclePath() string {
	return path.Join(p.path, "pwm"+p.pin, "duty_cycle")
}

// pwmPeriodPath returns period path for specified pin
func (p *PWMPin) pwmPeriodPath() string {
	return path.Join(p.path, "pwm"+p.pin, "period")
}

// pwmEnablePath returns enable path for specified pin
func (p *PWMPin) pwmEnablePath() string {
	return path.Join(p.path, "pwm"+p.pin, "enable")
}

// pwmPolarityPath returns polarity path for specified pin
func (p *PWMPin) pwmPolarityPath() string {
	return path.Join(p.path, "pwm"+p.pin, "polarity")
}

func writePwmFile(fs filesystem, path string, data []byte) (i int, err error) {
	file, err := fs.openFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readPwmFile(fs filesystem, path string) ([]byte, error) {
	file, err := fs.openFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i int
	i, err = file.Read(buf)
	if i == 0 {
		return []byte{}, err
	}
	return buf[:i], err
}

func (p *PWMPin) printState() {
	enabled, _ := p.enable()
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
		if polarity == polarityNormal {
			powerPercent = float64(duty) / float64(period) * 100
		} else {
			powerPercent = float64(period) / float64(duty) * 100
		}
	}
	fmt.Printf("Power: %.1f\n", powerPercent)
}
