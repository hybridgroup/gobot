package beaglebone

import (
	"fmt"
	"github.com/hybridgroup/gobot/sysfs"
	"os"
	"strings"
)

type pwmPin struct {
	pinNum    string
	pwmDevice string
}

// newPwmPin creates a new pwm pin with specified pin number
func newPwmPin(pinNum string, ocp string) (p *pwmPin, err error) {
	var fi sysfs.File

	p = &pwmPin{
		pinNum: strings.ToUpper(pinNum),
	}

	pwmDevice, err := glob(fmt.Sprintf("%v/pwm_test_%v.*", ocp, p.pinNum))
	if err != nil {
		return
	}

	p.pwmDevice = pwmDevice[0]

	for i := 0; i < 10; i++ {
		fi, err = sysfs.OpenFile(fmt.Sprintf("%v/run", p.pwmDevice), os.O_RDWR|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil && i == 9 {
			return
		} else {
			break
		}
	}

	_, err = fi.WriteString("1")
	if err != nil {
		return
	}
	err = fi.Sync()
	if err != nil {
		return
	}

	for {
		if _, err := sysfs.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_RDONLY, 0644); err == nil {
			break
		}
	}
	for {
		if _, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_RDONLY, 0644); err == nil {
			break
		}
	}
	return
}

// pwmWrite writes to a pwm pin with specified period and duty
func (p *pwmPin) pwmWrite(period string, duty string) (err error) {
	f1, err := sysfs.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f1.Close()
	if err != nil {
		return
	}
	_, err = f1.WriteString(period)
	if err != nil {
		return
	}

	f2, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f2.Close()
	if err != nil {
		return
	}
	_, err = f2.WriteString(duty)
	if err != nil {
		return
	}

	return
}

// releae writes string to close a pwm pin
func (p *pwmPin) release() (err error) {
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/run", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer fi.Close()
	if err != nil {
		return
	}
	_, err = fi.WriteString("0")
	return
}
