package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func TestNewBatteryDriver(t *testing.T) {
	// arrange
	d := NewBatteryDriver(testutil.NewBleTestAdaptor())
	// act & assert
	assert.True(t, strings.HasPrefix(d.Name(), "Battery"))
	assert.NotNil(t, d.Eventer)
}

func TestNewBatteryDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewBatteryDriver(a, WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestBatteryDriverRead(t *testing.T) {
	// arrange
	a := testutil.NewBleTestAdaptor()
	d := NewBatteryDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		if cUUID == "2a19" {
			return []byte{20}, nil
		}

		return nil, nil
	})
	// act
	level, err := d.GetBatteryLevel()
	// assert
	require.NoError(t, err)
	assert.Equal(t, uint8(20), level)
}
