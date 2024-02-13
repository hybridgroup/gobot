package microbit

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*TemperatureDriver)(nil)

func initTestTemperatureDriver() *TemperatureDriver {
	d := NewTemperatureDriver(testutil.NewBleTestAdaptor())
	return d
}

func TestNewTemperatureDriver(t *testing.T) {
	d := NewTemperatureDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &TemperatureDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Temperature"))
	assert.NotNil(t, d.Eventer)
}

func TestNewTemperatureDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewTemperatureDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestTemperatureStartAndHalt(t *testing.T) {
	d := initTestTemperatureDriver()
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestTemperatureReadData(t *testing.T) {
	sem := make(chan bool)
	a := testutil.NewBleTestAdaptor()
	d := NewTemperatureDriver(a)
	require.NoError(t, d.Start())
	err := d.On("temperature", func(data interface{}) {
		assert.Equal(t, int8(0x22), data)
		sem <- true
	})
	require.NoError(t, err)

	a.SendTestDataToSubscriber([]byte{0x22})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Microbit Event \"Temperature\" was not published")
	}
}
