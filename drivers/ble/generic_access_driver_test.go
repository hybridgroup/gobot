package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func TestNewGenericAccessDriver(t *testing.T) {
	d := NewGenericAccessDriver(testutil.NewBleTestAdaptor())
	assert.True(t, strings.HasPrefix(d.Name(), "GenericAccess"))
	assert.NotNil(t, d.Eventer)
}

func TestGenericAccessDriverGetDeviceName(t *testing.T) {
	d := NewGenericAccessDriver(testutil.NewBleTestAdaptor())
	assert.Equal(t, "2a00", d.GetDeviceName())
}

func TestGenericAccessDriverGetAppearance(t *testing.T) {
	a := testutil.NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		if cUUID == "2a01" {
			return []byte{128, 0}, nil
		}
		return nil, nil
	})

	assert.Equal(t, "Generic Computer", d.GetAppearance())
}
