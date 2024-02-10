package sphero

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*SPRKPlusDriver)(nil)

func TestNewSPRKPlusDriver(t *testing.T) {
	d := NewSPRKPlusDriver(testutil.NewBleTestAdaptor())
	assert.NotNil(t, d.OllieDriver)
	assert.True(t, strings.HasPrefix(d.Name(), "SPRK"))
	assert.NotNil(t, d.OllieDriver)
	assert.Equal(t, d.defaultCollisionConfig, sprkplusDefaultCollisionConfig())
}

func TestNewSPRKPlusDriverWithName(t *testing.T) {
	// This is a general test, that options are applied in constructor by using the common WithName() option.	Further
	// tests for options can also be done by call of "WithOption(val).apply(cfg)".
	// arrange
	const newName = "new name"
	a := testutil.NewBleTestAdaptor()
	// act
	d := NewSPRKPlusDriver(a, ble.WithName(newName))
	// assert
	assert.Equal(t, newName, d.Name())
}
