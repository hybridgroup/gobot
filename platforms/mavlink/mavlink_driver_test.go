package mavlink

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/hybridgroup/gobot"
	common "github.com/hybridgroup/gobot/platforms/mavlink/common"
)

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
	gobot.Assert(t, d.Name(), "myDriver")
	gobot.Assert(t, d.Connection().Name(), "myAdaptor")
	gobot.Assert(t, d.interval, 10*time.Millisecond)

	d = NewMavlinkDriver(m, "myDriver", 100*time.Millisecond)
	gobot.Assert(t, d.interval, 100*time.Millisecond)

}
func TestMavlinkDriverStart(t *testing.T) {
	d := initTestMavlinkDriver()
	err := make(chan error, 0)
	packet := make(chan *common.MAVLinkPacket, 0)
	message := make(chan common.MAVLinkMessage, 0)

	gobot.Once(d.Event("packet"), func(data interface{}) {
		packet <- data.(*common.MAVLinkPacket)
	})

	gobot.Once(d.Event("message"), func(data interface{}) {
		message <- data.(common.MAVLinkMessage)
	})
	gobot.Once(d.Event("errorIO"), func(data interface{}) {
		err <- data.(error)
	})
	gobot.Once(d.Event("errorMAVLink"), func(data interface{}) {
		err <- data.(error)
	})

	gobot.Assert(t, len(d.Start()), 0)

	select {
	case p := <-packet:
		gobot.Assert(t, d.SendPacket(p), nil)

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

	payload = []byte{0xFE, 0x09, 0x4E, 0x01, 0x01, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x03, 0x51, 0x04, 0x03, 0x1C, 0x7F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	select {
	case e := <-err:
		gobot.Assert(t, e, errors.New("Unknown Message ID: 255"))
	case <-time.After(100 * time.Millisecond):
		t.Errorf("error was not emitted")
	}

}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	gobot.Assert(t, len(d.Halt()), 0)
}
