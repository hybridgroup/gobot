package beaglebone

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hybridgroup/gobot/sysfs"
)

type pwmPin struct {
	pinNum    string
	pwmDevice string
}

// newPwmPin creates a new pwm pin with specified pin number
func newPwmPin(pinNum string, ocp string) (p *pwmPin, err error) {
	done := make(chan error, 0)
	p = &pwmPin{
		pinNum: strings.ToUpper(pinNum),
	}

	pwmDevice, err := glob(fmt.Sprintf("%v/pwm_test_%v.*", ocp, p.pinNum))
	if err != nil {
		return
	}

	p.pwmDevice = pwmDevice[0]

	go func() {
		for {
			if _, err := sysfs.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_RDONLY, 0644); err == nil {
				break
			}
		}
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				defer fi.Close()
				if _, err = fi.WriteString("0"); err != nil {
					done <- err
				}
				fi.Sync()
				break
			}
		}
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/polarity", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				defer fi.Close()
				if _, err = fi.WriteString("0"); err != nil {
					done <- err
				}
				fi.Sync()
				break
			}
		}
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/run", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				defer fi.Close()
				if _, err = fi.WriteString("1"); err != nil {
					done <- err
				}
				fi.Sync()
				break
			}
		}
		done <- nil
	}()

	select {
	case err = <-done:
		return p, err
	case <-time.After(500 * time.Millisecond):
		return p, errors.New("could not initialize pwm device")
	}
}

// pwmWrite writes to a pwm pin with specified period and duty
func (p *pwmPin) pwmWrite(period string, duty string) (err error) {
	f1, err := sysfs.OpenFile(fmt.Sprintf("%v/period", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f1.Close()
	if err != nil {
		return
	}
	if _, err = f1.WriteString(period); err != nil {
		return
	}

	f2, err := sysfs.OpenFile(fmt.Sprintf("%v/duty", p.pwmDevice), os.O_WRONLY|os.O_APPEND, 0666)
	defer f2.Close()
	if err != nil {
		return
	}

	_, err = f2.WriteString(duty)

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
