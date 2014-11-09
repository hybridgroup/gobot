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
func newPwmPin(pinNum string, ocp string) *pwmPin {
	var fi sysfs.File

	d := &pwmPin{
		pinNum: strings.ToUpper(pinNum),
	}

	pwmDevice, err := glob(fmt.Sprintf("%v/pwm_test_%v.*", ocp, d.pinNum))
	if err != nil {
		panic(err)
	}

	d.pwmDevice = pwmDevice[0]

	for i := 0; i < 10; i++ {
		fi, err = sysfs.OpenFile(fmt.Sprintf("%v/run", d.pwmDevice), os.O_RDWR|os.O_APPEND, 0666)
		defer fi.Close()
		if err != nil && i == 9 {
			panic(err)
		} else {
			break
		}
	}

	fi.WriteString("1")
	fi.Sync()

	for {
		if _, err := sysfs.OpenFile(fmt.Sprintf("%v/period", d.pwmDevice), os.O_RDONLY, 0644); err == nil {
			break
		}
	}
	for {
		if _, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", d.pwmDevice), os.O_RDONLY, 0644); err == nil {
			break
		}
	}
	return d
}

// pwmWrite writes to a pwm pin with specified period and duty
func (p *pwmPin) pwmWrite(period string, duty string) {
	f1, err := sysfs.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f1.Close()
	if err != nil {
		panic(err)
	}
	f1.WriteString(period)

	f2, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f2.Close()
	if err != nil {
		panic(err)
	}
	f2.WriteString(duty)
}

// releae writes string to close a pwm pin
func (p *pwmPin) release() {
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/run", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer fi.Close()
	if err != nil {
		panic(err)
	}
	fi.WriteString("0")
}
