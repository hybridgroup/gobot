package sysfs

import (
	"fmt"
	"os"
	"strconv"
)

// PWMPin is the interface for sysfs PWM interactions
type PWMPinner interface {
	// Export exports the pin for use by the operating system
	Export() error
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Enable enables/disables the PWM pin
	Enable(bool) (err error)
	// Period returns the current PWM period for pin
	Period() (period string, err error)
	// SetPeriod sets the current PWM period for pin
	SetPeriod(period string) (err error)
	// DutyCycle returns the duty cycle for the pin
	DutyCycle() (duty float64, err error)
	// SetDutyCycle writes the duty cycle to the pin
	SetDutyCycle(duty float64) (err error)
}

type PWMPin struct {
	pin     string
	Chip    string
	enabled bool
	write   func(path string, data []byte) (i int, err error)
	read    func(path string) ([]byte, error)
}

// NewPwmPin returns a new pwmPin
func NewPWMPin(pin int) *PWMPin {
	return &PWMPin{
		pin:     strconv.Itoa(pin),
		enabled: false,
		Chip:    "0",
		read:    readPwmFile,
		write:   writePwmFile}
}

// Export writes pin to pwm export path
func (p *PWMPin) Export() (err error) {
	_, err = p.write(p.pwmExportPath(), []byte(p.pin))
	return
}

// Unexport writes pin to pwm unexport path
func (p *PWMPin) Unexport() (err error) {
	_, err = p.write(p.pwmUnexportPath(), []byte(p.pin))
	return
}

// Enable writes value to pwm enable path
func (p *PWMPin) Enable(enable bool) (err error) {
	if p.enabled != enable {
		p.enabled = enable
		enableVal := 0
		if enable {
			enableVal = 1
		}
		_, err = p.write(p.pwmEnablePath(), []byte(fmt.Sprintf("%v", enableVal)))
	}
	return
}

// SetPolarityInverted writes value to pwm polarity path
func (p *PWMPin) SetPolarityInverted(invert bool) (err error) {
	if p.enabled {
		polarity := "normal"
		if invert {
			polarity = "inverted"
		}
		_, err = p.write(p.pwmPolarityPath(), []byte(polarity))
	}
	return
}

// Period reads from pwm period path and returns value
func (p *PWMPin) Period() (period string, err error) {
	buf, err := p.read(p.pwmPeriodPath())
	if err != nil {
		return
	}
	return string(buf), nil
}

// SetPeriod sets pwm period in nanoseconds
func (p *PWMPin) SetPeriod(period uint32) (err error) {
	_, err = p.write(p.pwmPeriodPath(), []byte(fmt.Sprintf("%v", period)))
	return
}

// DutyCycle reads from pwm duty cycle path and returns value
func (p *PWMPin) DutyCycle() (duty string, err error) {
	buf, err := p.read(p.pwmDutyCyclePath())
	if err != nil {
		return
	}
	return string(buf), nil
}

// SetDutyCycle writes value to pwm duty cycle path
// duty is in nanoseconds
func (p *PWMPin) SetDutyCycle(duty uint32) (err error) {
	_, err = p.write(p.pwmDutyCyclePath(), []byte(fmt.Sprintf("%v", duty)))
	return
}

// pwmPath returns pwm base path
func (p *PWMPin) pwmPath() string {
	return "/sys/class/pwm/pwmchip" + p.Chip
}

// pwmExportPath returns export path
func (p *PWMPin) pwmExportPath() string {
	return p.pwmPath() + "/export"
}

// pwmUnexportPath returns unexport path
func (p *PWMPin) pwmUnexportPath() string {
	return p.pwmPath() + "/unexport"
}

// pwmDutyCyclePath returns duty_cycle path for specified pin
func (p *PWMPin) pwmDutyCyclePath() string {
	return p.pwmPath() + "/pwm" + p.pin + "/duty_cycle"
}

// pwmPeriodPath returns period path for specified pin
func (p *PWMPin) pwmPeriodPath() string {
	return p.pwmPath() + "/pwm" + p.pin + "/period"
}

// pwmEnablePath returns enable path for specified pin
func (p *PWMPin) pwmEnablePath() string {
	return p.pwmPath() + "/pwm" + p.pin + "/enable"
}

// pwmPolarityPath returns polarity path for specified pin
func (p *PWMPin) pwmPolarityPath() string {
	return p.pwmPath() + "/pwm" + p.pin + "/polarity"
}

func writePwmFile(path string, data []byte) (i int, err error) {
	file, err := OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readPwmFile(path string) ([]byte, error) {
	file, err := OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i int
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}
