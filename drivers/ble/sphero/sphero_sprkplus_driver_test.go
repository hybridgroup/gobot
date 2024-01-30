package sphero

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"gobot.io/x/gobot/v2"
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
