package microbit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gobot.io/x/gobot/v2"
)

var _ gobot.Driver = (*LEDDriver)(nil)

func initTestLEDDriver() *LEDDriver {
	d := NewLEDDriver(NewBleTestAdaptor())
	return d
}

func TestLEDDriver(t *testing.T) {
	d := initTestLEDDriver()
	assert.True(t, strings.HasPrefix(d.Name(), "Microbit LED"))
	d.SetName("NewName")
	assert.Equal(t, "NewName", d.Name())
}

func TestLEDDriverStartAndHalt(t *testing.T) {
	d := initTestLEDDriver()
	require.NoError(t, d.Start())
	require.NoError(t, d.Halt())
}

func TestLEDDriverWriteMatrix(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	require.NoError(t, d.WriteMatrix([]byte{0x01, 0x02}))
}

func TestLEDDriverWriteText(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	require.NoError(t, d.WriteText("Hello"))
}

func TestLEDDriverCommands(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
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
