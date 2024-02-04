package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*BatteryDriver)(nil)

func TestNewBatteryDriver(t *testing.T) {
	d := NewBatteryDriver(testutil.NewBleTestAdaptor())
	assert.True(t, strings.HasPrefix(d.Name(), "Battery"))
	assert.NotNil(t, d.Eventer)
}

func TestBatteryDriverRead(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewBatteryDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		if cUUID == "2a19" {
			return []byte{20}, nil
		}

		return nil, nil
	})

	assert.Equal(t, uint8(20), d.GetBatteryLevel())
}
