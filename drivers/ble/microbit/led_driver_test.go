package microbit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/testutil"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(testutil.NewBleTestAdaptor())
	return d
}

func TestNewLEDDriver(t *testing.T) {
	d := NewLEDDriver(testutil.NewBleTestAdaptor())
	assert.IsType(t, &LEDDriver{}, d)
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit LED"))
	assert.NotNil(t, d.Eventer)
}

func TestLEDWriteMatrix(t *testing.T) {
	d := initTestLEDDriver()
	require.NoError(t, d.WriteMatrix([]byte{0x01, 0x02}))
}

func TestLEDWriteText(t *testing.T) {
	d := initTestLEDDriver()
	require.NoError(t, d.WriteText("Hello"))
}

func TestLEDCommands(t *testing.T) {
	d := initTestLEDDriver()
	require.NoError(t, d.Blank())
	require.NoError(t, d.Solid())
	require.NoError(t, d.UpRightArrow())
	require.NoError(t, d.UpLeftArrow())
	require.NoError(t, d.DownRightArrow())
	require.NoError(t, d.DownLeftArrow())
	require.NoError(t, d.Dimond())
	require.NoError(t, d.Smile())
	require.NoError(t, d.Wink())
}
