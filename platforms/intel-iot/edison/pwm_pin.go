package edison

import (
	"io/ioutil"
	"strconv"
)

func pwmPath() string {
	return "/sys/class/pwm/pwmchip0"
}
func pwmExportPath() string {
	return pwmPath() + "/export"
}
func pwmUnExportPath() string {
	return pwmPath() + "/unexport"
}
func pwmDutyCyclePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/duty_cycle"
}
func pwmPeriodPath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/period"
}
func pwmEnablePath(pin string) string {
	return pwmPath() + "/pwm" + pin + "/enable"
}

type pwmPin struct {
	pin string
}

func newPwmPin(pin int) *pwmPin {
	p := &pwmPin{pin: strconv.Itoa(pin)}
	p.export()
	p.enable("1")
	return p
}

func (p *pwmPin) enable(val string) {
	err := writeFile(pwmEnablePath(p.pin), val)
	if err != nil {
		panic(err)
	}
}

func (p *pwmPin) period() string {
	buf, err := ioutil.ReadFile(pwmPeriodPath(p.pin))
	if err != nil {
		panic(err)
	}
	return string(buf[0 : len(buf)-1])
}

func (p *pwmPin) writeDuty(duty string) {
	err := writeFile(pwmDutyCyclePath(p.pin), duty)
	if err != nil {
		panic(err)
	}
}

func (p *pwmPin) export() {
	writeFile(pwmExportPath(), p.pin)
}

func (p *pwmPin) unexport() {
	writeFile(pwmUnExportPath(), p.pin)
}

func (p *pwmPin) close() {
	p.enable("0")
	p.unexport()
}
