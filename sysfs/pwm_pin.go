package sysfs

import (
	"os"
	"strconv"
)

// PWMPin is the interface for sysfs PWM interactions
type PWMPin interface {
	// Unexport unexports the pin and releases the pin from the operating system
	Unexport() error
	// Export exports the pin for use by the operating system
	Export() error
	// Enable enables the PWM pin
	Enable(val string) (err error)
	// Period returns the current PWM period for pin
	Period() (period string, err error)
	// WriteDuty writes the duty cycle to the pin
	WriteDuty(duty string) (err error)
}

type pwmPin struct {
	pin   string
	write func(path string, data []byte) (i int, err error)
	read  func(path string) ([]byte, error)
}

// NewPwmPin returns a new pwmPin
func NewPwmPin(pin int) *pwmPin {
	p := &pwmPin{pin: strconv.Itoa(pin)}
	p.read = readPwmFile
	p.write = writePwmFile
	return p
}

// Enable writes value to pwm enable path
func (p *pwmPin) Enable(val string) (err error) {
	_, err = p.write(pwmEnablePath(p.pin), []byte(val))
	return
}

// Period reads from pwm period path and returns value
func (p *pwmPin) Period() (period string, err error) {
	buf, err := p.read(pwmPeriodPath(p.pin))
	if err != nil {
		return
	}
	return string(buf), nil
}

// WriteDuty writes value to pwm duty cycle path
func (p *pwmPin) WriteDuty(duty string) (err error) {
	_, err = p.write(pwmDutyCyclePath(p.pin), []byte(duty))
	return
}

// Export writes pin to pwm export path
func (p *pwmPin) Export() (err error) {
	_, err = p.write(pwmExportPath(), []byte(p.pin))
	return
}

// Unexport writes pin to pwm unexport path
func (p *pwmPin) Unexport() (err error) {
	_, err = p.write(pwmUnexportPath(), []byte(p.pin))
	return
}

// pwmPath returns pwm base path
func pwmPath() string {
	return "/sys/class/pwm/pwmchip0"
}

// pwmExportPath returns export path
func pwmExportPath() string {
	return pwmPath() + "/export"
}

// pwmUnexportPath returns unexport path
func pwmUnexportPath() string {
	return pwmPath() + "/unexport"
}

// pwmDutyCyclePath returns duty_cycle path for specified pin
func pwmDutyCyclePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/duty_cycle"
}

// pwmPeriodPath returns period path for specified pin
func pwmPeriodPath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/period"
}

// pwmEnablePath returns enable path for specified pin
func pwmEnablePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/enable"
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
