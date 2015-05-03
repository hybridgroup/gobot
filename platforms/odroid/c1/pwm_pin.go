package c1

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/hybridgroup/gobot/sysfs"
)

type pwmPin struct {
	pin       string
	gpioNum   int
	pwmNum    int
	pwmBase   string
}

// newPwmPin creates a new pwm pin with specified pin number
func newPwmPin(pin string, gpioNum int, pwmNum int, pwmBase string) (p *pwmPin, err error) {
	done := make(chan error, 0)
	p = &pwmPin{
		pin: pin,
		gpioNum: gpioNum,
		pwmNum: pwmNum,
		pwmBase: pwmBase,
	}

	go func() {
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/duty%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				defer fi.Close()
				if _, err = fi.WriteString("0"); err != nil {
					done <- err
				}
				fi.Sync()
				break
			}
		}
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/freq%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
				defer fi.Close()
				if _, err = fi.WriteString("0"); err != nil {
					done <- err
				}
				fi.Sync()
				break
			}
		}
		for {
			if fi, err := sysfs.OpenFile(fmt.Sprintf("%v/enable%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0644); err == nil {
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
		return p, errors.New(fmt.Sprintf("could not initialize pinNum: %v on pwmNum: %v", pin, pwmNum))
	}
}

// pwmWrite writes to a pwm pin with specified period and duty
func (p *pwmPin) pwmWrite(freq string, duty string) (err error) {
	f1, err := sysfs.OpenFile(fmt.Sprintf("%v/freq%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0666)
	defer f1.Close()
	if err != nil {
		return
	}
	if _, err = f1.WriteString(freq); err != nil {
		return
	}

	f2, err := sysfs.OpenFile(fmt.Sprintf("%v/duty%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0666)
	defer f2.Close()
	if err != nil {
		return
	}

	_, err = f2.WriteString(duty)

	return
}

// releae writes string to close a pwm pin
func (p *pwmPin) release() (err error) {
	fi, err := sysfs.OpenFile(fmt.Sprintf("%v/enable%v", p.pwmBase, p.pwmNum), os.O_WRONLY|os.O_APPEND, 0666)
	defer fi.Close()
	if err != nil {
		return
	}
	_, err = fi.WriteString("0")
	return
}
