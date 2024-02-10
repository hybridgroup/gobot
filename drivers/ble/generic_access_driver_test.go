package ble

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*GenericAccessDriver)(nil)

func TestNewGenericAccessDriver(t *testing.T) {
	// arrange
	d := NewGenericAccessDriver(testutil.NewBleTestAdaptor())
	// act
	assert.True(t, strings.HasPrefix(d.Name(), "GenericAccess"))
	assert.NotNil(t, d.Eventer)
}

func TestNewGenericAccessDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewGenericAccessDriver(a, WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}

func TestGenericAccessDriverGetDeviceName(t *testing.T) {
	// arrange
	d := NewGenericAccessDriver(testutil.NewBleTestAdaptor())
	// act
	devName, err := d.GetDeviceName()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "2a00", devName)
}

func TestGenericAccessDriverGetAppearance(t *testing.T) {
	// arrange
	a := testutil.NewBleTestAdaptor()
	d := NewGenericAccessDriver(a)
	a.SetReadCharacteristicTestFunc(func(cUUID string) ([]byte, error) {
		if cUUID == "2a01" {
			return []byte{128, 0}, nil
		}
		return nil, nil
	})
	// act
	app, err := d.GetAppearance()
	// assert
	require.NoError(t, err)
	assert.Equal(t, "Generic Computer", app)
}
