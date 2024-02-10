package sphero

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*BB8Driver)(nil)

func TestNewBB8Driver(t *testing.T) {
	d := NewBB8Driver(testutil.NewBleTestAdaptor())
	assert.NotNil(t, d.OllieDriver)
	assert.True(t, strings.HasPrefix(d.Name(), "BB8"))
	assert.NotNil(t, d.OllieDriver)
	assert.Equal(t, d.defaultCollisionConfig, bb8DefaultCollisionConfig())
}

func TestNewBB8DriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewBB8Driver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}
