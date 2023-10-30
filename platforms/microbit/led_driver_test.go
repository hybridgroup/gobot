package microbit

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, d.Start())
	assert.NoError(t, d.Halt())
}

func TestLEDDriverWriteMatrix(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.NoError(t, d.WriteMatrix([]byte{0x01, 0x02}))
}

func TestLEDDriverWriteText(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.NoError(t, d.WriteText("Hello"))
}

func TestLEDDriverCommands(t *testing.T) {
	d := initTestLEDDriver()
	_ = d.Start()
	assert.NoError(t, d.Blank())
	assert.NoError(t, d.Solid())
	assert.NoError(t, d.UpRightArrow())
	assert.NoError(t, d.UpLeftArrow())
	assert.NoError(t, d.DownRightArrow())
	assert.NoError(t, d.DownLeftArrow())
	assert.NoError(t, d.Dimond())
	assert.NoError(t, d.Smile())
	assert.NoError(t, d.Wink())
}
