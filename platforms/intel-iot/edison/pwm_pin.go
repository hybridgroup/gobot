package edison

import "strconv"

// pwmPath returns pwm base path
func pwmPath() string {
	return "/sys/class/pwm/pwmchip0"
}

// pwmExportPath returns export path
func pwmExportPath() string {
	return pwmPath() + "/export"
}

// pwmUnExportPath returns unexport path
func pwmUnExportPath() string {
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

type pwmPin struct {
	pin string
}

// newPwmPin returns an exported and enabled pwmPin
func newPwmPin(pin int) *pwmPin {
	p := &pwmPin{pin: strconv.Itoa(pin)}
	p.export()
	p.enable("1")
	return p
}

// enable writes value to pwm enable path
func (p *pwmPin) enable(val string) {
	_, err := writeFile(pwmEnablePath(p.pin), []byte(val))
	if err != nil {
		panic(err)
	}
}

// period reads from pwm period path and returns value
func (p *pwmPin) period() string {
	buf, err := readFile(pwmPeriodPath(p.pin))
	if err != nil {
		panic(err)
	}
	return string(buf[0 : len(buf)-1])
}

// writeDuty writes value to pwm duty cycle path
func (p *pwmPin) writeDuty(duty string) {
	_, err := writeFile(pwmDutyCyclePath(p.pin), []byte(duty))
	if err != nil {
		panic(err)
	}
}

// export writes pin to pwm export path
func (p *pwmPin) export() {
	writeFile(pwmExportPath(), []byte(p.pin))
}

// export writes pin to pwm unexport path
func (p *pwmPin) unexport() {
	writeFile(pwmUnExportPath(), []byte(p.pin))
}

// close disables and unexports pin
func (p *pwmPin) close() {
	p.enable("0")
	p.unexport()
}
