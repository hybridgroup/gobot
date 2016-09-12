package mavlink

import (
	"io"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/gobottest"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

var _ gobot.Driver = (*MavlinkDriver)(nil)

func initTestMavlinkDriver() *MavlinkDriver {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }
	m.sp = nullReadWriteCloser{}
	return NewMavlinkDriver(m, "myDriver")
}

func TestMavlinkDriver(t *testing.T) {
	m := NewMavlinkAdaptor("myAdaptor", "/dev/null")
	m.sp = nullReadWriteCloser{}
	m.connect = func(port string) (io.ReadWriteCloser, error) { return nil, nil }

	d := NewMavlinkDriver(m, "myDriver")
	gobottest.Assert(t, d.Name(), "myDriver")
	gobottest.Assert(t, d.Connection().Name(), "myAdaptor")
	gobottest.Assert(t, d.interval, 10*time.Millisecond)

	d = NewMavlinkDriver(m, "myDriver", 100*time.Millisecond)
	gobottest.Assert(t, d.interval, 100*time.Millisecond)
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

	gobottest.Assert(t, len(d.Start()), 0)

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
	gobottest.Assert(t, len(d.Halt()), 0)
}
