package arietta

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/internal"
	"github.com/hybridgroup/gobot/mocks"
	"testing"
)

var pc4 = &sysfsDigitalPinDesc{
	name:  "PC4",
	hwnum: 68,
	label: "pioC4",
}

type pinHarness struct {
	fs        *mocks.Filesystem
	export    *mocks.File
	value     *mocks.File
	direction *mocks.File
}

func newPinHarness() *pinHarness {
	fs := mocks.NewFilesystem()
	internal.SetFilesystem(fs)

	return &pinHarness{
		fs:        fs,
		export:    fs.Add("/sys/class/gpio/export"),
		value:     fs.Add("/sys/class/gpio/pioC4/value"),
		direction: fs.Add("/sys/class/gpio/pioC4/direction"),
	}
}

func TestNewSysfsDigitalPin(t *testing.T) {
	h := newPinHarness()

	newSysfsDigitalPin(pc4)
	gobot.Assert(t, h.export.Contents, "68")
}

func TestSysfsDigitalPinDigitalWrite(t *testing.T) {
	h := newPinHarness()

	p := newSysfsDigitalPin(pc4)

	p.DigitalWrite(1)
	gobot.Assert(t, h.direction.Contents, "out")
	gobot.Assert(t, h.value.Contents, "1")

	p.DigitalWrite(0)
	gobot.Assert(t, h.direction.Contents, "out")
	gobot.Assert(t, h.value.Contents, "0")
}

func TestSysfsDigitalPinDigitalRead(t *testing.T) {
	h := newPinHarness()

	p := newSysfsDigitalPin(pc4)

	h.value.Contents = "1"
	gobot.Assert(t, p.DigitalRead(), 1)
	gobot.Assert(t, h.direction.Contents, "in")

	h.value.Contents = "0"
	gobot.Assert(t, p.DigitalRead(), 0)
}

func TestSysfsDigitalPinDigitalReadWrite(t *testing.T) {
	h := newPinHarness()

	p := newSysfsDigitalPin(pc4)

	h.value.Contents = "1"
	gobot.Assert(t, p.DigitalRead(), 1)
	gobot.Assert(t, h.direction.Contents, "in")

	p.DigitalWrite(0)
	gobot.Assert(t, h.direction.Contents, "out")
	gobot.Assert(t, h.value.Contents, "0")
}
