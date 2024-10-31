//nolint:forcetypeassert,nilnil // ok here
package mavlink

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	common "gobot.io/x/gobot/v2/platforms/mavlink/common"
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
	assert.NotNil(t, d.Connection())
	assert.Equal(t, 10*time.Millisecond, d.interval)

	d = NewDriver(m, 100*time.Millisecond)
	assert.Equal(t, 100*time.Millisecond, d.interval)
}

func TestMavlinkDriverName(t *testing.T) {
	d := initTestMavlinkDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Mavlink"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestMavlinkDriverStart(t *testing.T) {
	d := initTestMavlinkDriver()
	err := make(chan error)
	packet := make(chan *common.MAVLinkPacket)
	message := make(chan common.MAVLinkMessage)

	_ = d.On(PacketEvent, func(data interface{}) {
		packet <- data.(*common.MAVLinkPacket)
	})
	_ = d.On(MessageEvent, func(data interface{}) {
		message <- data.(common.MAVLinkMessage)
	})
	_ = d.On(ErrorIOEvent, func(data interface{}) {
		err <- data.(error)
	})
	_ = d.On(ErrorMAVLinkEvent, func(data interface{}) {
		err <- data.(error)
	})

	require.NoError(t, d.Start())

	select {
	case p := <-packet:
		require.NoError(t, d.SendPacket(p))

	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "packet was not emitted")
	}
	select {
	case <-message:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "message was not emitted")
	}
	select {
	case <-err:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "error was not emitted")
	}
}

func TestMavlinkDriverHalt(t *testing.T) {
	d := initTestMavlinkDriver()
	require.NoError(t, d.Halt())
}
