package mavlink

import (
	"io"
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/gobottest"
	common "gobot.io/x/gobot/platforms/mavlink/common"
)

var _ gobot.Driver = (*Driver)(nil)

func initTestMavlinkDriver() *Driver {
	m := NewAdaptor("/dev/null")
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	m.sp = nullReadWriteCloser{}
	return NewDriver(m)
}

func TestMavlinkDriver(t *testing.T) {
	m := NewAdaptor("/dev/null")
	m.sp = nullReadWriteCloser{}
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }

	d := NewDriver(m)
	gobottest.Refute(t, d.Connection(), nil)
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	d = NewDriver(m, 100*time.Millisecond)
	gobottest.Assert(t, d.interval, 100*time.Millisecond)
}

func TestMavlinkDriverName(t *testing.T) {
	d := initTestMavlinkDriver()
	gobottest.Assert(t, strings.HasPrefix(d.Name(), "Mavlink"), true)
	d.SetName("NewName")
	gobottest.Assert(t, d.Name(), "NewName")
}

func TestMavlinkDriverStart(t *testing.T) {
	d := initTestMavlinkDriver()
	err := make(chan error, 0)
	packet := make(chan *common.MAVLinkPacket, 0)
	message := make(chan common.MAVLinkMessage, 0)

	d.On(PacketEvent, func(data interface{}) {
		packet <- data.(*common.MAVLinkPacket)
	})
	d.On(MessageEvent, func(data interface{}) {
		message <- data.(common.MAVLinkMessage)
	})
	d.On(ErrorIOEvent, func(data interface{}) {
		err <- data.(error)
	})
	d.On(ErrorMAVLinkEvent, func(data interface{}) {
		err <- data.(error)
	})

	gobottest.Assert(t, d.Start(), nil)

	select {
	case p := <-packet:
		gobottest.Assert(t, d.SendPacket(p), nil)

	case <-time.After(100 * time.Millisecond):
		t.Errorf("packet was not emitted")
	}
	select {
	case <-message:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("message was not emitted")
	}
	select {
	case <-err:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("error was not emitted")
	}
}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	gobottest.Assert(t, d.Halt(), nil)
}
