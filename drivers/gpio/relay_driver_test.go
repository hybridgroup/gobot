package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*RelayDriver)(nil)

// Helper to return low/high value for testing
func (l *RelayDriver) High() bool { return l.high }

func initTestRelayDriver() (*RelayDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) (err error) {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) (err error) {
		return nil
	}
	return NewRelayDriver(a, "1"), a
}

func TestRelayDriverDefaultName(t *testing.T) {
	g, _ := initTestRelayDriver()
	assert.NotNil(t, g.Connection())
	assert.True(t, strings.HasPrefix(g.Name(), "Relay"))
}

func TestRelayDriverSetName(t *testing.T) {
	g, _ := initTestRelayDriver()
	g.SetName("mybot")
	assert.Equal(t, "mybot", g.Name())
}

func TestRelayDriverStart(t *testing.T) {
	d, _ := initTestRelayDriver()
	assert.NoError(t, d.Start())
}

func TestRelayDriverHalt(t *testing.T) {
	d, _ := initTestRelayDriver()
	assert.NoError(t, d.Halt())
}

func TestRelayDriverToggle(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}

	_ = d.Off()
	assert.False(t, d.State())
	assert.Equal(t, byte(0), lastVal)
	_ = d.Toggle()
	assert.True(t, d.State())
	assert.Equal(t, byte(1), lastVal)
	_ = d.Toggle()
	assert.False(t, d.State())
	assert.Equal(t, byte(0), lastVal)
}

func TestRelayDriverToggleInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}

	d.Inverted = true
	_ = d.Off()
	assert.False(t, d.State())
	assert.Equal(t, byte(1), lastVal)
	_ = d.Toggle()
	assert.True(t, d.State())
	assert.Equal(t, byte(0), lastVal)
	_ = d.Toggle()
	assert.False(t, d.State())
	assert.Equal(t, byte(1), lastVal)
}

func TestRelayDriverCommands(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}

	assert.Nil(t, d.Command("Off")(nil))
	assert.False(t, d.State())
	assert.Equal(t, byte(0), lastVal)

	assert.Nil(t, d.Command("On")(nil))
	assert.True(t, d.State())
	assert.Equal(t, byte(1), lastVal)

	assert.Nil(t, d.Command("Toggle")(nil))
	assert.False(t, d.State())
	assert.Equal(t, byte(0), lastVal)
}

func TestRelayDriverCommandsInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}
	d.Inverted = true

	assert.Nil(t, d.Command("Off")(nil))
	assert.True(t, d.High())
	assert.False(t, d.State())
	assert.Equal(t, byte(1), lastVal)

	assert.Nil(t, d.Command("On")(nil))
	assert.False(t, d.High())
	assert.True(t, d.State())
	assert.Equal(t, byte(0), lastVal)

	assert.Nil(t, d.Command("Toggle")(nil))
	assert.True(t, d.High())
	assert.False(t, d.State())
	assert.Equal(t, byte(1), lastVal)
}
