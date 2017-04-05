package microbit

import (
	"strings"
	"testing"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(NewBleTestAdaptor())
	return d
}

func TestLEDDriver(t *testing.T) {
	d := initTestLEDDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Microbit LED"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestLEDDriverStartAndHalt(t *testing.T) {
	d := initTestLEDDriver()
	gobottest.Assert(t, d.Start(), nil)
	gobottest.Assert(t, d.Halt(), nil)
}

func TestLEDDriverWriteMatrix(t *testing.T) {
	d := initTestLEDDriver()
	d.Start()
	gobottest.Assert(t, d.WriteMatrix([]byte{0x01, 0x02}), nil)
}

func TestLEDDriverWriteText(t *testing.T) {
	d := initTestLEDDriver()
	d.Start()
	gobottest.Assert(t, d.WriteText("Hello"), nil)
}

func TestLEDDriverCommands(t *testing.T) {
	d := initTestLEDDriver()
	d.Start()
	gobottest.Assert(t, d.Blank(), nil)
	gobottest.Assert(t, d.Solid(), nil)
	gobottest.Assert(t, d.UpRightArrow(), nil)
	gobottest.Assert(t, d.UpLeftArrow(), nil)
	gobottest.Assert(t, d.DownRightArrow(), nil)
	gobottest.Assert(t, d.DownLeftArrow(), nil)
	gobottest.Assert(t, d.Dimond(), nil)
	gobottest.Assert(t, d.Smile(), nil)
	gobottest.Assert(t, d.Wink(), nil)
}
