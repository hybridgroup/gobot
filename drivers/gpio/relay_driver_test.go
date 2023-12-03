package gpio

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
)

var _ gobot.Driver = (*RelayDriver)(nil)

// Helper to return low/high value for testing
func (l *RelayDriver) High() bool { return l.high }

func initTestRelayDriver() (*RelayDriver, *gpioTestAdaptor) {
	a := newGpioTestAdaptor()
	a.digitalWriteFunc = func(string, byte) error {
		return nil
	}
	a.pwmWriteFunc = func(string, byte) error {
		return nil
	}
	return NewRelayDriver(a, "1"), a
}

func TestNewRelayDriver(t *testing.T) {
	// arrange
	a := newGpioTestAdaptor()
	// act
	d := NewRelayDriver(a, "10")
	// assert
	assert.IsType(t, &RelayDriver{}, d)
	// assert: gpio.driver attributes
	require.NotNil(t, d.driver)
	assert.True(t, strings.HasPrefix(d.driverCfg.name, "Relay"))
	assert.Equal(t, "10", d.driverCfg.pin)
	assert.Equal(t, a, d.connection)
	assert.NotNil(t, d.afterStart)
	assert.NotNil(t, d.beforeHalt)
	assert.NotNil(t, d.Commander)
	assert.NotNil(t, d.mutex)
	// assert: driver specific attributes
	assert.False(t, d.relayCfg.inverted)
	assert.False(t, d.high)
}

func TestNewRelayDriver_options(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option, least one
	// option of this driver and one of another driver (which should lead to panic). Further tests for options can also
	// be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const (
		myName = "switch alarm relay"
	)
	panicFunc := func() {
		NewRelayDriver(newGpioTestAdaptor(), "1", WithName("crazy"), aio.WithActuatorScaler(func(float64) int { return 0 }))
	}
	// act
	d := NewRelayDriver(newGpioTestAdaptor(), "1", WithName(myName), WithRelayInverted())
	// assert
	assert.True(t, d.relayCfg.inverted)
	assert.Equal(t, myName, d.Name())
	assert.PanicsWithValue(t, "'scaler option for analog actuators' can not be applied on 'crazy'", panicFunc)
}

func TestRelayToggle(t *testing.T) {
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

func TestRelayToggleInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}

	WithRelayInverted().apply(d.relayCfg)
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

func TestRelay_Commands(t *testing.T) {
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

func TestRelay_CommandsInverted(t *testing.T) {
	d, a := initTestRelayDriver()
	var lastVal byte
	a.digitalWriteFunc = func(pin string, val byte) error {
		lastVal = val
		return nil
	}
	WithRelayInverted().apply(d.relayCfg)

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
