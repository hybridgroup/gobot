package sphero

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
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
