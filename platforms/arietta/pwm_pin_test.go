package arietta

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/internal"
	"github.com/hybridgroup/gobot/mocks"
	"strconv"
	"testing"
)

var pb13 = &PwmDesc{
	name:  "PB13",
	chip:  0,
	hwnum: 2,
}

type pwmHarness struct {
	fs        *mocks.Filesystem
	export    *mocks.File
	enable    *mocks.File
	period    *mocks.File
	dutyCycle *mocks.File
}

func newPwmHarness() *pwmHarness {
	fs := mocks.NewFilesystem()
	internal.SetFilesystem(fs)

	return &pwmHarness{
		fs:        fs,
		export:    fs.Add("/sys/class/pwm/pwmchip0/export"),
		enable:    fs.Add("/sys/class/pwm/pwmchip0/pwm2/enable"),
		period:    fs.Add("/sys/class/pwm/pwmchip0/pwm2/period"),
		dutyCycle: fs.Add("/sys/class/pwm/pwmchip0/pwm2/duty_cycle"),
	}
}

func TestNewPwm(t *testing.T) {
	h := newPwmHarness()
	newPwm(pb13)

	gobot.Assert(t, h.export.Contents, "2")
}

func TestPwmWrite(t *testing.T) {
	h := newPwmHarness()
	p := newPwm(pb13)

	p.PwmWrite(123)
	gobot.Assert(t, h.enable.Contents, "1")
	gobot.Assert(t, h.period.Contents, strconv.Itoa(period))
	gobot.Assert(t, h.dutyCycle.Contents, strconv.Itoa(period*123/255))

	// The period needs to be set before the channel is enabled.
	gobot.Assert(t, h.period.Seq < h.enable.Seq, true)
}
