package beaglebone

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type pwmPin struct {
	pinNum    string
	pwmDevice string
}

func newPwmPin(pinNum string) *pwmPin {
	var err error
	var fi *os.File

	d := new(pwmPin)
	d.pinNum = strings.ToUpper(pinNum)

	ensureSlot("am33xx_pwm")
	ensureSlot(fmt.Sprintf("bone_pwm_%v", d.pinNum))

	ocp, err := filepath.Glob(Ocp)
	if err != nil {
		panic(err)
	}

	pwmDevice, err := filepath.Glob(fmt.Sprintf("%v/pwm_test_%v.*", ocp[0], d.pinNum))
	if err != nil {
		panic(err)
	}
	d.pwmDevice = pwmDevice[0]

	for i := 0; i < 10; i++ {
		fi, err = os.OpenFile(fmt.Sprintf("%v/run", d.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil && i == 9 {
			panic(err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	fi.WriteString("1")
	fi.Close()

	return d
}

func (p *pwmPin) pwmWrite(period string, duty string) {
	var err error
	var fi *os.File

	fi, err = os.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString(period)
	fi.Close()

	fi, err = os.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString(duty)
	fi.Close()
}

func (p *pwmPin) release() {
	fi, err := os.OpenFile(fmt.Sprintf("%v/run", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	fi.WriteString("0")
	fi.Close()
}
