//nolint:forcetypeassert,dupl // ok here
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

var _ gobot.Driver = (*MagnetometerDriver)(nil)

func initTestMagnetometerDriver() *MagnetometerDriver {
	d := NewMagnetometerDriver(testutil.NewBleTestAdaptor())
	return d
}

func TestMagnetometerDriver(t *testing.T) {
	d := NewMagnetometerDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &MagnetometerDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Magnetometer"))
	assert.NotNil(t, d.Eventer)
}

func TestNewMagnetometerDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewMagnetometerDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestMagnetometerStartAndHalt(t *testing.T) {
	d := initTestMagnetometerDriver()
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestMagnetometerReadData(t *testing.T) {
	sem := make(chan bool)
	a := testutil.NewBleTestAdaptor()
	d := NewMagnetometerDriver(a)
	require.NoError(t, d.Start())
	err := d.On("magnetometer", func(data interface{}) {
		assert.InDelta(t, float32(8.738), data.(*MagnetometerData).X, 0.0)
		assert.InDelta(t, float32(8.995), data.(*MagnetometerData).Y, 0.0)
		assert.InDelta(t, float32(9.252), data.(*MagnetometerData).Z, 0.0)
		sem <- true
	})
	require.NoError(t, err)

	a.SendTestDataToSubscriber([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Microbit Event \"Magnetometer\" was not published")
	}
}
