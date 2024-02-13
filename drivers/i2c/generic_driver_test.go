package i2c

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*GenericDriver)(nil)

func TestNewGenericDriver(t *testing.T) {
	// arrange
	a := newI2cTestAdaptor()
	// act
	var di interface{} = NewGenericDriver(a, "GenericI2C", 0x17)
	// assert
	d, ok := di.(*GenericDriver)
	if !ok {
		require.Fail(t, "NewGenericDriver() should have returned a *GenericDriver")
	}
	assert.NotNil(t, d.Driver)
	assert.True(t, strings.HasPrefix(d.Name(), "GenericI2C"))
}
