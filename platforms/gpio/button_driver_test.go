package gpio

import (
	"github.com/hybridgroup/gobot"
	"testing"
)

var b *ButtonDriver

func init() {
	b = NewButtonDriver(TestAdaptor{}, "bot", "1")
}

func TestButtonStart(t *testing.T) {
	gobot.Expect(t, a.Start(), true)
}

func TestButtonHalt(t *testing.T) {
	gobot.Expect(t, a.Halt(), true)
}

func TestButtonInit(t *testing.T) {
	gobot.Expect(t, a.Init(), true)
}

func TestButtonReadState(t *testing.T) {
	gobot.Expect(t, b.readState(), 1)
}

func TestButtonActive(t *testing.T) {
	b.update(1)
	gobot.Expect(t, b.Active, true)

	b.update(0)
	gobot.Expect(t, b.Active, false)
}
