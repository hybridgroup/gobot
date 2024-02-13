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

var _ gobot.Driver = (*AccelerometerDriver)(nil)

func TestNewAccelerometerDriver(t *testing.T) {
	d := NewAccelerometerDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &AccelerometerDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit Accelerometer"))
	assert.NotNil(t, d.Eventer)
}

func TestNewAccelerometerDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewAccelerometerDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestAccelerometerStartAndHalt(t *testing.T) {
	d := NewAccelerometerDriver(testutil.NewBleTestAdaptor())
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestAccelerometerReadData(t *testing.T) {
	sem := make(chan bool)
	a := testutil.NewBleTestAdaptor()
	d := NewAccelerometerDriver(a)
	require.NoError(t, d.Start())

	err := d.On("accelerometer", func(data interface{}) {
		assert.InDelta(t, float32(8.738), data.(*AccelerometerData).X, 0.0)
		assert.InDelta(t, float32(8.995), data.(*AccelerometerData).Y, 0.0)
		assert.InDelta(t, float32(9.252), data.(*AccelerometerData).Z, 0.0)
		sem <- true
	})

	require.NoError(t, err)

	a.SendTestDataToSubscriber([]byte{0x22, 0x22, 0x23, 0x23, 0x24, 0x24})

	select {
	case <-sem:
	case <-time.After(100 * time.Millisecond):
		require.Fail(t, "Microbit Event \"Accelerometer\" was not published")
	}
}
