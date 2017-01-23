package chip

import (
	"fmt"
	"io"
	"os"
)

const pwmSysfsPath = "/sys/class/pwm/pwmchip0"

type pwmControl struct {
	periodFile   *os.File
	dutyFile     *os.File
	polarityFile *os.File
	enableFile   *os.File

	duty        uint32
	periodNanos uint32
	enabled     bool
}

func exportPWM() (err error) {
	exporter, err := os.OpenFile(pwmSysfsPath + "/export", os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.WriteString(exporter, "0")
	return err
}

func unexportPWM() (err error) {
	exporter, err := os.OpenFile(pwmSysfsPath + "/unexport", os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.WriteString(exporter, "0")
	return err
}

func (c *Adaptor) initPWM() (err error) {
	const basePath = pwmSysfsPath + "/pwm0"

	if _, err = os.Stat(basePath); err != nil {
		if os.IsNotExist(err) {
			if err = exportPWM(); err != nil {
				return
			}
		} else {
			return
		}
	}

	var enableFile *os.File
	var periodFile *os.File
	var dutyFile *os.File
	var polarityFile *os.File

	defer func() {
		if enableFile != nil {
			enableFile.Close()
		}
		if periodFile != nil {
			periodFile.Close()
		}
		if dutyFile != nil {
			dutyFile.Close()
		}
		if polarityFile != nil {
			polarityFile.Close()
		}
	}()

	if enableFile, err = os.OpenFile(basePath + "/enable", os.O_WRONLY, 0666); err != nil {
		return
	}
	if periodFile, err = os.OpenFile(basePath + "/period", os.O_WRONLY, 0666); err != nil {
		return
	}
	if dutyFile, err = os.OpenFile(basePath + "/duty_cycle", os.O_WRONLY, 0666); err != nil {
		return
	}
	if polarityFile, err = os.OpenFile(basePath + "/polarity", os.O_WRONLY, 0666); err != nil {
		return
	}

	c.pwm = &pwmControl{
		enableFile:   enableFile,
		periodFile:   periodFile,
		dutyFile:     dutyFile,
		polarityFile: polarityFile,
	}

	enableFile = nil
	periodFile = nil
	dutyFile = nil
	polarityFile = nil

	return nil
}

func (c *Adaptor) closePWM() error {
	pwm := c.pwm
	if pwm != nil {
		pwm.setFrequency(0)
		pwm.setDutycycle(0)
		pwm.setEnable(false)

		if pwm.enableFile != nil {
			pwm.enableFile.Close()
		}
		if pwm.periodFile != nil {
			pwm.periodFile.Close()
		}
		if pwm.dutyFile != nil {
			pwm.dutyFile.Close()
		}
		if pwm.polarityFile != nil {
			pwm.polarityFile.Close()
		}
		if err := unexportPWM(); err != nil {
			return err
		}
		c.pwm = nil
	}
	return nil
}

func (p *pwmControl) setPolarityInverted(invPolarity bool) error {
	if p.enabled {
		polarityString := "normal"
		if invPolarity {
			polarityString = "inverted"
		}
		_, err := io.WriteString(p.polarityFile, polarityString)
		return err
	} else {
		return nil
	}
}

func (p *pwmControl) setDutycycle(duty float64) error {
	p.duty = uint32((float64(p.periodNanos) * (duty / 100.0)))
	if p.enabled {
		fmt.Printf("PWM: Setting duty cycle to %v (%v)\n", p.duty, duty)
		_, err := io.WriteString(p.dutyFile, fmt.Sprintf("%v", p.duty))
		return err
	} else {
		return nil
	}
}

func (p *pwmControl) setFrequency(freq float64) error {
	periodNanos := uint32(1e9 / freq)
	if p.enabled && (p.periodNanos != periodNanos) {
		p.periodNanos = periodNanos
		_, err := io.WriteString(p.periodFile, fmt.Sprintf("%v", periodNanos))
		return err
	} else {
		return nil
	}
}

func (p *pwmControl) setEnable(enabled bool) error {
	if p.enabled != enabled {
		p.enabled = enabled
		enableVal := 0
		if enabled {
			enableVal = 1
		}
		_, err := io.WriteString(p.enableFile, fmt.Sprintf("%v", enableVal))
		return err
	} else {
		return nil
	}
}
