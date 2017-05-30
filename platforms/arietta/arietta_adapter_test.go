package arietta

import (
	"github.com/hybridgroup/gobot"
	"strconv"
	"testing"
)

func initTestAriettaAdaptor() *AriettaAdaptor {
	b := NewAriettaAdaptor("bot")
	return b
}

func TestAriettaAdaptorFinalize(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Finalize(), true)
}
func TestAriettaAdaptorConnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Connect(), true)
}
func TestAriettaAdaptorDisconnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Disconnect(), true)
}
func TestAriettaAdaptorReconnect(t *testing.T) {
	a := initTestAriettaAdaptor()
	gobot.Assert(t, a.Reconnect(), true)
}

func TestAriettaAdaptorPwmWrite(t *testing.T) {
	a := initTestAriettaAdaptor()
	h := newPwmHarness()

	a.PwmWrite("PB13", 78)
	gobot.Assert(t, h.enable.Contents, "1")
	gobot.Assert(t, h.period.Contents, strconv.Itoa(period))
	gobot.Assert(t, h.dutyCycle.Contents, strconv.Itoa(period*78/255))
}

func TestAriettaAdaptorDigitalRead(t *testing.T) {
	a := initTestAriettaAdaptor()
	h := newPinHarness()

	h.value.Contents = "1"
	gobot.Assert(t, a.DigitalRead("PC4"), 1)
	gobot.Assert(t, h.direction.Contents, "in")
}

func TestAriettaAdaptorDigitalWrite(t *testing.T) {
	a := initTestAriettaAdaptor()
	h := newPinHarness()

	a.DigitalWrite("PC4", 0)
	gobot.Assert(t, h.value.Contents, "0")
	gobot.Assert(t, h.direction.Contents, "out")
}
