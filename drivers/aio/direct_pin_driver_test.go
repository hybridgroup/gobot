package aio

import (
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*DirectPinDriver)(nil)

func initTestDirectPinDriver(conn gobot.Connection) *DirectPinDriver {
	testAdaptorAnalogRead = func() (val int, err error) {
		val = 80
		return
	}
	return NewDirectPinDriver(conn, "1")
}

func TestDirectPinDriver(t *testing.T) {
	var ret map[string]interface{}

	d := initTestDirectPinDriver(newAioTestAdaptor())
	gobottest.Assert(t, d.Pin(), "1")
	gobottest.Refute(t, d.Connection(), nil)

	ret = d.Command("AnalogRead")(nil).(map[string]interface{})

	gobottest.Assert(t, ret["val"].(int), 80)
	gobottest.Assert(t, ret["err"], nil)
}

func TestDirectPinDriverStart(t *testing.T) {
	d := initTestDirectPinDriver(newAioTestAdaptor())
	gobottest.Assert(t, d.Start(), nil)
}

func TestDirectPinDriverHalt(t *testing.T) {
	d := initTestDirectPinDriver(newAioTestAdaptor())
	gobottest.Assert(t, d.Halt(), nil)
}

func TestDirectPinDriverAnalogRead(t *testing.T) {
	d := initTestDirectPinDriver(newAioTestAdaptor())
	ret, err := d.AnalogRead()
	gobottest.Assert(t, ret, 80)

	d = initTestDirectPinDriver(&aioTestBareAdaptor{})
	ret, err = d.AnalogRead()
	gobottest.Assert(t, err, ErrAnalogReadUnsupported)
}
