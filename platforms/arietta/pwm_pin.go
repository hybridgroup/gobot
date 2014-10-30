package arietta

import (
	"fmt"
	"github.com/hybridgroup/gobot/internal"
	"os"
)

const (
	pwmPath = "/sys/class/pwm"
	// Default period.  10 kHz.
	period = 100000
)

type PwmDesc struct {
	name  string
	chip  int
	hwnum int
}

type Pwm struct {
	desc   *PwmDesc
	path   string
	file   internal.File
	period int
}

func (d *PwmDesc) chipPath() string {
	return fmt.Sprintf("%s/pwmchip%d", pwmPath, d.chip)
}

func newPwm(desc *PwmDesc) *Pwm {
	p := &Pwm{
		desc: desc,
		path: fmt.Sprintf("%s/pwm%d", desc.chipPath(), desc.hwnum),
	}

	writeInt(desc.hwnum, desc.chipPath(), "export")
	p.file = openOrDie(os.O_WRONLY, p.path, "duty_cycle")
	return p
}

func (p *Pwm) PwmWrite(duty_ratio byte) {
	p.setPeriod(period)
	writeInt(p.period*int(duty_ratio)/255, p.path, "duty_cycle")
}

func (p *Pwm) setPeriod(period int) {
	if period != p.period {
		writeInt(0, p.path, "enable")
		writeInt(period, p.path, "period")
		writeInt(1, p.path, "enable")
		p.period = period
	}
}

func (p *Pwm) Finalize() bool {
	p.PwmWrite(0)
	p.file.Close()
	writeInt(p.desc.hwnum, p.desc.chipPath(), "unexport")
	return true
}
